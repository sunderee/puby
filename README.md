# PUBy

CLI app for updating dependencies in Dart/Flutter projects.

If you are familiar with
[`npm-check-updates`](https://www.npmjs.com/package/npm-check-updates) or [`yarn-upgrade-all`](https://www.npmjs.com/package/yarn-upgrade-all) from the
Node.js world, then imagine **PUBy** being the same for Dart/Flutter. This is a
CLI app which enables you to update `pub.dev`-hosted package dependencies
declared in a `pubspec.yaml` file to their latest version, as well as update the
Dart/Flutter SDK constraints (under `environment`).

## Usage

First, fetch dependencies and compile the app to a self-contained executable:

```bash
dart pub get
dart compile exe --output=puby bin/puby.dart
```

Use `-h/--help` to learn how to use the package:

```
$ ./puby --help   
-p, --path             Absolute path to the pubspec.yaml file.

Optional settings
-d, --dependencies     Which section of the pubspec.yaml file to check.
                       [production, development, all (default)]
-s, --sdk              From which channel (stable or beta) should the Dart/Flutter SDK constraints be checked for latest and updated.
                       [stable (default), beta]
-i, --include          Sets a list of comma-separated dependencies to check for updates. This option works together with -d/--dependencies. -i/--include and -e/--exclude cancel each other out, so only -d/--dependencies will be considered.
-e, --exclude          Sets a list of comme-separated dependencies to be excluded in version upgrade checks. This option works together with -d/--dependencies. -e/--exclude and -i/--include cancel each other out, so only --d/--dependencies will be considered.

Flags
-f, --[no-]flutter     Sets if this is a Flutter project (for setting SDK constraints).
-u, --[no-]unstable    Should you allow for unstable (alpha/beta/dev) versions.
-v, --[no-]verbose     Should you output more information.
-w, --[no-]write       Should you write changes to the pubspec.yaml file.
-h, --[no-]help        Show the usage syntax.
```