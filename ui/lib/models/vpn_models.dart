class ProxyNodeModel {
  final String name;
  final String server;
  final int port;
  final String protocol;
  final String country;
  final String countryFlag;
  final int latencyMs;

  ProxyNodeModel({
    required this.name,
    required this.server,
    required this.port,
    required this.protocol,
    required this.country,
    required this.countryFlag,
    required this.latencyMs,
  });

  factory ProxyNodeModel.fromJson(Map<String, dynamic> json) {
    int latency = 0;
    if (json['latency'] != null) {
      if (json['latency'] is int) {
        latency = (json['latency'] as int) ~/ 1000000; // if nanoseconds
      } else if (json['latency_ms'] != null) {
        latency = json['latency_ms'];
      }
    }
    return ProxyNodeModel(
      name: json['name'] ?? 'Unknown Server',
      server: json['server'] ?? '',
      port: json['port'] ?? 443,
      protocol: json['protocol'] ?? 'vless',
      country: json['country'] ?? 'GLOBAL',
      countryFlag: json['country_flag'] ?? '🌐',
      latencyMs: latency,
    );
  }
}

enum VpnConnectionState {
  disconnected,
  connecting,
  connected,
  error,
}

enum VpnTargetMode {
  fastPass,
  aiGemini,
  streaming,
  gaming,
}

extension VpnTargetModeExtension on VpnTargetMode {
  String get label {
    switch (this) {
      case VpnTargetMode.fastPass:
        return 'Fast Pass';
      case VpnTargetMode.aiGemini:
        return 'AI & Gemini';
      case VpnTargetMode.streaming:
        return 'Streaming';
      case VpnTargetMode.gaming:
        return 'Gaming';
    }
  }

  String get icon {
    switch (this) {
      case VpnTargetMode.fastPass:
        return '⚡';
      case VpnTargetMode.aiGemini:
        return '🤖';
      case VpnTargetMode.streaming:
        return '🎬';
      case VpnTargetMode.gaming:
        return '🎮';
    }
  }

  String get targetParam {
    switch (this) {
      case VpnTargetMode.fastPass:
        return '';
      case VpnTargetMode.aiGemini:
        return 'gemini';
      case VpnTargetMode.streaming:
        return 'streaming';
      case VpnTargetMode.gaming:
        return 'gaming';
    }
  }
}
