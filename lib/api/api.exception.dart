import 'dart:io';

import 'package:meta/meta.dart';

@immutable
final class ApiException implements Exception {
  final HttpClientResponse response;

  const ApiException(this.response);

  @override
  String toString() => '${response.statusCode}: ${response.reasonPhrase}';
}
