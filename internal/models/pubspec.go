package models

type Pubspec struct {
	Environment     *PubspecEnvironment `yaml:"environment"`
	Dependencies    map[string]any      `yaml:"dependencies"`
	DevDependencies map[string]any      `yaml:"dev_dependencies"`
}

type PubspecEnvironment struct {
	DartSDKVersion    *string `yaml:"sdk"`
	FlutterSDKVersion *string `yaml:"flutter"`
}
