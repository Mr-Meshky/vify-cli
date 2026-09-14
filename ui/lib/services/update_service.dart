import 'dart:convert';
import 'dart:io';
import 'package:flutter/foundation.dart';
import 'package:http/http.dart' as http;
import 'package:package_info_plus/package_info_plus.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:url_launcher/url_launcher.dart';

class UpdateInfo {
  final String currentVersion;
  final String latestVersion;
  final String title;
  final String changelog;
  final String htmlUrl;
  final String? directDownloadUrl;
  final bool isUpdateAvailable;

  const UpdateInfo({
    required this.currentVersion,
    required this.latestVersion,
    required this.title,
    required this.changelog,
    required this.htmlUrl,
    this.directDownloadUrl,
    required this.isUpdateAvailable,
  });

  String get downloadUrl => directDownloadUrl ?? htmlUrl;
}

class UpdateService {
  static const String repoOwner = 'Mr-Meshky';
  static const String repoName = 'vify-cli';

  static const String _keyDismissedVersion = 'vify_dismissed_update_version';
  static const String _keyLastDismissedTime = 'vify_dismissed_update_time';

  /// Compares two semver strings like "v1.2.0" and "1.2.3".
  /// Returns 1 if v1 > v2, -1 if v1 < v2, 0 if equal.
  static int compareSemVer(String v1, String v2) {
    List<int> parse(String v) {
      String clean = v.trim();
      if (clean.toLowerCase().startsWith('v')) {
        clean = clean.substring(1);
      }
      final int hyphenIdx = clean.indexOf('-');
      if (hyphenIdx != -1) {
        clean = clean.substring(0, hyphenIdx);
      }
      final int plusIdx = clean.indexOf('+');
      if (plusIdx != -1) {
        clean = clean.substring(0, plusIdx);
      }
      final parts = clean.split('.');
      final res = <int>[0, 0, 0];
      for (int i = 0; i < parts.length && i < 3; i++) {
        res[i] = int.tryParse(parts[i]) ?? 0;
      }
      return res;
    }

    final p1 = parse(v1);
    final p2 = parse(v2);

    for (int i = 0; i < 3; i++) {
      if (p1[i] > p2[i]) return 1;
      if (p1[i] < p2[i]) return -1;
    }
    return 0;
  }

  /// Checks GitHub releases for a new version.
  /// If [force] is true, ignores user dismissal and checks immediately.
  Future<UpdateInfo?> checkForUpdate({bool force = false}) async {
    try {
      // 1. Get current installed version
      String currentVersion = '1.2.3';
      try {
        final packageInfo = await PackageInfo.fromPlatform();
        if (packageInfo.version.isNotEmpty) {
          currentVersion = packageInfo.version;
        }
      } catch (e) {
        debugPrint('UpdateService: Could not read PackageInfo, using fallback: $e');
      }

      // 2. Query GitHub Releases API
      final uri = Uri.parse('https://api.github.com/repos/$repoOwner/$repoName/releases/latest');
      final response = await http.get(
        uri,
        headers: {
          'Accept': 'application/vnd.github.v3+json',
          'User-Agent': 'vify-app-updater',
        },
      ).timeout(const Duration(seconds: 8));

      if (response.statusCode != 200) {
        debugPrint('UpdateService: GitHub API returned status ${response.statusCode}');
        return null;
      }

      final Map<String, dynamic> data = jsonDecode(response.body) as Map<String, dynamic>;
      final String latestTag = (data['tag_name'] as String? ?? '').trim();
      final String releaseName = (data['name'] as String? ?? latestTag);
      final String changelog = (data['body'] as String? ?? '').trim();
      final String htmlUrl = (data['html_url'] as String? ?? 'https://github.com/$repoOwner/$repoName/releases/latest');

      final bool isNewer = compareSemVer(latestTag, currentVersion) > 0;

      // 3. Find matching platform download asset if available
      String? directDownloadUrl;
      final assets = data['assets'] as List<dynamic>? ?? [];
      final String currentPlatform = _currentPlatformKeyword();

      for (final asset in assets) {
        if (asset is Map<String, dynamic>) {
          final assetName = (asset['name'] as String? ?? '').toLowerCase();
          final downloadUrl = asset['browser_download_url'] as String?;
          if (downloadUrl != null && assetName.contains(currentPlatform)) {
            directDownloadUrl = downloadUrl;
            break;
          }
        }
      }

      // Fallback for Android APK
      if (directDownloadUrl == null && Platform.isAndroid) {
        for (final asset in assets) {
          if (asset is Map<String, dynamic>) {
            final assetName = (asset['name'] as String? ?? '').toLowerCase();
            if (assetName.endsWith('.apk')) {
              directDownloadUrl = asset['browser_download_url'] as String?;
              break;
            }
          }
        }
      }

      // 4. Check dismissal rules if not forced
      if (!force && isNewer) {
        final prefs = await SharedPreferences.getInstance();
        final dismissedVersion = prefs.getString(_keyDismissedVersion);
        final dismissedTime = prefs.getInt(_keyLastDismissedTime) ?? 0;
        final now = DateTime.now().millisecondsSinceEpoch;

        // If dismissed within the last 24 hours for the same version, skip popup
        if (dismissedVersion == latestTag && (now - dismissedTime) < 24 * 60 * 60 * 1000) {
          debugPrint('UpdateService: Version $latestTag was recently dismissed by user');
          return null;
        }
      }

      return UpdateInfo(
        currentVersion: currentVersion,
        latestVersion: latestTag,
        title: releaseName,
        changelog: changelog,
        htmlUrl: htmlUrl,
        directDownloadUrl: directDownloadUrl,
        isUpdateAvailable: isNewer,
      );
    } catch (e) {
      debugPrint('UpdateService: Error checking for update: $e');
      return null;
    }
  }

  /// Mark the version as dismissed so it won't prompt the user again for 24h
  Future<void> dismissVersion(String version) async {
    try {
      final prefs = await SharedPreferences.getInstance();
      await prefs.setString(_keyDismissedVersion, version);
      await prefs.setInt(_keyLastDismissedTime, DateTime.now().millisecondsSinceEpoch);
    } catch (e) {
      debugPrint('UpdateService: Failed to persist dismissal: $e');
    }
  }

  /// Open the download URL in the system browser or default installer
  static Future<bool> launchUpdateUrl(String url) async {
    try {
      final uri = Uri.parse(url);
      return await launchUrl(uri, mode: LaunchMode.externalApplication);
    } catch (e) {
      debugPrint('UpdateService: Failed to launch update URL $url: $e');
      return false;
    }
  }

  static String _currentPlatformKeyword() {
    if (Platform.isAndroid) return 'android';
    if (Platform.isWindows) return 'windows';
    if (Platform.isMacOS) return 'darwin';
    if (Platform.isLinux) return 'linux';
    if (Platform.isIOS) return 'ios';
    return '';
  }
}
