package config

type UpdateServiceConfig struct {
	UseBetaSDKVersions     *bool
	CheckFlutterSDKVersion *bool
	IncludePackages        *[]string
	ExcludePackages        *[]string
	WriteChangesToFile     *bool
}
