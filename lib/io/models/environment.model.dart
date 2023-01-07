import 'package:meta/meta.dart';

@immutable
class EnvironmentModel {
  final String dartSDK;
  final String? flutterSDK;

  EnvironmentModel({
    required this.dartSDK,
    this.flutterSDK,
  });
}
