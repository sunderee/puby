import 'package:dart_scope_functions/dart_scope_functions.dart';
import 'package:equatable/equatable.dart';
import 'package:meta/meta.dart';
import 'package:puby/utils/types.dart';

@immutable
final class PackageModel extends Equatable {
  final String name;
  final String latestStable;
  final String? latestUnstable;

  PackageModel._({
    required this.name,
    required this.latestStable,
    required this.latestUnstable,
  });

  factory PackageModel.fromJson(JsonObject json) {
    final latestStable = (json['latest'] as JsonObject)['version'] as String;
    final latestUnstable = (json['versions'] as List<dynamic>)
        .cast<JsonObject>()
        .map((item) => item['version'] as String)
        .last
        .let((it) => it == latestStable ? null : it);

    return PackageModel._(
      name: json['name'] as String,
      latestStable: latestStable,
      latestUnstable: latestUnstable,
    );
  }

  @override
  List<Object?> get props => [name, latestStable, latestUnstable];
}
