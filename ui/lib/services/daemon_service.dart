import 'dart:async';
import 'dart:convert';
import 'dart:io';
import 'package:flutter/foundation.dart';
import 'package:http/http.dart' as http;
import '../models/vpn_models.dart';

class DaemonService {
  static const String baseUrl = 'http://127.0.0.1:28190';

  static String findVifyBinary() {
    try {
      final currentExe = File(Platform.resolvedExecutable).resolveSymbolicLinksSync().toLowerCase();
      final appDir = File(Platform.resolvedExecutable).parent.path;

      // Candidate helper binaries (dedicated locations, never the GUI executable)
      final exeCandidates = [
        // macOS bundle helper location
        '$appDir/../Resources/vify',
        // Windows dedicated CLI/daemon names
        '$appDir/vify-cli.exe',
        '$appDir/bin/vify.exe',
        // Linux/Unix dedicated CLI/daemon names
        '$appDir/bin/vify',
        '$appDir/vify-cli',
        '$appDir/vify-daemon',
      ];
      for (final p in exeCandidates) {
        final f = File(p);
        if (f.existsSync()) {
          try {
            if (f.resolveSymbolicLinksSync().toLowerCase() != currentExe) {
              return p;
            }
          } catch (_) {
            return p;
          }
        }
      }
    } catch (_) {}

    final home = Platform.environment['HOME'] ?? Platform.environment['USERPROFILE'] ?? '';
    final candidates = [
      '$home/.local/bin/vify',
      '$home/.local/bin/vify.exe',
      '/usr/local/bin/vify',
      '/opt/homebrew/bin/vify',
      '/opt/vify/vify',
      '/usr/bin/vify',
    ];
    for (final p in candidates) {
      final f = File(p);
      if (f.existsSync()) {
        try {
          final currentExe = File(Platform.resolvedExecutable).resolveSymbolicLinksSync().toLowerCase();
          if (f.resolveSymbolicLinksSync().toLowerCase() != currentExe) {
            return p;
          }
        } catch (_) {
          return p;
        }
      }
    }
    try {
      final whichCmd = Platform.isWindows ? 'where' : 'which';
      final res = Process.runSync(whichCmd, ['vify']);
      final out = res.stdout.toString().trim().split('\n').first.trim();
      if (res.exitCode == 0 && out.isNotEmpty && File(out).existsSync()) {
        final currentExe = File(Platform.resolvedExecutable).resolveSymbolicLinksSync().toLowerCase();
        if (File(out).resolveSymbolicLinksSync().toLowerCase() != currentExe) {
          return out;
        }
      }
    } catch (_) {}
    return Platform.isWindows ? 'vify-cli.exe' : 'vify';
  }

  Future<bool> isDaemonRunning() async {
    try {
      final res = await http.get(Uri.parse('$baseUrl/api/ping')).timeout(
        const Duration(milliseconds: 800),
      );
      return res.statusCode == 200;
    } catch (_) {
      return false;
    }
  }

  Future<bool> ensureDaemonRunning({bool isTun = true}) async {
    if (await isDaemonRunning()) {
      return true;
    }

    final binary = findVifyBinary();
    debugPrint('Starting Vify daemon using: $binary (TUN: $isTun)');

    try {
      if (Platform.isMacOS) {
        if (isTun) {
          // Create a small helper script that runs the daemon detached.
          // nohup fails inside AppleScript's "do shell script", so we use
          // bash -c with explicit stdin/stdout redirection and disown.
          final helperScript = '/tmp/vify-start-daemon.sh';
          File(helperScript).writeAsStringSync(
            '#!/bin/bash\n'
            '\'$binary\' daemon </dev/null >/tmp/vify-daemon.log 2>&1 &\n'
            'disown\n'
            'exit 0\n',
          );
          // Make executable
          await Process.run('chmod', ['+x', helperScript]);

          // Native macOS admin prompt runs the helper script
          final appleScript = 'do shell script "/tmp/vify-start-daemon.sh" with administrator privileges';
          final res = await Process.run('osascript', ['-e', appleScript]);
          debugPrint('osascript exit=${res.exitCode} stdout=${res.stdout} stderr=${res.stderr}');
          if (res.exitCode != 0) {
            debugPrint('macOS authorization denied or failed: ${res.stderr}');
            return false;
          }
        } else {
          // System Proxy mode does not require root
          await Process.start(binary, ['daemon'],
              mode: ProcessStartMode.detached);
        }
      } else if (Platform.isLinux) {
        if (isTun) {
          await Process.start('pkexec', [binary, 'daemon'],
              mode: ProcessStartMode.detached);
        } else {
          await Process.start(binary, ['daemon'],
              mode: ProcessStartMode.detached);
        }
      } else {
        await Process.start(binary, ['daemon'],
            mode: ProcessStartMode.detached);
      }
    } catch (e) {
      debugPrint('Failed to start daemon process: $e');
      return false;
    }

    // Wait up to 8 seconds for daemon port to open
    for (int i = 0; i < 40; i++) {
      await Future.delayed(const Duration(milliseconds: 200));
      if (await isDaemonRunning()) {
        debugPrint('Daemon is ready after ${(i + 1) * 200}ms');
        return true;
      }
    }
    // Read log for debug info
    try {
      final logContent = File('/tmp/vify-daemon.log').readAsStringSync();
      debugPrint('Daemon log: $logContent');
    } catch (_) {}
    debugPrint('Daemon failed to start within 8s');
    return false;
  }

  Future<Map<String, dynamic>?> getStatus() async {
    try {
      final res = await http.get(Uri.parse('$baseUrl/api/status')).timeout(
        const Duration(seconds: 2),
      );
      if (res.statusCode == 200) {
        return jsonDecode(res.body);
      }
    } catch (e) {
      debugPrint('Error getting daemon status: $e');
    }
    return null;
  }

  Future<bool> connect({
    String country = '',
    String protocol = '',
    String mode = 'tun',
    String rawUrl = '',
  }) async {
    try {
      final res = await http.post(
        Uri.parse('$baseUrl/api/connect'),
        headers: {'Content-Type': 'application/json'},
        body: jsonEncode({
          'country': country,
          'protocol': protocol,
          'mode': mode,
          'fast_pass': true,
          'batch_size': 60,
          'raw_url': rawUrl,
        }),
      );
      return res.statusCode == 200;
    } catch (e) {
      debugPrint('Error connecting via daemon: $e');
      return false;
    }
  }

  Future<bool> disconnect() async {
    try {
      final res = await http.post(Uri.parse('$baseUrl/api/disconnect'));
      return res.statusCode == 200;
    } catch (e) {
      debugPrint('Error disconnecting via daemon: $e');
      return false;
    }
  }

  Future<List<ProxyNodeModel>> getNodes() async {
    try {
      final res = await http.get(Uri.parse('$baseUrl/api/nodes')).timeout(
        const Duration(seconds: 4),
      );
      if (res.statusCode == 200) {
        final data = jsonDecode(res.body);
        if (data['nodes'] != null) {
          return (data['nodes'] as List)
              .map((n) => ProxyNodeModel.fromJson(n))
              .toList();
        }
      }
    } catch (e) {
      debugPrint('Error fetching nodes: $e');
    }
    return [];
  }
}
