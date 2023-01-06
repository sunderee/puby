class Printer {
  static const String _red = '\u001b[31m';
  static const String _green = '\u001b[32m';
  static const String _yellow = '\u001b[33m';
  static const String _reset = '\u001b[0m';

  static String red(String message) => '$_red$message$_reset';
  static String green(String message) => '$_green$message$_reset';
  static String yellow(String message) => '$_yellow$message$_reset';

  const Printer._();
}
