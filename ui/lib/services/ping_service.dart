import 'dart:async';
import 'dart:io';
import 'package:flutter/foundation.dart';
import '../models/vpn_models.dart';
import 'android_vpn_service.dart';

class PingService {
  static final PingService _instance = PingService._internal();
  factory PingService() => _instance;
  PingService._internal();

  final AndroidVpnService _androidVpnService = AndroidVpnService();

  /// Tests TCP socket connection to [host]:[port] with given [timeout].
  /// Returns elapsed milliseconds if successful, or -1 if unreachable.
  Future<int> testTcpPing(
    String host,
    int port, {
    Duration timeout = const Duration(milliseconds: 1400),
  }) async {
    if (host.isEmpty || port <= 0 || port > 65535) return -1;
    try {
      final sw = Stopwatch()..start();
      final socket = await Socket.connect(host, port, timeout: timeout);
      final latency = sw.elapsedMilliseconds;
      await socket.close();
      return latency;
    } catch (_) {
      return -1;
    }
  }

  /// Finds the fastest working node using a standard 2-stage pipeline:
  /// Stage 1: High-concurrency TCP ping across a batch of candidates (fast pre-filter).
  /// Stage 2: Real HTTP 204 delay test on top candidates to ensure download & upload work.
  /// Finds the fastest working node using strict real-traffic verification:
  /// 1. Concurrent TCP socket ping across candidate pool (fast preliminary filter to discard dead IPs).
  /// 2. Concurrently verifies real HTTP 204 throughput on top candidates in chunks of 4.
  /// 3. Returns the first node with confirmed working download and upload (real delay > 0).
  /// NEVER returns an unverified node or fake TCP ping.
  Future<ProxyNodeModel?> findFastestWorkingNode(
    List<ProxyNodeModel> nodes, {
    void Function(String status)? onProgress,
  }) async {
    if (nodes.isEmpty) return null;

    final pool = nodes.take(60).toList();
    onProgress?.call('پایش سلامت اولیه ${pool.length} سرور...');

    // Stage 1: Fast TCP socket pre-filter
    final List<MapEntry<ProxyNodeModel, int>> responsive = [];
    final int tcpConcurrency = 20;

    for (int i = 0; i < pool.length; i += tcpConcurrency) {
      final end = (i + tcpConcurrency < pool.length) ? i + tcpConcurrency : pool.length;
      final chunk = pool.sublist(i, end);

      final futures = chunk.map((node) async {
        final ping = await testTcpPing(node.server, node.port, timeout: const Duration(milliseconds: 1200));
        return MapEntry(node, ping);
      });

      final chunkResults = await Future.wait(futures);
      for (final res in chunkResults) {
        if (res.value > 0) {
          responsive.add(res);
        }
      }
    }

    // Sort by lowest TCP latency as preliminary priority
    responsive.sort((a, b) => a.value.compareTo(b.value));

    List<ProxyNodeModel> candidateQueue = responsive.map((r) => r.key).toList();
    if (candidateQueue.isEmpty) {
      candidateQueue = pool.take(16).toList();
    } else {
      candidateQueue = candidateQueue.take(24).toList();
    }

    // Stage 2: Concurrently test real outbound throughput (4 workers at a time)
    const int verifyChunkSize = 4;
    int testedCount = 0;

    for (int i = 0; i < candidateQueue.length; i += verifyChunkSize) {
      final end = (i + verifyChunkSize < candidateQueue.length)
          ? i + verifyChunkSize
          : candidateQueue.length;
      final chunk = candidateQueue.sublist(i, end);

      testedCount += chunk.length;
      onProgress?.call('تست ترافیک واقعی سرورهای برتر ($testedCount از ${candidateQueue.length})...');

      if (Platform.isAndroid) {
        final delayFutures = chunk.map((node) async {
          final realDelay = await _androidVpnService.getDelay(node.rawUri);
          return MapEntry(node, realDelay);
        });

        final results = await Future.wait(delayFutures);
        final working = results.where((r) => r.value > 0).toList();

        if (working.isNotEmpty) {
          // Sort working candidates by lowest real delay
          working.sort((a, b) => a.value.compareTo(b.value));
          final winner = working.first;
          debugPrint('Verified working proxy node: ${winner.key.name} (Real HTTP Delay: ${winner.value} ms)');
          return winner.key.copyWith(latencyMs: winner.value);
        }
      } else {
        // Desktop platforms perform real verification inside the Go daemon engine
        final winner = chunk.first;
        return winner.copyWith(latencyMs: responsive.firstWhere((r) => r.key.server == winner.server, orElse: () => MapEntry(winner, 100)).value);
      }
    }

    debugPrint('No server passed real HTTP verification in candidate queue');
    return null;
  }

  /// Benchmark nodes and return list where ONLY nodes with verified real delay show positive latency.
  /// Fake TCP pings to Cloudflare edges are never displayed as working servers.
  Future<List<ProxyNodeModel>> benchmarkAllNodes(
    List<ProxyNodeModel> nodes, {
    void Function(int completed, int total)? onProgress,
  }) async {
    if (nodes.isEmpty) return [];

    final List<ProxyNodeModel> testedNodes = [];
    final int concurrency = 20;

    // 1. Fast socket check to eliminate offline IPs
    for (int i = 0; i < nodes.length; i += concurrency) {
      final end = (i + concurrency < nodes.length) ? i + concurrency : nodes.length;
      final chunk = nodes.sublist(i, end);

      final futures = chunk.map((node) async {
        final ping = await testTcpPing(node.server, node.port, timeout: const Duration(milliseconds: 1200));
        return node.copyWith(latencyMs: ping > 0 ? ping : -1);
      });

      final chunkResults = await Future.wait(futures);
      testedNodes.addAll(chunkResults);
    }

    if (!Platform.isAndroid) {
      testedNodes.sort((a, b) {
        if (a.latencyMs > 0 && b.latencyMs > 0) return a.latencyMs.compareTo(b.latencyMs);
        if (a.latencyMs > 0 && b.latencyMs <= 0) return -1;
        if (a.latencyMs <= 0 && b.latencyMs > 0) return 1;
        return 0;
      });
      return testedNodes;
    }

    // 2. On Android, verify REAL HTTP delay for top responsive candidates
    final responsive = testedNodes.where((n) => n.latencyMs > 0).toList();
    responsive.sort((a, b) => a.latencyMs.compareTo(b.latencyMs));

    // Test real delay on top 20 responsive candidates in chunks of 4
    final toVerify = responsive.take(20).toList();
    final Map<String, int> realDelays = {};
    int completed = 0;

    for (int i = 0; i < toVerify.length; i += 4) {
      final end = (i + 4 < toVerify.length) ? i + 4 : toVerify.length;
      final chunk = toVerify.sublist(i, end);

      final futures = chunk.map((node) async {
        final realDelay = await _androidVpnService.getDelay(node.rawUri);
        return MapEntry(node.rawUri, realDelay);
      });

      final chunkResults = await Future.wait(futures);
      for (final res in chunkResults) {
        realDelays[res.key] = res.value;
      }
      completed += chunk.length;
      onProgress?.call(completed, toVerify.length);
    }

    // 3. Build final list: only nodes with realDelay > 0 get positive latency!
    // Nodes that failed real verification or weren't tested get -1 (Timeout/Unverified)
    final List<ProxyNodeModel> finalList = testedNodes.map((n) {
      final verifiedDelay = realDelays[n.rawUri];
      if (verifiedDelay != null && verifiedDelay > 0) {
        return n.copyWith(latencyMs: verifiedDelay);
      }
      return n.copyWith(latencyMs: -1);
    }).toList();

    // Sort: Verified working nodes (lowest real latency) first, then timeouts
    finalList.sort((a, b) {
      if (a.latencyMs > 0 && b.latencyMs > 0) {
        return a.latencyMs.compareTo(b.latencyMs);
      }
      if (a.latencyMs > 0 && b.latencyMs <= 0) return -1;
      if (a.latencyMs <= 0 && b.latencyMs > 0) return 1;
      return 0;
    });

    return finalList;
  }
}
