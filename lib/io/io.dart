import 'dart:io';

import 'package:dart_scope_functions/dart_scope_functions.dart';
import 'package:puby/io/models/dependency.model.dart';
import 'package:puby/io/models/environment.model.dart';
import 'package:puby/utils/types.dart';
import 'package:yaml/yaml.dart';

final class IO {
  Future<PubspecFile> readAndParsePubspec(File file) async {
    final rawFile = await file.readAsString();
    final result = loadYaml(rawFile) as YamlMap;

    final environment = result.entries
        .firstWhere((item) => item.key == 'environment')
        .let((it) => it.value as YamlMap)
        .let((it) => EnvironmentModel(
              dartSDK: (it['sdk'] as String)
                  .replaceFirst('>', '')
                  .replaceFirst('<', '')
                  .replaceFirst('=', ''),
              flutterSDK: (it['flutter'] as String?)
                  ?.replaceFirst('>', '')
                  .replaceFirst('<', '')
                  .replaceFirst('=', ''),
            ));

    final dependencies = (result['dependencies'] as YamlMap?)
        ?.let((it) => it.entries)
        .where((item) => item.key is String && item.value is String)
        .map((item) => DependencyModel(
              isProduction: true,
              name: item.key as String,
              version: item.value as String,
            ));
    final devDependencies = (result['dev_dependencies'] as YamlMap?)
        ?.let((it) => it.entries)
        .where((item) => item.key is String && item.value is String)
        .map((item) => DependencyModel(
              isProduction: true,
              name: item.key as String,
              version: item.value as String,
            ));

    final allDependencies = [dependencies, devDependencies]
        .whereType<Iterable<DependencyModel>>()
        .expand((item) => item)
        .toSet();

    return (environment, allDependencies);
  }

  Future<void> writeChangesToPubspec(
    File file,
    Environment environment,
    UpgradableDependencies upgradableDependencies,
  ) async {
    final rawFile = await file.readAsString();
    final lines = rawFile.split('\n');

    final startOfEnvironment = lines.indexOf('environment:');
    environment.$1
        .let((it) => lines[startOfEnvironment + 1] = '  sdk: ">=$it"');
    environment.$2
        ?.let((it) => lines[startOfEnvironment + 2] = '  flutter: ">=$it"');

    for (var item in upgradableDependencies) {
      lines
          .firstWhere((line) => line.contains(item.$1))
          .let((it) => lines.indexOf(it))
          .let((it) => lines[it] = '  ${item.$1}: ^${item.$3}');
    }

    await file.writeAsString(lines.join('\n'));
  }
}
