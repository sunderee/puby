package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/sunderee/puby/internal/models"
)

const (
	DEFAULT_SDK_RELEASE_URL = "https://storage.googleapis.com/flutter_infra_release/releases/releases_macos.json"
	DEFAULT_PACKAGE_URL     = "https://pub.dev/api/packages/%s"
	HTTP_METHOD             = "GET"
)

type APIService struct {
	Client        *http.Client
	SDKReleaseURL string
	PackageURL    string
}

func NewAPIService() *APIService {
	return &APIService{
		Client:        &http.Client{},
		SDKReleaseURL: DEFAULT_SDK_RELEASE_URL,
		PackageURL:    DEFAULT_PACKAGE_URL,
	}
}

// GetSDKRelease fetches the latest SDK release from the Flutter repository
func (s *APIService) GetSDKRelease() (*models.SDKReleaseWrapper, error) {
	request, err := http.NewRequest(HTTP_METHOD, s.SDKReleaseURL, nil)
	if err != nil {
		return nil, err
	}

	response, err := s.Client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var sdkRelease models.SDKReleaseWrapper
	err = json.Unmarshal(body, &sdkRelease)
	if err != nil {
		return nil, err
	}

	return &sdkRelease, nil
}

// GetPackage fetches the latest package data from the pub.dev packages repository
func (s *APIService) GetPackage(packageName string) (*models.PackageWrapper, error) {
	if packageName == "" {
		return nil, fmt.Errorf("package name cannot be empty")
	}

	request, err := http.NewRequest(HTTP_METHOD, fmt.Sprintf(s.PackageURL, packageName), nil)
	if err != nil {
		return nil, err
	}

	response, err := s.Client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode >= 400 {
		return nil, fmt.Errorf("server returned status code %d", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var packageWrapper models.PackageWrapper
	err = json.Unmarshal(body, &packageWrapper)
	if err != nil {
		return nil, err
	}

	return &packageWrapper, nil
}
