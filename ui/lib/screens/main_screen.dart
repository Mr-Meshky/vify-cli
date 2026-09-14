import 'dart:async';
import 'package:flutter/material.dart';
import '../services/update_service.dart';
import '../theme/vify_theme.dart';
import '../widgets/custom_title_bar.dart';
import '../widgets/metrics_panel.dart';
import '../widgets/power_button.dart';
import '../widgets/server_card.dart';
import '../widgets/update_dialog.dart';

class MainScreen extends StatefulWidget {
  const MainScreen({super.key});

  @override
  State<MainScreen> createState() => _MainScreenState();
}

class _MainScreenState extends State<MainScreen> {
  Timer? _updateTimer;

  @override
  void initState() {
    super.initState();
    // Check for updates shortly after first frame
    WidgetsBinding.instance.addPostFrameCallback((_) {
      _updateTimer = Timer(const Duration(milliseconds: 600), () {
        if (mounted) {
          _checkUpdate();
        }
      });
    });
  }

  @override
  void dispose() {
    _updateTimer?.cancel();
    super.dispose();
  }

  Future<void> _checkUpdate() async {
    if (!mounted) return;
    final update = await UpdateService().checkForUpdate();
    if (update != null && update.isUpdateAvailable && mounted) {
      UpdateDialog.show(context, update);
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: VifyTheme.bgDark,
      body: SafeArea(
        child: Column(
          children: [
            // Top Custom Title Bar with Window Controls
            const CustomTitleBar(),

            // Scrollable Content Area for responsiveness across small and large screens
            Expanded(
              child: SingleChildScrollView(
                physics: const BouncingScrollPhysics(),
                child: Padding(
                  padding: const EdgeInsets.symmetric(vertical: 8),
                  child: Column(
                    children: const [
                      SizedBox(height: 16),

                      // Central Glowing One-Tap Power Button
                      Center(child: PowerButton()),

                      SizedBox(height: 28),

                      // Live Download & Upload Speed Metrics Panel
                      MetricsPanel(),

                      SizedBox(height: 18),

                      // Active Server Selector Card (Click to open Drawer/Sheet)
                      ServerCard(),

                      SizedBox(height: 18),
                    ],
                  ),
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }
}

