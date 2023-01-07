import 'package:equatable/equatable.dart';
import 'package:meta/meta.dart';

@immutable
class DependencyModel extends Equatable {
  final String name;
  final String version;

  @override
  List<Object?> get props => [name, version];

  DependencyModel({
    required this.name,
    required this.version,
  });
}
