import 'dart:io';
import 'package:flutter/material.dart';
import 'package:window_manager/window_manager.dart';

class WindowService with WindowListener {
  static final WindowService _instance = WindowService._internal();
  factory WindowService() => _instance;
  WindowService._internal();

  static Future<void> initialize() async {
    if (Platform.isMacOS || Platform.isWindows || Platform.isLinux) {
      await windowManager.ensureInitialized();
      windowManager.addListener(_instance);

      WindowOptions windowOptions = const WindowOptions(
        size: Size(380, 680),
        minimumSize: Size(360, 620),
        maximumSize: Size(440, 780),
        center: true,
        backgroundColor: Colors.transparent,
        skipTaskbar: false,
        titleBarStyle: TitleBarStyle.hidden,
        title: 'Vify',
      );

      await windowManager.waitUntilReadyToShow(windowOptions, () async {
        await windowManager.setPreventClose(true);
        await windowManager.show();
        await windowManager.focus();
      });
    }
  }

  @override
  void onWindowClose() async {
    bool isPreventClose = await windowManager.isPreventClose();
    if (isPreventClose) {
      await windowManager.hide();
    }
  }

  static Future<void> minimize() async {
    if (Platform.isMacOS || Platform.isWindows || Platform.isLinux) {
      await windowManager.minimize();
    }
  }

  static Future<void> closeToTray() async {
    if (Platform.isMacOS || Platform.isWindows || Platform.isLinux) {
      await windowManager.hide();
    }
  }

  static Future<void> showWindow() async {
    if (Platform.isMacOS || Platform.isWindows || Platform.isLinux) {
      await windowManager.show();
      await windowManager.focus();
    }
  }
}
