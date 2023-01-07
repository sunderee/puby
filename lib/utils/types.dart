import 'package:ktuples/ktuples.dart';
import 'package:puby/io/models/dependency.model.dart';
import 'package:puby/io/models/environment.model.dart';

typedef JsonObject = Map<String, dynamic>;
typedef Headers = Map<String, String>;
typedef DecodingCallback<T> = T Function(JsonObject json);

typedef PubspecFile = Pair<EnvironmentModel, Set<DependencyModel>>;

typedef Environment = Pair<String, String?>;
typedef UpgradableDependencies = Iterable<Triple<String, String, String>>;
