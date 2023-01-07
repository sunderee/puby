import 'package:puby/api/api.service.dart';
import 'package:test/test.dart';

void main() {
  group('API service', () {
    final service = IApiService.instance();

    test('should should fetch latest versions', () async {
      final packages = [
        'flutter_bloc',
        'get_it',
        'injectable',
        'http',
        'http_parser',
        'hive_flutter',
        'equatable',
        'meta'
      ];

      final results = await service.getLatestPackageVersions(packages);
      expect(results.length, 8);
    });

    test('should fetch SDK releases', () async {
      final result = await service.getSDKReleases();

      expect(result.beta.dartSDK.isNotEmpty, true);
      expect(result.beta.flutterSDK.isNotEmpty, true);
      expect(result.stable.dartSDK.isNotEmpty, true);
      expect(result.stable.flutterSDK.isNotEmpty, true);
    });
  });
}
