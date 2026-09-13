import 'package:flutter/material.dart';
import 'package:google_fonts/google_fonts.dart';

class VifyTheme {
  // Color Palette
  static const Color bgDark = Color(0xFF090C13);
  static const Color bgSurface = Color(0xFF101522);
  static const Color bgCard = Color(0xFF151C2C);
  static const Color bgGlass = Color(0xCC151C2C);
  static const Color borderGlass = Color(0x332A3854);
  
  // Neon Accents
  static const Color neonCyan = Color(0xFF00F5D4);
  static const Color neonEmerald = Color(0xFF05FFA1);
  static const Color neonBlue = Color(0xFF00BBF9);
  static const Color neonPurple = Color(0xFF9B5DE5);
  static const Color neonAmber = Color(0xFFFFB703);
  static const Color neonRed = Color(0xFFFF0054);
  
  // Text Colors
  static const Color textPrimary = Color(0xFFFFFFFF);
  static const Color textSecondary = Color(0xFF8E9BB0);
  static const Color textMuted = Color(0xFF5A667A);

  static ThemeData get darkTheme {
    return ThemeData(
      brightness: Brightness.dark,
      scaffoldBackgroundColor: bgDark,
      primaryColor: neonCyan,
      textTheme: GoogleFonts.outfitTextTheme(ThemeData.dark().textTheme).copyWith(
        displayLarge: GoogleFonts.outfit(color: textPrimary, fontWeight: FontWeight.bold),
        titleLarge: GoogleFonts.outfit(color: textPrimary, fontWeight: FontWeight.w600),
        bodyLarge: GoogleFonts.outfit(color: textPrimary),
        bodyMedium: GoogleFonts.outfit(color: textSecondary),
      ),
      colorScheme: const ColorScheme.dark(
        primary: neonCyan,
        secondary: neonEmerald,
        surface: bgSurface,
      ),
    );
  }

  static BoxDecoration glassCard({double radius = 16, Color? borderColor}) {
    return BoxDecoration(
      color: bgGlass,
      borderRadius: BorderRadius.circular(radius),
      border: Border.all(
        color: borderColor ?? borderGlass,
        width: 1.2,
      ),
      boxShadow: [
        BoxShadow(
          color: Colors.black.withAlpha(50),
          blurRadius: 16,
          offset: const Offset(0, 8),
        ),
      ],
    );
  }
}
