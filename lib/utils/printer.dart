class Printer {
  static String red(String message) => '\u001b[31m$message\u001b[0m';
  static String green(String message) => '\u001b[32m$message\u001b[0m';
  static String yellow(String message) => '\u001b[33m$message\u001b[0m';

  const Printer._();
}
