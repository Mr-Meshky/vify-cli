import 'package:flutter/material.dart';
import 'package:flutter_animate/flutter_animate.dart';
import 'package:provider/provider.dart';
import '../providers/vpn_provider.dart';
import '../theme/vify_theme.dart';

class PowerButton extends StatelessWidget {
  const PowerButton({super.key});

  @override
  Widget build(BuildContext context) {
    final vpn = context.watch<VpnProvider>();
    final isConnected = vpn.isConnected;
    final isConnecting = vpn.isConnecting;
    final hasError = vpn.hasError;

    final glowColor = isConnected
        ? VifyTheme.neonEmerald
        : hasError
            ? VifyTheme.neonRed
            : isConnecting
                ? VifyTheme.neonCyan
                : VifyTheme.neonBlue.withAlpha(80);

    return Column(
      mainAxisSize: MainAxisSize.min,
      children: [
        // Status Badge & Latency Pill
        Container(
          padding: const EdgeInsets.symmetric(horizontal: 14, vertical: 6),
          decoration: BoxDecoration(
            color: VifyTheme.bgCard,
            borderRadius: BorderRadius.circular(20),
            border: Border.all(
              color: isConnected
                  ? VifyTheme.neonEmerald.withAlpha(100)
                  : hasError
                      ? VifyTheme.neonRed.withAlpha(100)
                      : VifyTheme.borderGlass,
            ),
          ),
          child: Row(
            mainAxisSize: MainAxisSize.min,
            children: [
              Container(
                width: 8,
                height: 8,
                decoration: BoxDecoration(
                  shape: BoxShape.circle,
                  color: isConnected
                      ? VifyTheme.neonEmerald
                      : hasError
                          ? VifyTheme.neonRed
                          : isConnecting
                              ? VifyTheme.neonAmber
                              : VifyTheme.textMuted,
                  boxShadow: isConnected || isConnecting || hasError
                      ? [
                          BoxShadow(
                            color: glowColor.withAlpha(180),
                            blurRadius: 6,
                          ),
                        ]
                      : [],
                ),
              ),
              const SizedBox(width: 8),
              Text(
                vpn.statusMessage.toUpperCase(),
                style: TextStyle(
                  fontSize: 11,
                  fontWeight: FontWeight.bold,
                  letterSpacing: 1.2,
                  color: isConnected
                      ? VifyTheme.neonEmerald
                      : hasError
                          ? VifyTheme.neonRed
                          : isConnecting
                              ? VifyTheme.neonCyan
                              : VifyTheme.textSecondary,
                ),
              ),
              if (isConnected && vpn.latencyMs > 0) ...[
                const SizedBox(width: 8),
                Container(
                  width: 1,
                  height: 12,
                  color: VifyTheme.borderGlass,
                ),
                const SizedBox(width: 8),
                Text(
                  '${vpn.latencyMs} ms',
                  style: const TextStyle(
                    fontSize: 11,
                    fontWeight: FontWeight.w600,
                    color: VifyTheme.neonCyan,
                  ),
                ),
              ],
            ],
          ),
        ),

        const SizedBox(height: 32),

        // Central Multi-Ring Power Button
        GestureDetector(
          onTap: () => vpn.toggleConnection(),
          child: SizedBox(
            width: 200,
            height: 200,
            child: Stack(
              alignment: Alignment.center,
              children: [
                // Outer Pulse Ring 1
                if (isConnected || isConnecting)
                  Container(
                    width: 190,
                    height: 190,
                    decoration: BoxDecoration(
                      shape: BoxShape.circle,
                      border: Border.all(
                        color: glowColor.withAlpha(40),
                        width: 1.5,
                      ),
                    ),
                  )
                      .animate(onPlay: (c) => c.repeat(reverse: true))
                      .scale(begin: const Offset(0.92, 0.92), end: const Offset(1.05, 1.05), duration: 1800.ms),

                // Outer Pulse Ring 2
                if (isConnected || isConnecting)
                  Container(
                    width: 160,
                    height: 160,
                    decoration: BoxDecoration(
                      shape: BoxShape.circle,
                      border: Border.all(
                        color: glowColor.withAlpha(70),
                        width: 2,
                      ),
                    ),
                  )
                      .animate(onPlay: (c) => c.repeat(reverse: true))
                      .scale(begin: const Offset(0.96, 0.96), end: const Offset(1.03, 1.03), duration: 1400.ms),

                // Main Central Button
                AnimatedContainer(
                  duration: const Duration(milliseconds: 350),
                  width: 130,
                  height: 130,
                  decoration: BoxDecoration(
                    shape: BoxShape.circle,
                    gradient: LinearGradient(
                      begin: Alignment.topLeft,
                      end: Alignment.bottomRight,
                      colors: isConnected
                          ? [
                              const Color(0xFF0D2818),
                              const Color(0xFF041F14),
                            ]
                          : hasError
                              ? [
                                  const Color(0xFF2B0A12),
                                  const Color(0xFF19060B),
                                ]
                              : [
                                  VifyTheme.bgSurface,
                                  VifyTheme.bgDark,
                                ],
                    ),
                    border: Border.all(
                      color: isConnected
                          ? VifyTheme.neonEmerald
                          : hasError
                              ? VifyTheme.neonRed
                              : isConnecting
                                  ? VifyTheme.neonCyan
                                  : VifyTheme.borderGlass,
                      width: isConnected || hasError ? 3.0 : 2.0,
                    ),
                    boxShadow: [
                      BoxShadow(
                        color: isConnected
                            ? VifyTheme.neonEmerald.withAlpha(100)
                            : hasError
                                ? VifyTheme.neonRed.withAlpha(90)
                                : isConnecting
                                    ? VifyTheme.neonCyan.withAlpha(80)
                                    : Colors.black.withAlpha(80),
                        blurRadius: isConnected || isConnecting || hasError ? 32 : 16,
                        spreadRadius: isConnected ? 2 : 0,
                      ),
                    ],
                  ),
                  child: Center(
                    child: isConnecting
                        ? const SizedBox(
                            width: 38,
                            height: 38,
                            child: CircularProgressIndicator(
                              strokeWidth: 3,
                              valueColor: AlwaysStoppedAnimation<Color>(
                                VifyTheme.neonCyan,
                              ),
                            ),
                          )
                        : Icon(
                            hasError ? Icons.refresh : Icons.power_settings_new,
                            size: 48,
                            color: isConnected
                                ? VifyTheme.neonEmerald
                                : hasError
                                    ? VifyTheme.neonRed
                                    : VifyTheme.textSecondary,
                          ),
                  ),
                ),
              ],
            ),
          ),
        ),

        // Error message details (if any)
        if (hasError && vpn.lastError.isNotEmpty) ...[
          const SizedBox(height: 16),
          Padding(
            padding: const EdgeInsets.symmetric(horizontal: 24),
            child: Container(
              padding: const EdgeInsets.all(10),
              decoration: BoxDecoration(
                color: VifyTheme.neonRed.withAlpha(15),
                borderRadius: BorderRadius.circular(10),
                border: Border.all(color: VifyTheme.neonRed.withAlpha(50)),
              ),
              child: Text(
                vpn.lastError,
                textAlign: TextAlign.center,
                style: const TextStyle(
                  fontSize: 11,
                  color: VifyTheme.neonRed,
                ),
              ),
            ),
          ),
        ],
      ],
    );
  }
}
