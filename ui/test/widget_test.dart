import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:provider/provider.dart';
import 'package:vify/providers/vpn_provider.dart';
import 'package:vify/screens/main_screen.dart';
import 'package:vify/theme/vify_theme.dart';

void main() {
  testWidgets('Vify App initial smoke test', (WidgetTester tester) async {
    final vpnProvider = VpnProvider();

    await tester.pumpWidget(
      MultiProvider(
        providers: [
          ChangeNotifierProvider.value(value: vpnProvider),
        ],
        child: MaterialApp(
          theme: VifyTheme.darkTheme,
          home: const MainScreen(),
        ),
      ),
    );

    expect(find.text('VIFY'), findsOneWidget);
    expect(find.textContaining('Fast-Pass'), findsOneWidget);

    vpnProvider.dispose();
  });
}
