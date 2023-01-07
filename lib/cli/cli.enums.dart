enum DependencyEnum {
  production,
  development,
  all;

  static DependencyEnum from(String dependencyString) {
    switch (dependencyString) {
      case "production":
        return DependencyEnum.production;
      case "development":
        return DependencyEnum.development;
      case "":
      case "all":
        return DependencyEnum.all;
      default:
        throw ArgumentError("Invalid dependency string: $dependencyString");
    }
  }
}

enum SDKEnum {
  stable,
  beta;

  static SDKEnum from(String sdkString) {
    switch (sdkString) {
      case "":
      case "stable":
        return SDKEnum.stable;
      case "beta":
        return SDKEnum.beta;
      default:
        throw ArgumentError("Invalid SDK string: $sdkString");
    }
  }
}
