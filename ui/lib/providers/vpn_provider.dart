import 'dart:async';
import 'package:flutter/material.dart';
import '../models/vpn_models.dart';
import '../services/daemon_service.dart';

class VpnProvider extends ChangeNotifier {
  final DaemonService _daemonService = DaemonService();

  VpnConnectionState _state = VpnConnectionState.disconnected;
  VpnTargetMode _selectedTarget = VpnTargetMode.fastPass;
  ProxyNodeModel? _currentNode;
  
  int _downloadSpeed = 0; // Bytes/sec
  int _uploadSpeed = 0;   // Bytes/sec
  int _latencyMs = 0;
  String _statusMessage = 'Disconnected';
  String _lastError = '';
  bool _isTunMode = true;
  bool _isDaemonConnected = false;

  List<ProxyNodeModel> _nodes = [];
  Timer? _pollingTimer;

  VpnConnectionState get state => _state;
  VpnTargetMode get selectedTarget => _selectedTarget;
  ProxyNodeModel? get currentNode => _currentNode;
  int get downloadSpeed => _downloadSpeed;
  int get uploadSpeed => _uploadSpeed;
  int get latencyMs => _latencyMs;
  String get statusMessage => _statusMessage;
  String get lastError => _lastError;
  bool get isTunMode => _isTunMode;
  bool get isDaemonConnected => _isDaemonConnected;
  List<ProxyNodeModel> get nodes => _nodes;

  bool get isConnected => _state == VpnConnectionState.connected;
  bool get isConnecting => _state == VpnConnectionState.connecting;
  bool get hasError => _state == VpnConnectionState.error;

  VpnProvider() {
    _startPolling();
    fetchNodes();
  }

  void selectTarget(VpnTargetMode target) {
    if (_selectedTarget == target) return;
    _selectedTarget = target;
    notifyListeners();
    if (isConnected) {
      toggleConnection();
    }
  }

  void toggleTunMode(bool isTun) {
    if (_isTunMode == isTun) return;
    _isTunMode = isTun;
    notifyListeners();
    if (isConnected) {
      disconnect();
    }
  }

  Future<void> toggleConnection() async {
    if (isConnected || isConnecting) {
      await disconnect();
    } else {
      await connect();
    }
  }

  Future<void> connect({String rawUrl = '', ProxyNodeModel? preselectedNode}) async {
    _state = VpnConnectionState.connecting;
    _statusMessage = _isTunMode ? 'Requesting Privileges...' : 'Starting connection...';
    _lastError = '';
    notifyListeners();

    // 1. Ensure Go Core Daemon is running (prompt admin if TUN mode on macOS/Linux)
    final daemonReady = await _daemonService.ensureDaemonRunning(isTun: _isTunMode);
    if (!daemonReady) {
      _state = VpnConnectionState.error;
      _statusMessage = _isTunMode ? 'Admin Privileges Denied' : 'Daemon failed to start';
      _lastError = _isTunMode
          ? 'TUN mode requires root privileges. Please grant administrator access or switch to System Proxy.'
          : 'Could not communicate with background engine.';
      notifyListeners();
      return;
    }

    _statusMessage = 'Finding fastest server...';
    notifyListeners();

    // 2. Trigger real connection via Go daemon
    final success = await _daemonService.connect(
      target: _selectedTarget,
      mode: _isTunMode ? 'tun' : 'system_proxy',
      rawUrl: rawUrl,
    );

    if (!success) {
      _state = VpnConnectionState.error;
      _statusMessage = 'Connection Failed';
      _lastError = 'Failed to dispatch connection request to engine.';
      notifyListeners();
    }
  }

  Future<void> disconnect() async {
    _state = VpnConnectionState.connecting;
    _statusMessage = 'Disconnecting...';
    notifyListeners();

    await _daemonService.disconnect();

    _state = VpnConnectionState.disconnected;
    _statusMessage = 'Disconnected';
    _downloadSpeed = 0;
    _uploadSpeed = 0;
    _latencyMs = 0;
    _currentNode = null;
    notifyListeners();
  }

  Future<void> fetchNodes() async {
    final fetched = await _daemonService.getNodes();
    if (fetched.isNotEmpty) {
      _nodes = fetched;
      notifyListeners();
    }
  }

  void _startPolling() {
    _pollingTimer = Timer.periodic(const Duration(seconds: 1), (timer) async {
      final daemonAlive = await _daemonService.isDaemonRunning();
      _isDaemonConnected = daemonAlive;

      if (daemonAlive) {
        final status = await _daemonService.getStatus();
        if (status != null) {
          final bool isConn = status['connected'] ?? false;
          final bool isConnPending = status['connecting'] ?? false;
          final String err = status['last_error'] ?? '';

          if (err.isNotEmpty && !isConn && !isConnPending) {
            _state = VpnConnectionState.error;
            _statusMessage = 'Connection Error';
            _lastError = err;
          } else if (isConn) {
            _state = VpnConnectionState.connected;
            _statusMessage = 'Connected';
            if (status['current_node'] != null) {
              _currentNode = ProxyNodeModel.fromJson(status['current_node']);
              _latencyMs = _currentNode!.latencyMs;
            }
            _downloadSpeed = status['download_speed'] ?? 0;
            _uploadSpeed = status['upload_speed'] ?? 0;
          } else if (isConnPending) {
            _state = VpnConnectionState.connecting;
            _statusMessage = status['status_message'] ?? 'Connecting...';
          } else {
            if (_state != VpnConnectionState.disconnected && _state != VpnConnectionState.error) {
              _state = VpnConnectionState.disconnected;
              _statusMessage = 'Disconnected';
              _downloadSpeed = 0;
              _uploadSpeed = 0;
            }
          }
          notifyListeners();
        }
      }
    });
  }

  @override
  void dispose() {
    _pollingTimer?.cancel();
    super.dispose();
  }
}
