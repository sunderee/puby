import 'dart:io';

import 'package:args/args.dart';
import 'package:dart_scope_functions/dart_scope_functions.dart';
import 'package:meta/meta.dart';
import 'package:puby/utils/printer.dart';

enum DependencyEnum {
  production,
  development,
  all;

  static DependencyEnum from(String dependencyString) {
    switch (dependencyString) {
      case "production":
        return DependencyEnum.production;
      case "development":
        return DependencyEnum.development;
      case "":
      case "all":
        return DependencyEnum.all;
      default:
        throw ArgumentError("Invalid dependency string: $dependencyString");
    }
  }
}

enum SDKEnum {
  stable,
  beta;

  static SDKEnum from(String sdkString) {
    switch (sdkString) {
      case "":
      case "stable":
        return SDKEnum.stable;
      case "beta":
        return SDKEnum.beta;
      default:
        throw ArgumentError("Invalid SDK string: $sdkString");
    }
  }
}

@immutable
class AppArguments {
  final String path;
  final DependencyEnum dependency;
  final SDKEnum sdk;
  final List<String>? includeList;
  final List<String>? excludeList;
  final bool shouldUseUnstable;
  final bool beVerbose;
  final bool shouldWriteChanges;

  const AppArguments({
    required this.path,
    required this.dependency,
    required this.sdk,
    this.includeList,
    this.excludeList,
    required this.shouldUseUnstable,
    required this.beVerbose,
    required this.shouldWriteChanges,
  });
}

AppArguments parseArguments(List<String> arguments) {
  final parser = ArgParser()
    ..addOption(
      'path',
      abbr: 'p',
      help: 'Absolute path to the pubspec.yaml file.',
    )
    ..addOption(
      'dependencies',
      abbr: 'd',
      allowed: ['production', 'development', 'all'],
      defaultsTo: 'all',
      help: 'Which section of the pubspec.yaml file to check.',
    )
    ..addOption(
      'sdk',
      abbr: 's',
      allowed: ['stable', 'beta'],
      defaultsTo: 'stable',
      help: 'From which channel (stable or beta) should the Dart/Flutter SDK '
          'constraints be checked for latest and updated.',
    )
    ..addOption(
      'include',
      abbr: 'i',
      help: 'Sets a list of comma-separated dependencies to check for '
          'updates. This option works together with -d/--dependencies. '
          '-i/--include and -o/--omit cancel each other out, so only '
          '-d/--dependencies will be considered.',
    )
    ..addOption(
      'omit',
      abbr: 'o',
      help: 'Sets a list of comme-separated dependencies to be ignored in '
          'version upgrade checks. This option works together with '
          '-d/--dependencies. -o/--omit and -i/--include cancel each other out, '
          'so only --d/--dependencies will be considered.',
    )
    ..addFlag(
      'unstable',
      abbr: 'u',
      defaultsTo: false,
      help: 'Should you allow for unstable (alpha/beta/dev) versions.',
    )
    ..addFlag(
      'verbose',
      abbr: 'v',
      defaultsTo: false,
      help: 'Should you output more information.',
    )
    ..addFlag(
      'write',
      abbr: 'w',
      defaultsTo: true,
      help: 'Should you write changes to the pubspec.yaml file.',
    )
    ..addFlag(
      'help',
      abbr: 'h',
      help: 'Show the usage syntax.',
    );

  try {
    final results = parser.parse(arguments);
    if (results.wasParsed('help')) {
      print(Printer.yellow('Dart-based CLI app for managing dependencies.'));
      print(parser.usage);
      exit(0);
    }

    if (!results.wasParsed('path')) {
      print(Printer.red('You need to provide the absolute path.'));
      print(parser.usage);
      exit(1);
    }

    return AppArguments(
      path: results['path'] as String,
      dependency: DependencyEnum.from(results['dependencies'] as String),
      sdk: SDKEnum.from(results['sdk'] as String),
      includeList: (results['include'] as String?)
          ?.let((it) => it.split(','))
          .map((item) => item.trim())
          .toList(),
      excludeList: (results['omit'] as String?)
          ?.let((it) => it.split(','))
          .map((item) => item.trim())
          .toList(),
      shouldUseUnstable: results['unstable'] as bool,
      beVerbose: results['verbose'] as bool,
      shouldWriteChanges: results['write'] as bool,
    );
  } catch (e) {
    print(Printer.red('Invalid arguments.'));
    print(parser.usage);
    exit(1);
  }
}
