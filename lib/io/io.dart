import 'dart:io';

import 'package:dart_scope_functions/dart_scope_functions.dart';
import 'package:ktuples/ktuples.dart';
import 'package:puby/io/models/dependency.model.dart';
import 'package:puby/io/models/environment.model.dart';
import 'package:puby/utils/types.dart';
import 'package:yaml/yaml.dart';

class IO {
  static Future<PubspecFile> readAndParsePubspec(File file) async {
    final rawFile = await file.readAsString();
    final result = loadYaml(rawFile) as YamlMap;

    final environment = result.entries
        .firstWhere((item) => item.key == 'environment')
        .let((it) => it.value as YamlMap)
        .let((it) => EnvironmentModel(
              dartSDK: it['sdk'] as String,
              flutterSDK: it['flutter'] as String?,
            ));

    final firstMap = result['dependencies'] as YamlMap?;
    final secondMap = result['dev_dependencies'] as YamlMap?;
    final dependencies = [firstMap, secondMap]
        .whereType<YamlMap>()
        .map((item) => item.entries)
        .expand((item) => item)
        .where((item) => item.key is String && item.value is String)
        .map((item) => [item.key as String, item.value as String])
        .map((item) => DependencyModel(name: item.first, version: item.last))
        .toSet();

    return Pair(environment, dependencies);
  }

  Future<void> writeChangesToPubspec(
    File file,
    Pair<String, String?> environment,
    Iterable<Triple<String, String, String>> upgradableDependencies,
  ) async {
    final rawFile = await file.readAsString();
    final lines = rawFile.split('\n');

    final startOfEnvironment = lines.indexOf('environment:');
    environment.first
        .let((it) => lines[startOfEnvironment + 1] = '  sdk: "$it"');
    environment.second
        ?.let((it) => lines[startOfEnvironment + 2] = '  flutter: "$it"');

    for (var item in upgradableDependencies) {
      lines
          .firstWhere((line) => line.contains(item.first))
          .let((it) => lines.indexOf(it))
          .let((it) => lines[it] = '  ${item.first}: ^${item.third}');
    }

    await file.writeAsString(lines.join('\n'));
  }
}
