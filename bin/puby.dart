import 'dart:async';

import 'package:puby/puby.dart';

void main(List<String> arguments) => runZoned(() async {
      final cli = parseArguments(arguments);
      final puby = PUBy.fromConfiguration(cli);

      await puby.run();
    });
