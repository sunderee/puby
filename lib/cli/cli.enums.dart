enum DependencyEnum {
  production,
  development,
  all;

  static DependencyEnum from(String input) => switch (input) {
        'production' => DependencyEnum.production,
        'development' => DependencyEnum.development,
        '' || 'all' => DependencyEnum.all,
        _ => throw ArgumentError("Invalid dependency string: $input"),
      };
}

enum SDKEnum {
  stable,
  beta;

  static SDKEnum from(String input) => switch (input) {
        '' || 'stable' => SDKEnum.stable,
        'beta' => SDKEnum.beta,
        _ => throw ArgumentError("Invalid SDK string: $input"),
      };
}
