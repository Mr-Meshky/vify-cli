import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../providers/vpn_provider.dart';
import '../theme/vify_theme.dart';
import 'server_bottom_sheet.dart';

class ServerCard extends StatelessWidget {
  const ServerCard({super.key});

  @override
  Widget build(BuildContext context) {
    final vpn = context.watch<VpnProvider>();
    final node = vpn.currentNode;

    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 20),
      child: Material(
        color: Colors.transparent,
        child: InkWell(
          borderRadius: BorderRadius.circular(16),
          onTap: () {
            showModalBottomSheet(
              context: context,
              backgroundColor: Colors.transparent,
              isScrollControlled: true,
              builder: (ctx) => const ServerBottomSheet(),
            );
          },
          child: Container(
            padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 14),
            decoration: VifyTheme.glassCard(
              radius: 16,
              borderColor: vpn.isConnected
                  ? VifyTheme.neonEmerald.withAlpha(60)
                  : VifyTheme.borderGlass,
            ),
            child: Row(
              children: [
                // Flag Avatar
                Container(
                  width: 36,
                  height: 36,
                  decoration: BoxDecoration(
                    color: VifyTheme.bgSurface,
                    shape: BoxShape.circle,
                    border: Border.all(color: VifyTheme.borderGlass),
                  ),
                  child: Center(
                    child: Text(
                      node?.countryFlag ?? '⚡',
                      style: const TextStyle(fontSize: 18),
                    ),
                  ),
                ),
                const SizedBox(width: 12),

                // Server Details
                Expanded(
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(
                        node?.name ?? 'Fast-Pass Automatic Server',
                        style: const TextStyle(
                          fontSize: 13,
                          fontWeight: FontWeight.bold,
                          color: VifyTheme.textPrimary,
                        ),
                        maxLines: 1,
                        overflow: TextOverflow.ellipsis,
                      ),
                      const SizedBox(height: 2),
                      Text(
                        node != null
                            ? '${node.protocol.toUpperCase()} • ${node.latencyMs > 0 ? '${node.latencyMs}ms' : 'Optimized for speed'}'
                            : 'Auto-select lowest latency node',
                        style: const TextStyle(
                          fontSize: 11,
                          color: VifyTheme.textSecondary,
                        ),
                      ),
                    ],
                  ),
                ),

                const Icon(
                  Icons.arrow_forward_ios,
                  size: 14,
                  color: VifyTheme.textSecondary,
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }
}
