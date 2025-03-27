package models

type Update struct {
	EnvironmentUpdate *EnvironmentUpdate
	DependencyUpdates []DependencyUpdate
}

type EnvironmentUpdate struct {
	DartSDKVersion    *string
	FlutterSDKVersion *string
}

type DependencyUpdate struct {
	Name           string
	CurrentVersion string
	LatestVersion  string
}
