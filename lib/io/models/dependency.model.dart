import 'package:equatable/equatable.dart';
import 'package:meta/meta.dart';

@immutable
final class DependencyModel extends Equatable {
  final bool isProduction;
  final String name;
  final String version;

  DependencyModel({
    required this.isProduction,
    required this.name,
    required this.version,
  });

  @override
  List<Object?> get props => [isProduction, name, version];
}
