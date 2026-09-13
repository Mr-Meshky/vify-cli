import 'package:flutter/material.dart';
import '../theme/vify_theme.dart';
import '../widgets/custom_title_bar.dart';
import '../widgets/metrics_panel.dart';
import '../widgets/power_button.dart';
import '../widgets/server_card.dart';

class MainScreen extends StatelessWidget {
  const MainScreen({super.key});

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
