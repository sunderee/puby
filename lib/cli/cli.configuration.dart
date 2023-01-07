import 'package:meta/meta.dart';
import 'package:puby/cli/cli.enums.dart';

@immutable
class CLIConfiguration {
  final String path;
  final DependencyEnum targetDependencies;
  final SDKEnum projectSDKChannel;
  final Set<String>? includeDependenciesSet;
  final Set<String>? excludeDependenciesSet;
  final bool allowUnstable;
  final bool useVerboseOutput;
  final bool writeToFile;

  const CLIConfiguration({
    required this.path,
    required this.targetDependencies,
    required this.projectSDKChannel,
    required this.includeDependenciesSet,
    required this.excludeDependenciesSet,
    required this.allowUnstable,
    required this.useVerboseOutput,
    required this.writeToFile,
  });
}
