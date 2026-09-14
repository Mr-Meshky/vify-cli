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
  Future<ProxyNodeModel?> findFastestWorkingNode(
    List<ProxyNodeModel> nodes, {
    void Function(String status)? onProgress,
  }) async {
    if (nodes.isEmpty) return null;

    // Take up to 40 nodes to test quickly without overloading the device
    final testBatch = nodes.take(40).toList();
    onProgress?.call('بررسی سلامت ${testBatch.length} سرور...');

    // Stage 1: Parallel TCP ping with concurrency limit
    final List<MapEntry<ProxyNodeModel, int>> tcpResults = [];
    final int concurrency = 15;

    for (int i = 0; i < testBatch.length; i += concurrency) {
      final end = (i + concurrency < testBatch.length) ? i + concurrency : testBatch.length;
      final chunk = testBatch.sublist(i, end);

      final futures = chunk.map((node) async {
        final ping = await testTcpPing(node.server, node.port);
        return MapEntry(node, ping);
      });

      final chunkResults = await Future.wait(futures);
      for (final res in chunkResults) {
        if (res.value > 0) {
          tcpResults.add(res);
        }
      }

      // Fast-pass: if we already have at least 3 low-latency candidates, break early
      final lowLatencyCount = tcpResults.where((r) => r.value > 0 && r.value < 500).length;
      if (lowLatencyCount >= 3) break;
    }

    if (tcpResults.isEmpty) {
      debugPrint('No nodes responded to TCP ping, falling back to first node');
      return nodes.first;
    }

    // Sort by lowest TCP latency
    tcpResults.sort((a, b) => a.value.compareTo(b.value));

    // Stage 2: Real verification on the top candidates (up to 5 candidates)
    final topCandidates = tcpResults.take(5).toList();

    for (int i = 0; i < topCandidates.length; i++) {
      final candidate = topCandidates[i].key;
      final tcpPing = topCandidates[i].value;

      onProgress?.call('اعتبارسنجی کیفیت اینترنت سرور ${i + 1} از ${topCandidates.length}...');

      if (Platform.isAndroid) {
        try {
          // Measure real outbound HTTP 204 delay through Xray core
          final realDelay = await _androidVpnService
              .getDelay(candidate.rawUri)
              .timeout(const Duration(milliseconds: 2500), onTimeout: () => -1);

          if (realDelay > 0) {
            debugPrint('Found verified working server: ${candidate.name} (Real delay: $realDelay ms)');
            return candidate.copyWith(latencyMs: realDelay);
          }
        } catch (e) {
          debugPrint('Error measuring real delay for ${candidate.name}: $e');
        }
      } else {
        // Desktop platforms handle real verification inside the Go daemon
        return candidate.copyWith(latencyMs: tcpPing);
      }
    }

    // If real delay timed out or was blocked by local sandbox, fallback to lowest TCP latency
    final bestTcp = topCandidates.first;
    debugPrint('Fallback to lowest TCP latency candidate: ${bestTcp.key.name} (${bestTcp.value} ms)');
    return bestTcp.key.copyWith(latencyMs: bestTcp.value);
  }

  /// Benchmark all nodes concurrently and return updated list sorted by latency (lowest first).
  Future<List<ProxyNodeModel>> benchmarkAllNodes(
    List<ProxyNodeModel> nodes, {
    void Function(int completed, int total)? onProgress,
  }) async {
    if (nodes.isEmpty) return [];

    final List<ProxyNodeModel> testedNodes = [];
    final int concurrency = 20;
    int completed = 0;

    for (int i = 0; i < nodes.length; i += concurrency) {
      final end = (i + concurrency < nodes.length) ? i + concurrency : nodes.length;
      final chunk = nodes.sublist(i, end);

      final futures = chunk.map((node) async {
        final ping = await testTcpPing(node.server, node.port);
        return node.copyWith(latencyMs: ping > 0 ? ping : -1);
      });

      final chunkResults = await Future.wait(futures);
      testedNodes.addAll(chunkResults);
      completed += chunk.length;
      onProgress?.call(completed, nodes.length);
    }

    // Sort: positive latencies ascending, then -1 (offline/timeout) at the end
    testedNodes.sort((a, b) {
      if (a.latencyMs > 0 && b.latencyMs > 0) {
        return a.latencyMs.compareTo(b.latencyMs);
      }
      if (a.latencyMs > 0 && b.latencyMs <= 0) return -1;
      if (a.latencyMs <= 0 && b.latencyMs > 0) return 1;
      return 0;
    });

    return testedNodes;
  }
}
