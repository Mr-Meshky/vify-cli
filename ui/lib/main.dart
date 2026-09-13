import 'dart:io';
import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'providers/vpn_provider.dart';
import 'screens/main_screen.dart';
import 'services/tray_service.dart';
import 'services/window_service.dart';
import 'theme/vify_theme.dart';

void main() async {
  WidgetsFlutterBinding.ensureInitialized();

  // Initialize Desktop Window (macOS / Windows / Linux)
  await WindowService.initialize();

  final vpnProvider = VpnProvider();

  // Initialize System Tray / Menu Bar
  await TrayService().initialize(
    onToggle: () => vpnProvider.toggleConnection(),
    quit: () => exit(0),
  );

  runApp(
    MultiProvider(
      providers: [
        ChangeNotifierProvider.value(value: vpnProvider),
      ],
      child: const VifyApp(),
    ),
  );
}

class VifyApp extends StatelessWidget {
  const VifyApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Vify',
      debugShowCheckedModeBanner: false,
      theme: VifyTheme.darkTheme,
      home: const MainScreen(),
    );
  }
}
