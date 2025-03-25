package models

type SDKReleaseWrapper struct {
	CurrentRelease SDKReleaseHashes `json:"current_release"`
	Releases       []SDKRelease     `json:"releases"`
}

type SDKReleaseHashes struct {
	Beta   string `json:"beta"`
	Stable string `json:"stable"`
}

type SDKRelease struct {
	Hash              string `json:"hash"`
	Channel           string `json:"channel"`
	FlutterSDKVersion string `json:"version"`
	DartSDKVersion    string `json:"dart_sdk_version"`
}
