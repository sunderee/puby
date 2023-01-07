import 'package:ktuples/ktuples.dart';
import 'package:puby/api/models/sdk_release.model.dart';
import 'package:puby/cli/cli.enums.dart';
import 'package:puby/io/models/environment.model.dart';

class Printer {
  static String red(String message) => '\u001b[31m$message\u001b[0m';
  static String green(String message) => '\u001b[32m$message\u001b[0m';
  static String yellow(String message) => '\u001b[33m$message\u001b[0m';

  static String sdkPrint(
    SDKEnum sdk,
    Pair<EnvironmentModel, SDKReleaseModel> input,
  ) {
    final buffer = StringBuffer();

    if (sdk == SDKEnum.stable) {
      if (input.first.dartSDK != input.second.stable.dartSDK) {
        final oldSDK = '\u001b[31m${input.first.dartSDK}\u001b[0m';
        final newSDK = '\u001b[32m${input.second.stable.dartSDK}\u001b[0m';
        buffer.writeln('Dart SDK: $oldSDK -> $newSDK');
      }

      if (input.first.flutterSDK != null &&
          (input.first.flutterSDK != input.second.stable.flutterSDK)) {
        final oldSDK = '\u001b[31m${input.first.flutterSDK}\u001b[0m';
        final newSDK = '\u001b[32m${input.second.stable.flutterSDK}\u001b[0m';
        buffer.writeln('Flutter SDK: $oldSDK -> $newSDK');
      }
    } else {
      if (input.first.dartSDK != input.second.beta.dartSDK) {
        final oldSDK = '\u001b[31m${input.first.dartSDK}\u001b[0m';
        final newSDK = '\u001b[32m${input.second.beta.dartSDK}\u001b[0m';
        buffer.writeln('Dart SDK: $oldSDK -> $newSDK');
      }

      if (input.first.flutterSDK != null &&
          (input.first.flutterSDK != input.second.beta.flutterSDK)) {
        final oldSDK = '\u001b[31m${input.first.flutterSDK}\u001b[0m';
        final newSDK = '\u001b[32m${input.second.beta.flutterSDK}\u001b[0m';
        buffer.writeln('Flutter SDK: $oldSDK -> $newSDK');
      }
    }
    return buffer.toString();
  }

  static String upgradablePrint(Triple<String, String, String> input) {
    final oldVersion = '\u001b[31m${input.second}\u001b[0m';
    final newVersion = '\u001b[32m${input.third}\u001b[0m';
    return '${input.first}: $oldVersion -> $newVersion';
  }

  const Printer._();
}
