import 'package:dart_scope_functions/dart_scope_functions.dart';
import 'package:puby/cli/cli.enums.dart';
import 'package:puby/puby.dart';
import 'package:test/test.dart';

void main() {
  group('argument parser', () {
    test('parse path', () {
      final arguments = <String>['--path=path/to/file'];
      final result = parseArguments(arguments);

      expect(result.path, 'path/to/file');
    });

    test('parse dependencies', () {
      final allDependencies =
          ['--path=path/to/file'].let((it) => parseArguments(it));
      final productionDependencies = [
        '--path=path/to/file',
        '--dependencies=production',
      ].let((it) => parseArguments(it));
      final developmentDependencies = [
        '--path=path/to/file',
        '-d=development',
      ].let((it) => parseArguments(it));

      expect(
        allDependencies.targetDependencies,
        DependencyEnum.all,
      );
      expect(
        productionDependencies.targetDependencies,
        DependencyEnum.production,
      );
      expect(
        developmentDependencies.targetDependencies,
        DependencyEnum.development,
      );
    });
  });
}
