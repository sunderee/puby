package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/sunderee/puby/internal/models"
)

const (
	SDK_RELEASE_URL = "https://storage.googleapis.com/flutter_infra_release/releases/releases_macos.json"
	PACKAGE_URL     = "https://pub.dev/api/packages/%s"
	HTTP_METHOD     = "GET"
)

type APIService struct {
	Client *http.Client
}

func NewAPIService() *APIService {
	return &APIService{
		Client: &http.Client{},
	}
}

func (s *APIService) GetSDKRelease() (*models.SDKReleaseWrapper, error) {
	request, err := http.NewRequest(HTTP_METHOD, SDK_RELEASE_URL, nil)
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

func (s *APIService) GetPackage(packageName string) (*models.PackageWrapper, error) {
	request, err := http.NewRequest(HTTP_METHOD, fmt.Sprintf(PACKAGE_URL, packageName), nil)
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

	var packageWrapper models.PackageWrapper
	err = json.Unmarshal(body, &packageWrapper)
	if err != nil {
		return nil, err
	}

	return &packageWrapper, nil
}
