package config

// This configuration is the result of CLI argument parsing. It instructs the
// tool on its behavior.
type CLIConfig struct {
	// If this flag is set, we will be comparing Dart (and Flutter) SDK versions
	// against the latest beta versions available.
	UseBetaSDKVersions *bool

	// If this flag is set, we will be checking if the Flutter SDK version is up
	// to date. If the flag is missing, only the Dart SDK version will be checked.
	CheckFlutterSDKVersion *bool

	// This slice contains the packages that we need to check for updates. If it's
	// empty or not set, all packages will be checked unless the exclusion list
	// is also set. In any case, the inclusion list will take precedence over the
	// exclusion list. If there's a package both in the inclusion and exclusion list,
	// it will cause a conflict and the program will exit with an error.
	IncludePackages *[]string

	// This slice contains the packages that we need to exclude from update check.
	// In case this slice is empty or not set, then either all packages will be checked
	// or only the ones in the inclusion list. When it is set, the exclusion list simply
	// tells which packages should not be checked.
	ExcludePackages *[]string

	// If this flag is not null or set to true, we will be writing the changes to the
	// pubspec.yaml file. Otherwise, we are running in the dry-run mode and only
	// printing the changes to the console.
	WriteChangesToFile *bool
}
