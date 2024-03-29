import 'package:meta/meta.dart';
import 'package:puby/cli/cli.enums.dart';

@immutable
final class CLIConfiguration {
  final String path;
  final DependencyEnum targetDependencies;
  final SDKEnum projectSDKChannel;
  final Set<String>? includeDependenciesSet;
  final Set<String>? excludeDependenciesSet;
  final bool isFlutter;
  final bool allowUnstable;
  final bool writeToFile;

  const CLIConfiguration({
    required this.path,
    required this.targetDependencies,
    required this.projectSDKChannel,
    required this.includeDependenciesSet,
    required this.excludeDependenciesSet,
    required this.isFlutter,
    required this.allowUnstable,
    required this.writeToFile,
  });
}
