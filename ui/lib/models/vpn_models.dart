class ProxyNodeModel {
  final String name;
  final String server;
  final int port;
  final String protocol;
  final String country;
  final String countryFlag;
  final int latencyMs;

  final String rawUri;

  ProxyNodeModel({
    required this.name,
    required this.server,
    required this.port,
    required this.protocol,
    required this.country,
    required this.countryFlag,
    required this.latencyMs,
    this.rawUri = '',
  });

  ProxyNodeModel copyWith({
    String? name,
    String? server,
    int? port,
    String? protocol,
    String? country,
    String? countryFlag,
    int? latencyMs,
    String? rawUri,
  }) {
    return ProxyNodeModel(
      name: name ?? this.name,
      server: server ?? this.server,
      port: port ?? this.port,
      protocol: protocol ?? this.protocol,
      country: country ?? this.country,
      countryFlag: countryFlag ?? this.countryFlag,
      latencyMs: latencyMs ?? this.latencyMs,
      rawUri: rawUri ?? this.rawUri,
    );
  }

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
      rawUri: json['raw_uri'] ?? json['rawUri'] ?? '',
    );
  }
}

enum VpnConnectionState {
  disconnected,
  connecting,
  connected,
  error,
}
