import 'dart:io';

import 'package:collection/collection.dart';
import 'package:dart_scope_functions/dart_scope_functions.dart';
import 'package:ktuples/ktuples.dart';
import 'package:puby/api/api.service.dart';
import 'package:puby/cli/cli.configuration.dart';
import 'package:puby/cli/cli.enums.dart';
import 'package:puby/io/io.dart';
import 'package:puby/utils/printer.dart';

class PUBy {
  final CLIConfiguration cli;
  final IApiService _apiService;
  final IO _io;

  PUBy.fromConfiguration(this.cli)
      : _apiService = IApiService.instance(),
        _io = IO();

  Future<void> run() async {
    final pubspecFile = File(cli.path);
    final file = await _io.readAndParsePubspec(pubspecFile);

    final currentSDK = file.first;
    final latestSDK = await _apiService.getSDKReleases();
    final fullEnvironment = Pair(currentSDK, latestSDK);
    Printer.sdkPrint(
      cli.projectSDKChannel,
      fullEnvironment,
    ).also((it) => print(it));

    var allDependencies = file.second;
    if (cli.includeDependenciesSet != null ||
        cli.excludeDependenciesSet != null) {
      if (cli.includeDependenciesSet != null) {
        allDependencies
            .where((item) => cli.includeDependenciesSet!.contains(item.name))
            .let((it) => allDependencies = it.toSet());
      } else if (cli.excludeDependenciesSet != null) {
        allDependencies
            .where((item) => !cli.excludeDependenciesSet!.contains(item.name))
            .let((it) => allDependencies = it.toSet());
      }
    }

    if (cli.targetDependencies == DependencyEnum.production) {
      allDependencies
          .where((item) => item.isProduction)
          .let((it) => allDependencies = it.toSet());
    } else if (cli.targetDependencies == DependencyEnum.development) {
      allDependencies
          .where((item) => !item.isProduction)
          .let((it) => allDependencies = it.toSet());
    }

    final packages = allDependencies.map((item) => item.name).toList();
    final latest = await _apiService.getLatestPackageVersions(packages);

    final upgradableDependencies = allDependencies
        .map((item) => latest
            .firstWhereOrNull((e) => item.name == e.name)
            ?.let((it) => Triple(
                  it.name,
                  item.version,
                  cli.allowUnstable
                      ? it.latestUnstable ?? it.latestStable
                      : it.latestStable,
                )))
        .whereType<Triple<String, String, String>>()
        .map((item) => Triple(
              item.first,
              item.second.replaceFirst('^', ''),
              item.third,
            ))
        .where((item) => item.second != item.third);

    for (var item in upgradableDependencies) {
      print(Printer.upgradablePrint(item));
    }

    if (cli.writeToFile) {
      // TODO: figure this part out
    }
  }
}
