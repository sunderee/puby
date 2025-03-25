package models

type PackageWrapper struct {
	Name          string  `json:"name"`
	LatestVersion Package `json:"latest"`
}

type Package struct {
	Version         string         `json:"version"`
	Dependencies    map[string]any `json:"dependencies"`
	DevDependencies map[string]any `json:"dev_dependencies"`
}
