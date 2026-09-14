import 'package:flutter_test/flutter_test.dart';
import 'package:vify/services/update_service.dart';

void main() {
  group('UpdateService semver comparison tests', () {
    test('detects newer major version', () {
      expect(UpdateService.compareSemVer('v2.0.0', '1.9.9'), 1);
    });

    test('detects newer minor version', () {
      expect(UpdateService.compareSemVer('v1.3.0', 'v1.2.9'), 1);
    });

    test('detects newer patch version', () {
      expect(UpdateService.compareSemVer('v1.2.4', '1.2.3'), 1);
    });

    test('detects older version', () {
      expect(UpdateService.compareSemVer('v1.2.0', 'v1.2.1'), -1);
    });

    test('detects equal versions with different formatting', () {
      expect(UpdateService.compareSemVer('v1.2.3', '1.2.3'), 0);
      expect(UpdateService.compareSemVer('1.2.3+4', 'v1.2.3'), 0);
      expect(UpdateService.compareSemVer('v1.2.3-beta.1', '1.2.3'), 0);
    });
  });
}
