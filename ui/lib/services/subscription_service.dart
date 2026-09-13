import 'dart:convert';
import 'package:flutter/foundation.dart';
import 'package:http/http.dart' as http;
import 'package:shared_preferences/shared_preferences.dart';
import '../models/vpn_models.dart';

class SubscriptionService {
  static const String _cacheKey = 'vify_cached_nodes';
  static const String _subsKey = 'vify_subscription_urls';

  static final List<String> defaultSubscriptions = [
    'https://raw.githubusercontent.com/Mr-Meshky/vify/main/configs/vless.txt',
    'https://raw.githubusercontent.com/Mr-Meshky/vify/main/configs/vmess.txt',
    'https://raw.githubusercontent.com/Mr-Meshky/vify/main/configs/trojan.txt',
    'https://raw.githubusercontent.com/Mr-Meshky/vify/main/configs/ss.txt',
  ];

  Future<List<String>> getSubscriptionUrls() async {
    final prefs = await SharedPreferences.getInstance();
    return prefs.getStringList(_subsKey) ?? defaultSubscriptions;
  }

  Future<void> saveSubscriptionUrls(List<String> urls) async {
    final prefs = await SharedPreferences.getInstance();
    await prefs.setStringList(_subsKey, urls);
  }

  /// Load cached nodes from SharedPreferences for instant boot
  Future<List<ProxyNodeModel>> loadCachedNodes() async {
    try {
      final prefs = await SharedPreferences.getInstance();
      final raw = prefs.getString(_cacheKey);
      if (raw != null && raw.isNotEmpty) {
        final list = jsonDecode(raw) as List;
        return list.map((item) => ProxyNodeModel.fromJson(item)).toList();
      }
    } catch (e) {
      debugPrint('Error loading cached nodes: $e');
    }
    return [];
  }

  /// Save parsed nodes to cache
  Future<void> cacheNodes(List<ProxyNodeModel> nodes) async {
    try {
      final prefs = await SharedPreferences.getInstance();
      final data = nodes.map((n) => {
        'name': n.name,
        'server': n.server,
        'port': n.port,
        'protocol': n.protocol,
        'country': n.country,
        'country_flag': n.countryFlag,
        'latency_ms': n.latencyMs,
        'raw_uri': n.rawUri,
      }).toList();
      await prefs.setString(_cacheKey, jsonEncode(data));
    } catch (e) {
      debugPrint('Error saving cached nodes: $e');
    }
  }

  /// Fetch nodes from all subscription URLs
  Future<List<ProxyNodeModel>> fetchAllNodes() async {
    final urls = await getSubscriptionUrls();
    final List<ProxyNodeModel> result = [];
    final Set<String> seenUris = {};

    for (final url in urls) {
      try {
        final res = await http.get(Uri.parse(url)).timeout(
          const Duration(seconds: 8),
        );
        if (res.statusCode == 200 && res.body.trim().isNotEmpty) {
          final nodes = parseSubscriptionContent(res.body.trim());
          for (final node in nodes) {
            if (node.rawUri.isNotEmpty && !seenUris.contains(node.rawUri)) {
              seenUris.add(node.rawUri);
              result.add(node);
            }
          }
        }
      } catch (e) {
        debugPrint('Error fetching subscription from $url: $e');
      }
    }

    if (result.isNotEmpty) {
      await cacheNodes(result);
      return result;
    }

    // Fallback to cache if network failed
    return await loadCachedNodes();
  }

  /// Parse subscription text which may be plain URLs or base64 encoded
  List<ProxyNodeModel> parseSubscriptionContent(String content) {
    String text = content.trim();
    // Try base64 decoding if content looks like base64
    if (!text.contains('\n') && !text.contains('://')) {
      try {
        text = utf8.decode(base64Decode(base64.normalize(text)));
      } catch (_) {}
    }

    final lines = text.split(RegExp(r'[\r\n]+'));
    final List<ProxyNodeModel> nodes = [];

    for (var line in lines) {
      line = line.trim();
      if (line.isEmpty) continue;
      final node = parseRawUri(line);
      if (node != null) {
        nodes.add(node);
      }
    }

    return nodes;
  }

  /// Parse a single URI (vless, vmess, trojan, ss)
  ProxyNodeModel? parseRawUri(String rawUri) {
    final uriStr = rawUri.trim();
    if (uriStr.isEmpty) return null;

    try {
      if (uriStr.startsWith('vless://')) {
        return _parseVless(uriStr);
      } else if (uriStr.startsWith('vmess://')) {
        return _parseVmess(uriStr);
      } else if (uriStr.startsWith('trojan://')) {
        return _parseTrojan(uriStr);
      } else if (uriStr.startsWith('ss://')) {
        return _parseShadowsocks(uriStr);
      }
    } catch (e) {
      debugPrint('Failed to parse URI $uriStr: $e');
    }
    return null;
  }

  ProxyNodeModel? _parseVless(String rawUri) {
    final uri = Uri.parse(rawUri);
    final port = uri.port > 0 ? uri.port : 443;
    final server = uri.host;
    String name = uri.fragment.isNotEmpty
        ? Uri.decodeComponent(uri.fragment)
        : 'VLESS-$server:$port';

    final flag = detectCountryFlag(name);
    final country = detectCountryName(name);

    return ProxyNodeModel(
      name: name,
      server: server,
      port: port,
      protocol: 'vless',
      country: country,
      countryFlag: flag,
      latencyMs: 0,
      rawUri: rawUri,
    );
  }

  ProxyNodeModel? _parseVmess(String rawUri) {
    final b64 = rawUri.substring('vmess://'.length).trim();
    try {
      final jsonStr = utf8.decode(base64Decode(base64.normalize(b64)));
      final data = jsonDecode(jsonStr);
      final server = data['add']?.toString() ?? '';
      final port = int.tryParse(data['port']?.toString() ?? '443') ?? 443;
      final name = (data['ps'] != null && data['ps'].toString().isNotEmpty)
          ? data['ps'].toString()
          : 'VMess-$server:$port';

      final flag = detectCountryFlag(name);
      final country = detectCountryName(name);

      return ProxyNodeModel(
        name: name,
        server: server,
        port: port,
        protocol: 'vmess',
        country: country,
        countryFlag: flag,
        latencyMs: 0,
        rawUri: rawUri,
      );
    } catch (_) {
      return null;
    }
  }

  ProxyNodeModel? _parseTrojan(String rawUri) {
    final uri = Uri.parse(rawUri);
    final port = uri.port > 0 ? uri.port : 443;
    final server = uri.host;
    String name = uri.fragment.isNotEmpty
        ? Uri.decodeComponent(uri.fragment)
        : 'Trojan-$server:$port';

    final flag = detectCountryFlag(name);
    final country = detectCountryName(name);

    return ProxyNodeModel(
      name: name,
      server: server,
      port: port,
      protocol: 'trojan',
      country: country,
      countryFlag: flag,
      latencyMs: 0,
      rawUri: rawUri,
    );
  }

  ProxyNodeModel? _parseShadowsocks(String rawUri) {
    final uri = Uri.parse(rawUri);
    String name = uri.fragment.isNotEmpty
        ? Uri.decodeComponent(uri.fragment)
        : 'Shadowsocks';
    String server = uri.host;
    int port = uri.port > 0 ? uri.port : 8388;

    final flag = detectCountryFlag(name);
    final country = detectCountryName(name);

    return ProxyNodeModel(
      name: name,
      server: server,
      port: port,
      protocol: 'shadowsocks',
      country: country,
      countryFlag: flag,
      latencyMs: 0,
      rawUri: rawUri,
    );
  }

  static String detectCountryFlag(String text) {
    final upper = text.toUpperCase();
    if (upper.contains('DE') || upper.contains('GERMAN')) return '🇩🇪';
    if (upper.contains('NL') || upper.contains('NETHERLAND')) return '🇳🇱';
    if (upper.contains('US') || upper.contains('UNITED STATES') || upper.contains('AMERICA')) return '🇺🇸';
    if (upper.contains('UK') || upper.contains('GB') || upper.contains('BRITAIN')) return '🇬🇧';
    if (upper.contains('FR') || upper.contains('FRANCE')) return '🇫🇷';
    if (upper.contains('TR') || upper.contains('TURKEY')) return '🇹🇷';
    if (upper.contains('FI') || upper.contains('FINLAND')) return '🇫🇮';
    if (upper.contains('SE') || upper.contains('SWEDEN')) return '🇸🇪';
    if (upper.contains('CA') || upper.contains('CANADA')) return '🇨🇦';
    if (upper.contains('SG') || upper.contains('SINGAPORE')) return '🇸🇬';
    if (upper.contains('JP') || upper.contains('JAPAN')) return '🇯🇵';
    if (upper.contains('CH') || upper.contains('SWISS')) return '🇨🇭';
    if (upper.contains('AE') || upper.contains('DUBAI') || upper.contains('UAE')) return '🇦🇪';
    if (upper.contains('PL') || upper.contains('POLAND')) return '🇵🇱';
    if (upper.contains('IT') || upper.contains('ITALY')) return '🇮🇹';
    return '🌐';
  }

  static String detectCountryName(String text) {
    final upper = text.toUpperCase();
    if (upper.contains('DE') || upper.contains('GERMAN')) return 'Germany';
    if (upper.contains('NL') || upper.contains('NETHERLAND')) return 'Netherlands';
    if (upper.contains('US') || upper.contains('UNITED STATES')) return 'United States';
    if (upper.contains('UK') || upper.contains('GB')) return 'United Kingdom';
    if (upper.contains('FR') || upper.contains('FRANCE')) return 'France';
    if (upper.contains('TR') || upper.contains('TURKEY')) return 'Turkey';
    if (upper.contains('FI') || upper.contains('FINLAND')) return 'Finland';
    if (upper.contains('SE') || upper.contains('SWEDEN')) return 'Sweden';
    if (upper.contains('CA') || upper.contains('CANADA')) return 'Canada';
    if (upper.contains('SG') || upper.contains('SINGAPORE')) return 'Singapore';
    if (upper.contains('JP') || upper.contains('JAPAN')) return 'Japan';
    if (upper.contains('CH') || upper.contains('SWISS')) return 'Switzerland';
    if (upper.contains('AE') || upper.contains('UAE')) return 'UAE';
    if (upper.contains('PL') || upper.contains('POLAND')) return 'Poland';
    if (upper.contains('IT') || upper.contains('ITALY')) return 'Italy';
    return 'GLOBAL';
  }
}
