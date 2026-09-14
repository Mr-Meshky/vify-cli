import 'package:flutter/material.dart';
import 'package:flutter_animate/flutter_animate.dart';
import '../services/update_service.dart';
import '../theme/vify_theme.dart';

class UpdateDialog extends StatelessWidget {
  final UpdateInfo updateInfo;

  const UpdateDialog({
    super.key,
    required this.updateInfo,
  });

  static Future<void> show(BuildContext context, UpdateInfo info) {
    return showDialog<void>(
      context: context,
      barrierDismissible: true,
      barrierColor: Colors.black.withAlpha(180),
      builder: (context) => UpdateDialog(updateInfo: info),
    );
  }

  @override
  Widget build(BuildContext context) {
    return Dialog(
      backgroundColor: Colors.transparent,
      insetPadding: const EdgeInsets.symmetric(horizontal: 24, vertical: 32),
      child: Container(
        constraints: const BoxConstraints(maxWidth: 440),
        decoration: BoxDecoration(
          color: VifyTheme.bgCard,
          borderRadius: BorderRadius.circular(24),
          border: Border.all(color: VifyTheme.neonCyan.withAlpha(80), width: 1.5),
          boxShadow: [
            BoxShadow(
              color: VifyTheme.neonCyan.withAlpha(35),
              blurRadius: 32,
              spreadRadius: 2,
              offset: const Offset(0, 8),
            ),
          ],
        ),
        padding: const EdgeInsets.all(24),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            // Glowing Top Icon Header
            Center(
              child: Container(
                width: 64,
                height: 64,
                decoration: BoxDecoration(
                  shape: BoxShape.circle,
                  gradient: LinearGradient(
                    colors: [
                      VifyTheme.neonCyan.withAlpha(50),
                      VifyTheme.neonPurple.withAlpha(50),
                    ],
                    begin: Alignment.topLeft,
                    end: Alignment.bottomRight,
                  ),
                  border: Border.all(color: VifyTheme.neonCyan, width: 1.5),
                  boxShadow: [
                    BoxShadow(
                      color: VifyTheme.neonCyan.withAlpha(60),
                      blurRadius: 16,
                    ),
                  ],
                ),
                child: const Icon(
                  Icons.rocket_launch_rounded,
                  color: VifyTheme.neonCyan,
                  size: 32,
                ),
              ),
            ).animate().scale(duration: 350.ms, curve: Curves.easeOutBack),

            const SizedBox(height: 18),

            // Title
            const Center(
              child: Text(
                'نسخه جدید در دسترس است!',
                textAlign: TextAlign.center,
                style: TextStyle(
                  color: VifyTheme.textPrimary,
                  fontSize: 20,
                  fontWeight: FontWeight.bold,
                  letterSpacing: -0.5,
                ),
              ),
            ),

            const SizedBox(height: 12),

            // Version Comparison Badge
            Center(
              child: Container(
                padding: const EdgeInsets.symmetric(horizontal: 14, vertical: 6),
                decoration: BoxDecoration(
                  color: VifyTheme.bgSurface,
                  borderRadius: BorderRadius.circular(20),
                  border: Border.all(color: VifyTheme.borderGlass),
                ),
                child: Row(
                  mainAxisSize: MainAxisSize.min,
                  children: [
                    Text(
                      'v${updateInfo.currentVersion.replaceAll(RegExp(r'^v'), '')}',
                      style: const TextStyle(
                        color: VifyTheme.textMuted,
                        fontSize: 13,
                        fontWeight: FontWeight.w600,
                      ),
                    ),
                    const Padding(
                      padding: EdgeInsets.symmetric(horizontal: 8),
                      child: Icon(Icons.arrow_forward_rounded, size: 14, color: VifyTheme.neonEmerald),
                    ),
                    Text(
                      updateInfo.latestVersion,
                      style: const TextStyle(
                        color: VifyTheme.neonEmerald,
                        fontSize: 14,
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                  ],
                ),
              ),
            ),

            const SizedBox(height: 18),

            // Changelog Box (if any)
            if (updateInfo.changelog.isNotEmpty) ...[
              const Text(
                'تغییرات این نسخه:',
                style: TextStyle(
                  color: VifyTheme.textSecondary,
                  fontSize: 13,
                  fontWeight: FontWeight.w600,
                ),
              ),
              const SizedBox(height: 8),
              Flexible(
                child: Container(
                  constraints: const BoxConstraints(maxHeight: 160),
                  padding: const EdgeInsets.all(12),
                  decoration: BoxDecoration(
                    color: VifyTheme.bgDark.withAlpha(180),
                    borderRadius: BorderRadius.circular(12),
                    border: Border.all(color: VifyTheme.borderGlass),
                  ),
                  child: SingleChildScrollView(
                    physics: const BouncingScrollPhysics(),
                    child: Text(
                      updateInfo.changelog,
                      style: const TextStyle(
                        color: VifyTheme.textSecondary,
                        fontSize: 12.5,
                        height: 1.5,
                      ),
                    ),
                  ),
                ),
              ),
              const SizedBox(height: 20),
            ] else ...[
              const SizedBox(height: 8),
            ],

            // Action Buttons
            Row(
              children: [
                // "Later" Button
                Expanded(
                  child: TextButton(
                    onPressed: () async {
                      await UpdateService().dismissVersion(updateInfo.latestVersion);
                      if (context.mounted) {
                        Navigator.of(context).pop();
                      }
                    },
                    style: TextButton.styleFrom(
                      padding: const EdgeInsets.symmetric(vertical: 14),
                      shape: RoundedRectangleBorder(
                        borderRadius: BorderRadius.circular(12),
                      ),
                    ),
                    child: const Text(
                      'بعداً',
                      style: TextStyle(
                        color: VifyTheme.textMuted,
                        fontSize: 14,
                        fontWeight: FontWeight.w600,
                      ),
                    ),
                  ),
                ),

                const SizedBox(width: 12),

                // "Download & Update" Button
                Expanded(
                  flex: 2,
                  child: Container(
                    decoration: BoxDecoration(
                      borderRadius: BorderRadius.circular(12),
                      gradient: const LinearGradient(
                        colors: [VifyTheme.neonCyan, VifyTheme.neonBlue],
                      ),
                      boxShadow: [
                        BoxShadow(
                          color: VifyTheme.neonCyan.withAlpha(70),
                          blurRadius: 14,
                          offset: const Offset(0, 4),
                        ),
                      ],
                    ),
                    child: ElevatedButton(
                      onPressed: () async {
                        final url = updateInfo.downloadUrl;
                        await UpdateService.launchUpdateUrl(url);
                        if (context.mounted) {
                          Navigator.of(context).pop();
                        }
                      },
                      style: ElevatedButton.styleFrom(
                        backgroundColor: Colors.transparent,
                        shadowColor: Colors.transparent,
                        padding: const EdgeInsets.symmetric(vertical: 14),
                        shape: RoundedRectangleBorder(
                          borderRadius: BorderRadius.circular(12),
                        ),
                      ),
                      child: const Row(
                        mainAxisAlignment: MainAxisAlignment.center,
                        children: [
                          Icon(Icons.download_rounded, color: Colors.black, size: 18),
                          SizedBox(width: 6),
                          Text(
                            'دانلود و به‌روزرسانی',
                            style: TextStyle(
                              color: Colors.black,
                              fontWeight: FontWeight.bold,
                              fontSize: 13.5,
                            ),
                          ),
                        ],
                      ),
                    ),
                  ),
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }
}
