import 'dart:io';
import 'package:flutter/material.dart';
import 'package:tray_manager/tray_manager.dart';
import 'window_service.dart';

class TrayService with TrayListener {
  static final TrayService _instance = TrayService._internal();
  factory TrayService() => _instance;
  TrayService._internal();

  VoidCallback? onToggleConnect;
  VoidCallback? onQuit;

  Future<void> initialize({VoidCallback? onToggle, VoidCallback? quit}) async {
    if (Platform.isMacOS || Platform.isWindows || Platform.isLinux) {
      onToggleConnect = onToggle;
      onQuit = quit;
      trayManager.addListener(this);

      String iconPath = Platform.isWindows
          ? 'assets/app_icon.ico'
          : 'assets/app_icon.png';

      // Fallback or set tray icon
      try {
        await trayManager.setIcon(iconPath);
      } catch (_) {}

      await updateContextMenu(isConnected: false);
    }
  }

  Future<void> updateContextMenu({required bool isConnected}) async {
    if (Platform.isMacOS || Platform.isWindows || Platform.isLinux) {
      Menu menu = Menu(
        items: [
          MenuItem(
            key: 'show_window',
            label: 'Open Vify',
          ),
          MenuItem.separator(),
          MenuItem(
            key: 'toggle_vpn',
            label: isConnected ? '⚡ Disconnect VPN' : '⚡ Connect VPN (Fast-Pass)',
          ),
          MenuItem.separator(),
          MenuItem(
            key: 'quit_app',
            label: 'Quit Vify',
          ),
        ],
      );
      await trayManager.setContextMenu(menu);
    }
  }

  @override
  void onTrayIconMouseDown() {
    WindowService.showWindow();
  }

  @override
  void onTrayIconRightMouseDown() {
    trayManager.popUpContextMenu();
  }

  @override
  void onTrayMenuItemClick(MenuItem menuItem) {
    if (menuItem.key == 'show_window') {
      WindowService.showWindow();
    } else if (menuItem.key == 'toggle_vpn') {
      onToggleConnect?.call();
    } else if (menuItem.key == 'quit_app') {
      if (onQuit != null) {
        onQuit!();
      } else {
        exit(0);
      }
    }
  }
}
