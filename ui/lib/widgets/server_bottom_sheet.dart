import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:provider/provider.dart';
import '../providers/vpn_provider.dart';
import '../theme/vify_theme.dart';

class ServerBottomSheet extends StatefulWidget {
  const ServerBottomSheet({super.key});

  @override
  State<ServerBottomSheet> createState() => _ServerBottomSheetState();
}

class _ServerBottomSheetState extends State<ServerBottomSheet> {
  String _searchQuery = '';
  final TextEditingController _searchCtrl = TextEditingController();

  void _showAddConfigDialog() {
    final TextEditingController urlCtrl = TextEditingController();
    showDialog(
      context: context,
      builder: (ctx) => AlertDialog(
        backgroundColor: VifyTheme.bgSurface,
        shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(16)),
        title: const Text(
          'Import Custom Config',
          style: TextStyle(color: VifyTheme.textPrimary, fontSize: 16, fontWeight: FontWeight.bold),
        ),
        content: Column(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            const Text(
              'Paste a VLESS, VMess, Trojan, or Shadowsocks URL:',
              style: TextStyle(color: VifyTheme.textSecondary, fontSize: 12),
            ),
            const SizedBox(height: 12),
            TextField(
              controller: urlCtrl,
              maxLines: 3,
              style: const TextStyle(color: VifyTheme.textPrimary, fontSize: 12),
              decoration: InputDecoration(
                hintText: 'vless://... or vmess://...',
                hintStyle: const TextStyle(color: VifyTheme.textMuted, fontSize: 12),
                filled: true,
                fillColor: VifyTheme.bgCard,
                border: OutlineInputBorder(
                  borderRadius: BorderRadius.circular(10),
                  borderSide: const BorderSide(color: VifyTheme.borderGlass),
                ),
              ),
            ),
            const SizedBox(height: 8),
            TextButton.icon(
              icon: const Icon(Icons.content_paste, size: 16, color: VifyTheme.neonCyan),
              label: const Text('Paste from Clipboard', style: TextStyle(color: VifyTheme.neonCyan, fontSize: 12)),
              onPressed: () async {
                final data = await Clipboard.getData(Clipboard.kTextPlain);
                if (data != null && data.text != null) {
                  urlCtrl.text = data.text!.trim();
                }
              },
            ),
          ],
        ),
        actions: [
          TextButton(
            child: const Text('Cancel', style: TextStyle(color: VifyTheme.textMuted)),
            onPressed: () => Navigator.pop(ctx),
          ),
          ElevatedButton(
            style: ElevatedButton.styleFrom(
              backgroundColor: VifyTheme.neonCyan,
              foregroundColor: Colors.black,
              shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(8)),
            ),
            child: const Text('Connect Now', style: TextStyle(fontWeight: FontWeight.bold)),
            onPressed: () {
              final link = urlCtrl.text.trim();
              if (link.isNotEmpty) {
                Navigator.pop(ctx);
                Navigator.pop(context);
                context.read<VpnProvider>().connect(rawUrl: link);
              }
            },
          ),
        ],
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    final vpn = context.watch<VpnProvider>();
    final nodes = vpn.nodes.where((n) {
      if (_searchQuery.isEmpty) return true;
      return n.name.toLowerCase().contains(_searchQuery.toLowerCase()) ||
          n.country.toLowerCase().contains(_searchQuery.toLowerCase()) ||
          n.protocol.toLowerCase().contains(_searchQuery.toLowerCase());
    }).toList();

    return Container(
      height: MediaQuery.of(context).size.height * 0.75,
      decoration: const BoxDecoration(
        color: VifyTheme.bgSurface,
        borderRadius: BorderRadius.vertical(top: Radius.circular(24)),
        border: Border(top: BorderSide(color: VifyTheme.borderGlass, width: 1.5)),
      ),
      child: Column(
        children: [
          // Drag handle
          Container(
            margin: const EdgeInsets.symmetric(vertical: 12),
            width: 40,
            height: 4,
            decoration: BoxDecoration(
              color: VifyTheme.textMuted,
              borderRadius: BorderRadius.circular(2),
            ),
          ),

          // Header
          Padding(
            padding: const EdgeInsets.symmetric(horizontal: 20),
            child: Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                Text(
                  'Locations & Servers (${vpn.nodes.length})',
                  style: Theme.of(context).textTheme.titleLarge?.copyWith(
                        fontSize: 16,
                        fontWeight: FontWeight.bold,
                      ),
                ),
                Row(
                  children: [
                    IconButton(
                      icon: const Icon(Icons.add_link, size: 22, color: VifyTheme.neonCyan),
                      onPressed: _showAddConfigDialog,
                      tooltip: 'Import Custom Config',
                    ),
                    IconButton(
                      icon: const Icon(Icons.refresh, size: 20, color: VifyTheme.textSecondary),
                      onPressed: () => vpn.fetchNodes(),
                      tooltip: 'Refresh Servers',
                    ),
                  ],
                ),
              ],
            ),
          ),

          // Search Bar
          Padding(
            padding: const EdgeInsets.symmetric(horizontal: 20, vertical: 8),
            child: TextField(
              controller: _searchCtrl,
              onChanged: (val) => setState(() => _searchQuery = val),
              style: const TextStyle(color: VifyTheme.textPrimary, fontSize: 13),
              decoration: InputDecoration(
                hintText: 'Search by country, name, or protocol...',
                hintStyle: const TextStyle(color: VifyTheme.textMuted, fontSize: 13),
                prefixIcon: const Icon(Icons.search, size: 18, color: VifyTheme.textSecondary),
                filled: true,
                fillColor: VifyTheme.bgCard,
                contentPadding: const EdgeInsets.symmetric(vertical: 0, horizontal: 16),
                border: OutlineInputBorder(
                  borderRadius: BorderRadius.circular(12),
                  borderSide: const BorderSide(color: VifyTheme.borderGlass),
                ),
                enabledBorder: OutlineInputBorder(
                  borderRadius: BorderRadius.circular(12),
                  borderSide: const BorderSide(color: VifyTheme.borderGlass),
                ),
                focusedBorder: OutlineInputBorder(
                  borderRadius: BorderRadius.circular(12),
                  borderSide: const BorderSide(color: VifyTheme.neonCyan),
                ),
              ),
            ),
          ),

          // Fast-Pass Auto Option
          Padding(
            padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 4),
            child: Material(
              color: Colors.transparent,
              child: InkWell(
                borderRadius: BorderRadius.circular(12),
                onTap: () {
                  Navigator.pop(context);
                  vpn.connect();
                },
                child: Container(
                  padding: const EdgeInsets.symmetric(horizontal: 14, vertical: 10),
                  decoration: BoxDecoration(
                    color: VifyTheme.neonEmerald.withAlpha(15),
                    borderRadius: BorderRadius.circular(12),
                    border: Border.all(color: VifyTheme.neonEmerald.withAlpha(60)),
                  ),
                  child: Row(
                    children: const [
                      Text('⚡', style: TextStyle(fontSize: 18)),
                      SizedBox(width: 10),
                      Expanded(
                        child: Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            Text(
                              'Fast-Pass Automatic Selection',
                              style: TextStyle(
                                fontSize: 13,
                                fontWeight: FontWeight.bold,
                                color: VifyTheme.neonEmerald,
                              ),
                            ),
                            Text(
                              'Automatically benchmark & connect to fastest server',
                              style: TextStyle(fontSize: 11, color: VifyTheme.textSecondary),
                            ),
                          ],
                        ),
                      ),
                      Icon(Icons.bolt, color: VifyTheme.neonEmerald, size: 18),
                    ],
                  ),
                ),
              ),
            ),
          ),

          const SizedBox(height: 4),

          // Server List
          Expanded(
            child: nodes.isEmpty
                ? const Center(
                    child: Text(
                      'No servers matched',
                      style: TextStyle(color: VifyTheme.textMuted),
                    ),
                  )
                : ListView.builder(
                    itemCount: nodes.length,
                    padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
                    itemBuilder: (ctx, index) {
                      final node = nodes[index];
                      final isSelected = vpn.currentNode?.server == node.server &&
                          vpn.currentNode?.port == node.port;

                      return Container(
                        margin: const EdgeInsets.only(bottom: 8),
                        decoration: BoxDecoration(
                          color: isSelected
                              ? VifyTheme.neonCyan.withAlpha(20)
                              : VifyTheme.bgCard,
                          borderRadius: BorderRadius.circular(12),
                          border: Border.all(
                            color: isSelected
                                ? VifyTheme.neonCyan
                                : VifyTheme.borderGlass,
                          ),
                        ),
                        child: ListTile(
                          dense: true,
                          leading: Text(
                            node.countryFlag,
                            style: const TextStyle(fontSize: 22),
                          ),
                          title: Text(
                            node.name,
                            style: TextStyle(
                              fontSize: 13,
                              fontWeight: isSelected ? FontWeight.bold : FontWeight.w500,
                              color: isSelected
                                  ? VifyTheme.neonCyan
                                  : VifyTheme.textPrimary,
                            ),
                            maxLines: 1,
                            overflow: TextOverflow.ellipsis,
                          ),
                          subtitle: Text(
                            '${node.protocol.toUpperCase()} • ${node.server}',
                            style: const TextStyle(fontSize: 11, color: VifyTheme.textMuted),
                          ),
                          trailing: node.latencyMs > 0
                              ? Container(
                                  padding: const EdgeInsets.symmetric(
                                      horizontal: 8, vertical: 4),
                                  decoration: BoxDecoration(
                                    color: VifyTheme.bgSurface,
                                    borderRadius: BorderRadius.circular(8),
                                    border: Border.all(
                                      color: node.latencyMs < 200
                                          ? VifyTheme.neonEmerald.withAlpha(60)
                                          : VifyTheme.neonAmber.withAlpha(60),
                                    ),
                                  ),
                                  child: Text(
                                    '${node.latencyMs} ms',
                                    style: TextStyle(
                                      fontSize: 11,
                                      fontWeight: FontWeight.bold,
                                      color: node.latencyMs < 200
                                          ? VifyTheme.neonEmerald
                                          : VifyTheme.neonAmber,
                                    ),
                                  ),
                                )
                              : const Icon(Icons.arrow_forward_ios, size: 12, color: VifyTheme.textMuted),
                          onTap: () {
                            Navigator.pop(context);
                            vpn.connect(preselectedNode: node);
                          },
                        ),
                      );
                    },
                  ),
          ),
        ],
      ),
    );
  }
}
