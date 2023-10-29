import 'package:puby/io/models/dependency.model.dart';
import 'package:puby/io/models/environment.model.dart';

typedef JsonObject = Map<String, dynamic>;
typedef Headers = Map<String, String>;
typedef DecodingCallback<T> = T Function(JsonObject json);

typedef PubspecFile = (EnvironmentModel, Set<DependencyModel>);

typedef Environment = (String, String?);
typedef UpgradableDependencies = Iterable<(String, String, String)>;
