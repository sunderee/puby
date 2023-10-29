import 'package:equatable/equatable.dart';
import 'package:meta/meta.dart';

@immutable
final class EnvironmentModel extends Equatable {
  final String dartSDK;
  final String? flutterSDK;

  EnvironmentModel({
    required this.dartSDK,
    this.flutterSDK,
  });

  @override
  List<Object?> get props => [dartSDK, flutterSDK];
}
