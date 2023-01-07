import 'dart:async';

import 'package:puby/puby.dart';

void main(List<String> arguments) {
  final cli = parseArguments(arguments);
  final puby = PUBy.fromConfiguration(cli);

  runZoned(() async {
    await puby.run();
  });
}
