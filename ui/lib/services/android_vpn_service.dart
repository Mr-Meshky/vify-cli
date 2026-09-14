import 'dart:async';
import 'dart:io';
import 'package:flutter/foundation.dart';
import 'package:flutter_v2ray/flutter_v2ray.dart';
import '../models/vpn_models.dart';

class AndroidVpnService {
  static final AndroidVpnService _instance = AndroidVpnService._internal();
  factory AndroidVpnService() => _instance;
  AndroidVpnService._internal();

  FlutterV2ray? _flutterV2ray;
  bool _isInitialized = false;

  void Function(V2RayStatus status)? onStatusUpdate;

  Future<void> initialize({void Function(V2RayStatus status)? onStatus}) async {
    if (!Platform.isAndroid) return;
    if (_isInitialized) return;

    onStatusUpdate = onStatus;

    _flutterV2ray = FlutterV2ray(
      onStatusChanged: (status) {
        debugPrint('V2Ray status: ${status.state}, Up: ${status.uploadSpeed}, Down: ${status.downloadSpeed}');
        onStatusUpdate?.call(status);
      },
    );

    try {
      await _flutterV2ray!.initializeV2Ray();
      _isInitialized = true;
      debugPrint('V2Ray initialized successfully on Android');
    } catch (e) {
      debugPrint('Failed to initialize V2Ray: $e');
    }
  }

  Future<bool> requestPermission() async {
    if (!Platform.isAndroid) return true;
    if (_flutterV2ray == null) await initialize();
    try {
      return await _flutterV2ray!.requestPermission();
    } catch (e) {
      debugPrint('Error requesting VPN permission: $e');
      return false;
    }
  }

  Future<bool> connect({
    required ProxyNodeModel node,
    bool isTun = true,
  }) async {
    if (!Platform.isAndroid) return false;
    if (_flutterV2ray == null) await initialize();

    final hasPerm = await requestPermission();
    if (!hasPerm) {
      debugPrint('Android VPN permission not granted');
      return false;
    }

    try {
      final parser = FlutterV2ray.parseFromURL(node.rawUri);
      final config = parser.getFullConfiguration();

      // Start V2Ray service:
      // Passing bypassSubnets as null sets builder.addRoute("0.0.0.0", 0) in Android VpnService,
      // guaranteeing full device TUN routing without route fragmentation or packet drops.
      await _flutterV2ray!.startV2Ray(
        remark: node.name.isNotEmpty ? node.name : 'Vify-${node.country}',
        config: config,
        proxyOnly: !isTun,
        bypassSubnets: null,
        notificationDisconnectButtonName: 'DISCONNECT',
      );

      return true;
    } catch (e) {
      debugPrint('Error starting Android V2Ray: $e');
      return false;
    }
  }

  Future<void> disconnect() async {
    if (!Platform.isAndroid) return;
    try {
      await _flutterV2ray?.stopV2Ray();
    } catch (e) {
      debugPrint('Error stopping Android V2Ray: $e');
    }
  }

  Future<int> getDelay(String rawUri) async {
    if (!Platform.isAndroid || _flutterV2ray == null) return -1;
    try {
      final parser = FlutterV2ray.parseFromURL(rawUri);
      final config = parser.getFullConfiguration();
      return await _flutterV2ray!.getServerDelay(
        config: config,
        url: 'https://cp.cloudflare.com/generate_204',
      );
    } catch (e) {
      return -1;
    }
  }
}
