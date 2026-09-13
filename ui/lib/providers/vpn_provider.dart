import 'dart:async';
import 'dart:io';
import 'package:flutter/material.dart';
import '../models/vpn_models.dart';
import '../services/android_vpn_service.dart';
import '../services/daemon_service.dart';
import '../services/subscription_service.dart';

class VpnProvider extends ChangeNotifier {
  final DaemonService _daemonService = DaemonService();
  final AndroidVpnService _androidVpnService = AndroidVpnService();
  final SubscriptionService _subService = SubscriptionService();

  VpnConnectionState _state = VpnConnectionState.disconnected;
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
    _initPlatform();
  }

  void _initPlatform() async {
    if (Platform.isAndroid) {
      await _androidVpnService.initialize(
        onStatus: (status) {
          final stateStr = status.state.toUpperCase();
          if (stateStr == 'CONNECTED') {
            _state = VpnConnectionState.connected;
            _statusMessage = 'Connected';
          } else if (stateStr == 'CONNECTING') {
            _state = VpnConnectionState.connecting;
            _statusMessage = 'Connecting...';
          } else if (stateStr == 'DISCONNECTED' || stateStr == 'STOPPED') {
            if (_state == VpnConnectionState.connected || _state == VpnConnectionState.connecting) {
              _state = VpnConnectionState.disconnected;
              _statusMessage = 'Disconnected';
              _downloadSpeed = 0;
              _uploadSpeed = 0;
            }
          }
          _uploadSpeed = status.uploadSpeed;
          _downloadSpeed = status.downloadSpeed;
          notifyListeners();
        },
      );
      // Load cached nodes immediately then fetch latest
      _nodes = await _subService.loadCachedNodes();
      notifyListeners();
      fetchNodes();
    } else {
      _startDesktopPolling();
      fetchNodes();
    }
  }

  void selectNode(ProxyNodeModel node) {
    _currentNode = node;
    _latencyMs = node.latencyMs;
    notifyListeners();
    if (isConnected) {
      connect(preselectedNode: node);
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
    _lastError = '';

    if (Platform.isAndroid) {
      _statusMessage = 'Preparing VPN...';
      notifyListeners();

      ProxyNodeModel? targetNode = preselectedNode;

      if (rawUrl.isNotEmpty) {
        targetNode = _subService.parseRawUri(rawUrl);
      }

      if (targetNode == null && _currentNode != null) {
        targetNode = _currentNode;
      }

      if (targetNode == null) {
        if (_nodes.isEmpty) {
          _statusMessage = 'Fetching servers...';
          notifyListeners();
          _nodes = await _subService.fetchAllNodes();
        }

        if (_nodes.isNotEmpty) {
          targetNode = _nodes.first;
        }
      }

      if (targetNode == null || targetNode.rawUri.isEmpty) {
        _state = VpnConnectionState.error;
        _statusMessage = 'No Servers Available';
        _lastError = 'Could not find any proxy server to connect to.';
        notifyListeners();
        return;
      }

      _currentNode = targetNode;
      _statusMessage = 'Connecting to ${targetNode.name}...';
      notifyListeners();

      final success = await _androidVpnService.connect(
        node: targetNode,
        isTun: _isTunMode,
      );

      if (success) {
        _state = VpnConnectionState.connected;
        _statusMessage = 'Connected';
        _lastError = '';
      } else {
        _state = VpnConnectionState.error;
        _statusMessage = 'Connection Failed';
        _lastError = 'Android VPN connection was refused or cancelled.';
      }
      notifyListeners();
      return;
    }

    // Desktop platforms (macOS / Linux / Windows)
    _statusMessage = _isTunMode ? 'Requesting Privileges...' : 'Starting connection...';
    notifyListeners();

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

    final success = await _daemonService.connect(
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

    if (Platform.isAndroid) {
      await _androidVpnService.disconnect();
    } else {
      await _daemonService.disconnect();
    }

    _state = VpnConnectionState.disconnected;
    _statusMessage = 'Disconnected';
    _downloadSpeed = 0;
    _uploadSpeed = 0;
    _latencyMs = 0;
    notifyListeners();
  }

  Future<void> fetchNodes() async {
    if (Platform.isAndroid) {
      final fetched = await _subService.fetchAllNodes();
      if (fetched.isNotEmpty) {
        _nodes = fetched;
        notifyListeners();
      }
    } else {
      final fetched = await _daemonService.getNodes();
      if (fetched.isNotEmpty) {
        _nodes = fetched;
        notifyListeners();
      }
    }
  }

  void _startDesktopPolling() {
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
