import 'dart:io';

import 'package:args/args.dart';
import 'package:dart_scope_functions/dart_scope_functions.dart';
import 'package:puby/cli/cli.configuration.dart';
import 'package:puby/cli/cli.enums.dart';
import 'package:puby/utils/printer.dart';

CLIConfiguration parseArguments(List<String> arguments) {
  final parser = ArgParser()
    ..addOption(
      'path',
      abbr: 'p',
      help: 'Absolute path to the pubspec.yaml file.',
    )
    ..addSeparator('Optional settings')
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
          '-i/--include and -e/--exclude cancel each other out, so only '
          '-d/--dependencies will be considered.',
    )
    ..addOption(
      'exclude',
      abbr: 'e',
      help: 'Sets a list of comme-separated dependencies to be excluded in '
          'version upgrade checks. This option works together with '
          '-d/--dependencies. -e/--exclude and -i/--include cancel each other '
          'out, so only --d/--dependencies will be considered.',
    )
    ..addSeparator('Flags')
    ..addFlag(
      'flutter',
      abbr: 'f',
      defaultsTo: false,
      help: 'Sets if this is a Flutter project (for setting SDK constraints).',
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
      defaultsTo: false,
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
      print(parser.usage);
      exit(0);
    }

    final path = results['path'] as String?;
    if (path == null) {
      print(Printer.red('Option path is mandatory.'));
      exit(1);
    }

    if (path.isEmpty) {
      print(Printer.red('Option path cannot be empty.'));
      exit(1);
    }

    final dependencies = (results['dependencies'] as String)
        .let((it) => DependencyEnum.from(it));
    final sdk = (results['sdk'] as String).let((it) => SDKEnum.from(it));
    var include = (results['include'] as String?)
        ?.split(',')
        .map((item) => item.trim())
        .toSet();
    var exclude = (results['exclude'] as String?)
        ?.split(',')
        .map((item) => item.trim())
        .toSet();

    if (include != null && exclude != null) {
      print(Printer.yellow('include and exclude together take no effect'));
      include = null;
      exclude = null;
    }

    final isFlutter = results['flutter'] as bool;
    final allowUnstable = results['unstable'] as bool;
    final useVerboseOutput = results['verbose'] as bool;
    final writeToFile = results['write'] as bool;

    return CLIConfiguration(
      path: path,
      targetDependencies: dependencies,
      projectSDKChannel: sdk,
      includeDependenciesSet: include,
      excludeDependenciesSet: exclude,
      isFlutter: isFlutter,
      allowUnstable: allowUnstable,
      useVerboseOutput: useVerboseOutput,
      writeToFile: writeToFile,
    );
  } on FormatException catch (e) {
    print(Printer.red(e.message));
    exit(1);
  }
}
