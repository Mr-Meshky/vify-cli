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

      // Domestic & LAN bypass subnets
      final List<String> bypassSubnets = [
        "0.0.0.0/5",
        "8.0.0.0/7",
        "11.0.0.0/8",
        "12.0.0.0/6",
        "16.0.0.0/4",
        "32.0.0.0/3",
        "64.0.0.0/2",
        "128.0.0.0/3",
        "160.0.0.0/5",
        "168.0.0.0/6",
        "172.0.0.0/12",
        "172.32.0.0/11",
        "172.64.0.0/10",
        "172.128.0.0/9",
        "173.0.0.0/8",
        "174.0.0.0/7",
        "176.0.0.0/4",
        "192.0.0.0/9",
        "192.128.0.0/11",
        "192.160.0.0/13",
        "192.169.0.0/16",
        "192.170.0.0/15",
        "192.172.0.0/14",
        "192.176.0.0/12",
        "192.192.0.0/10",
        "193.0.0.0/8",
        "194.0.0.0/7",
        "196.0.0.0/6",
        "200.0.0.0/5",
        "208.0.0.0/4",
        "240.0.0.0/4",
      ];

      await _flutterV2ray!.startV2Ray(
        remark: node.name.isNotEmpty ? node.name : 'Vify-${node.country}',
        config: config,
        proxyOnly: !isTun,
        bypassSubnets: isTun ? bypassSubnets : null,
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
      return await _flutterV2ray!.getServerDelay(
        config: parser.getFullConfiguration(),
        url: 'https://cp.cloudflare.com/generate_204',
      );
    } catch (e) {
      return -1;
    }
  }
}
