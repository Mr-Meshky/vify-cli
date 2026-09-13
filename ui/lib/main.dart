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

  final vpnProvider = VpnProvider();

  // Initialize Desktop Window and Tray only on Desktop
  if (Platform.isMacOS || Platform.isWindows || Platform.isLinux) {
    await WindowService.initialize();
    await TrayService().initialize(
      onToggle: () => vpnProvider.toggleConnection(),
      quit: () => exit(0),
    );
  }

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
