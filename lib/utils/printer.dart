import 'package:puby/api/models/sdk_release.model.dart';
import 'package:puby/cli/cli.enums.dart';
import 'package:puby/io/models/environment.model.dart';

final class Printer {
  static String red(String message) => '\u001b[31m$message\u001b[0m';
  static String green(String message) => '\u001b[32m$message\u001b[0m';
  static String yellow(String message) => '\u001b[33m$message\u001b[0m';

  static String sdkPrint(
    SDKEnum sdk,
    bool isFlutter,
    (EnvironmentModel, SDKReleaseModel) input,
  ) {
    final buffer = StringBuffer();

    if (sdk == SDKEnum.stable) {
      if (input.$1.dartSDK != input.$2.stable.dartSDK) {
        final oldSDK = '\u001b[31m${input.$1.dartSDK}\u001b[0m';
        final newSDK = '\u001b[32m${input.$2.stable.dartSDK}\u001b[0m';
        buffer.writeln('Dart SDK: $oldSDK -> $newSDK');
      }

      if (isFlutter && (input.$1.flutterSDK != input.$2.stable.flutterSDK)) {
        final oldSDK = '\u001b[31m${input.$1.flutterSDK}\u001b[0m';
        final newSDK = '\u001b[32m${input.$2.stable.flutterSDK}\u001b[0m';
        buffer.writeln('Flutter SDK: $oldSDK -> $newSDK');
      }
    } else {
      if (input.$1.dartSDK != input.$2.beta.dartSDK) {
        final oldSDK = '\u001b[31m${input.$1.dartSDK}\u001b[0m';
        final newSDK = '\u001b[32m${input.$2.beta.dartSDK}\u001b[0m';
        buffer.writeln('Dart SDK: $oldSDK -> $newSDK');
      }

      if (isFlutter && (input.$1.flutterSDK != input.$2.beta.flutterSDK)) {
        final oldSDK = '\u001b[31m${input.$1.flutterSDK}\u001b[0m';
        final newSDK = '\u001b[32m${input.$2.beta.flutterSDK}\u001b[0m';
        buffer.writeln('Flutter SDK: $oldSDK -> $newSDK');
      }
    }
    return buffer.toString();
  }

  static String upgradablePrint((String, String, String) input) {
    final oldVersion = '\u001b[31m${input.$2}\u001b[0m';
    final newVersion = '\u001b[32m${input.$3}\u001b[0m';
    return '${input.$1}: $oldVersion -> $newVersion';
  }

  const Printer._();
}
