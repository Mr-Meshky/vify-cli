import 'dart:io';
import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:window_manager/window_manager.dart';
import '../providers/vpn_provider.dart';
import '../services/window_service.dart';
import '../theme/vify_theme.dart';

class CustomTitleBar extends StatelessWidget {
  const CustomTitleBar({super.key});

  @override
  Widget build(BuildContext context) {
    final vpn = context.watch<VpnProvider>();
    final isDesktop = Platform.isMacOS || Platform.isWindows || Platform.isLinux;
    final isMac = Platform.isMacOS;

    Widget branding = Row(
      mainAxisSize: MainAxisSize.min,
      children: [
        Container(
          width: 28,
          height: 28,
          decoration: BoxDecoration(
            gradient: const LinearGradient(
              colors: [VifyTheme.neonCyan, VifyTheme.neonEmerald],
              begin: Alignment.topLeft,
              end: Alignment.bottomRight,
            ),
            borderRadius: BorderRadius.circular(8),
            boxShadow: [
              BoxShadow(
                color: VifyTheme.neonCyan.withAlpha(80),
                blurRadius: 10,
                offset: const Offset(0, 2),
              ),
            ],
          ),
          child: const Center(
            child: Icon(
              Icons.bolt,
              size: 17,
              color: Colors.black,
            ),
          ),
        ),
        const SizedBox(width: 10),
        Text(
          'VIFY',
          style: Theme.of(context).textTheme.titleLarge?.copyWith(
                fontSize: 16,
                fontWeight: FontWeight.w900,
                letterSpacing: 2.2,
                color: VifyTheme.textPrimary,
              ),
        ),
        if (Platform.isAndroid) ...[
          const SizedBox(width: 8),
          Container(
            padding: const EdgeInsets.symmetric(horizontal: 6, vertical: 2),
            decoration: BoxDecoration(
              color: VifyTheme.neonEmerald.withAlpha(25),
              borderRadius: BorderRadius.circular(6),
              border: Border.all(color: VifyTheme.neonEmerald.withAlpha(80), width: 0.8),
            ),
            child: const Text(
              'ANDROID',
              style: TextStyle(
                fontSize: 9,
                fontWeight: FontWeight.w800,
                color: VifyTheme.neonEmerald,
                letterSpacing: 0.8,
              ),
            ),
          ),
        ],
      ],
    );

    return Container(
      height: 52,
      padding: EdgeInsets.only(
        left: isMac ? 84 : 16,
        right: 16,
      ),
      child: Row(
        children: [
          // Logo & Title (Draggable window area on Desktop)
          isDesktop ? DragToMoveArea(child: branding) : branding,

          // Draggable spacer across empty area on Desktop
          Expanded(
            child: isDesktop
                ? const DragToMoveArea(
                    child: SizedBox(
                      height: 52,
                      width: double.infinity,
                    ),
                  )
                : const SizedBox(),
          ),

          // Android: Quick Refresh Servers Button
          if (Platform.isAndroid) ...[
            IconButton(
              icon: const Icon(Icons.refresh, size: 20, color: VifyTheme.textSecondary),
              tooltip: 'Refresh Subscriptions',
              splashRadius: 18,
              onPressed: () => vpn.fetchNodes(),
            ),
            const SizedBox(width: 4),
          ],

          // Mode Switch (TUN Mode / System Proxy Mode)
          Material(
            color: Colors.transparent,
            child: InkWell(
              borderRadius: BorderRadius.circular(8),
              onTap: () => vpn.toggleTunMode(!vpn.isTunMode),
              child: Tooltip(
                message: vpn.isTunMode
                    ? 'TUN Mode (Full Device VPN). Click to switch to Proxy.'
                    : 'Proxy Mode (Local Port). Click to switch to TUN.',
                child: AnimatedContainer(
                  duration: const Duration(milliseconds: 200),
                  padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 5),
                  decoration: BoxDecoration(
                    color: vpn.isTunMode
                        ? VifyTheme.neonCyan.withAlpha(20)
                        : VifyTheme.neonPurple.withAlpha(20),
                    borderRadius: BorderRadius.circular(8),
                    border: Border.all(
                      color: vpn.isTunMode
                          ? VifyTheme.neonCyan.withAlpha(100)
                          : VifyTheme.neonPurple.withAlpha(100),
                      width: 1.2,
                    ),
                    boxShadow: [
                      BoxShadow(
                        color: (vpn.isTunMode ? VifyTheme.neonCyan : VifyTheme.neonPurple)
                            .withAlpha(25),
                        blurRadius: 8,
                      ),
                    ],
                  ),
                  child: Row(
                    mainAxisSize: MainAxisSize.min,
                    children: [
                      Icon(
                        vpn.isTunMode ? Icons.shield_outlined : Icons.public_rounded,
                        size: 13,
                        color: vpn.isTunMode
                            ? VifyTheme.neonCyan
                            : VifyTheme.neonPurple,
                      ),
                      const SizedBox(width: 5),
                      Text(
                        vpn.isTunMode ? 'TUN' : 'PROXY',
                        style: TextStyle(
                          fontSize: 11,
                          fontWeight: FontWeight.w700,
                          letterSpacing: 0.8,
                          color: vpn.isTunMode
                              ? VifyTheme.neonCyan
                              : VifyTheme.neonPurple,
                        ),
                      ),
                      const SizedBox(width: 4),
                      Icon(
                        Icons.swap_horiz,
                        size: 12,
                        color: (vpn.isTunMode ? VifyTheme.neonCyan : VifyTheme.neonPurple)
                            .withAlpha(180),
                      ),
                    ],
                  ),
                ),
              ),
            ),
          ),

          // Window Controls (Only on Windows / Linux)
          if (isDesktop && !isMac) ...[
            const SizedBox(width: 8),
            IconButton(
              icon: const Icon(Icons.remove, size: 18, color: VifyTheme.textSecondary),
              splashRadius: 16,
              tooltip: 'Minimize',
              onPressed: () => WindowService.minimize(),
            ),
            IconButton(
              icon: const Icon(Icons.close, size: 18, color: VifyTheme.textSecondary),
              splashRadius: 16,
              tooltip: 'Close to Tray',
              onPressed: () => WindowService.closeToTray(),
            ),
          ],
        ],
      ),
    );
  }
}
