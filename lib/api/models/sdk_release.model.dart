import 'package:dart_scope_functions/dart_scope_functions.dart';
import 'package:meta/meta.dart';
import 'package:puby/utils/types.dart';

@immutable
final class SDKReleaseModel {
  final ReleaseModel beta;
  final ReleaseModel stable;

  SDKReleaseModel._({
    required this.beta,
    required this.stable,
  });

  factory SDKReleaseModel.fromJson(JsonObject json) {
    final betaHash = json['current_release']['beta'] as String;
    final stableHash = json['current_release']['stable'] as String;

    final releases = (json['releases'] as List<dynamic>).cast<JsonObject>();
    final beta = releases
        .firstWhere((item) => item['hash'] as String == betaHash)
        .let((it) => ReleaseModel.fromJson(it, isBeta: true));
    final stable = releases
        .firstWhere((item) => item['hash'] as String == stableHash)
        .let((it) => ReleaseModel.fromJson(it, isBeta: false));

    return SDKReleaseModel._(
      beta: beta,
      stable: stable,
    );
  }
}

@immutable
final class ReleaseModel {
  final String dartSDK;
  final String flutterSDK;

  ReleaseModel._({
    required this.dartSDK,
    required this.flutterSDK,
  });

  factory ReleaseModel.fromJson(
    JsonObject json, {
    required bool isBeta,
  }) {
    final dartSDK = isBeta
        ? (json['dart_sdk_version'] as String)
            .let((it) => RegExp(r'(?<=build )[\w\.\-]+(?=\))').stringMatch(it))
            .let((it) => it ?? '')
        : json['dart_sdk_version'] as String;
    final flutterSDK = json['version'] as String;

    return ReleaseModel._(
      dartSDK: dartSDK,
      flutterSDK: flutterSDK,
    );
  }
}
