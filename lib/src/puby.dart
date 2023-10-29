import 'dart:io';

import 'package:collection/collection.dart';
import 'package:dart_scope_functions/dart_scope_functions.dart';
import 'package:puby/api/api.service.dart';
import 'package:puby/api/models/sdk_release.model.dart';
import 'package:puby/cli/cli.configuration.dart';
import 'package:puby/cli/cli.enums.dart';
import 'package:puby/io/io.dart';
import 'package:puby/io/models/environment.model.dart';
import 'package:puby/utils/printer.dart';
import 'package:puby/utils/types.dart';

final class PUBy {
  final CLIConfiguration cli;
  final IApiService _apiService;
  final IO _io;

  PUBy.fromConfiguration(this.cli)
      : _apiService = IApiService.instance(),
        _io = IO();

  Future<void> run() async {
    final pubspecFile = File(cli.path);
    final file = await _io.readAndParsePubspec(pubspecFile);

    final currentSDK = file.$1;
    final latestSDK = await _apiService.getSDKReleases();
    final fullEnvironment = (currentSDK, latestSDK);
    Printer.sdkPrint(
      cli.projectSDKChannel,
      cli.isFlutter,
      fullEnvironment,
    ).also((it) => print(it));

    var allDependencies = file.$2;
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
        .map((item) =>
            latest.firstWhereOrNull((e) => item.name == e.name)?.let((it) => (
                  it.name,
                  item.version,
                  cli.allowUnstable
                      ? it.latestUnstable ?? it.latestStable
                      : it.latestStable,
                )))
        .where((item) => item != null)
        .cast<(String, String, String)>()
        .map((item) => (
              item.$1,
              item.$2.replaceFirst('^', ''),
              item.$3,
            ))
        .where((item) => item.$2 != item.$3);

    for (var item in upgradableDependencies) {
      print(Printer.upgradablePrint(item));
    }

    if (cli.writeToFile) {
      await _io.writeChangesToPubspec(
        pubspecFile,
        _getEnvironment(
          cli.projectSDKChannel,
          cli.isFlutter,
          fullEnvironment,
        ),
        upgradableDependencies,
      );
    }
  }

  Environment _getEnvironment(
    SDKEnum sdk,
    bool isFlutter,
    (EnvironmentModel, SDKReleaseModel) input,
  ) {
    var dartSDK = input.$1.dartSDK;
    var flutterSDK = input.$1.flutterSDK;

    if (sdk == SDKEnum.stable) {
      if (input.$1.dartSDK != input.$2.stable.dartSDK) {
        dartSDK = input.$2.stable.dartSDK;
      }

      if (isFlutter && (input.$1.flutterSDK != input.$2.stable.flutterSDK)) {
        flutterSDK = input.$2.stable.flutterSDK;
      }
    } else {
      if (input.$1.dartSDK != input.$2.beta.dartSDK) {
        dartSDK = input.$2.beta.dartSDK;
      }

      if (isFlutter && (input.$1.flutterSDK != input.$2.beta.flutterSDK)) {
        flutterSDK = input.$2.beta.flutterSDK;
      }
    }
    return (dartSDK, flutterSDK);
  }
}
