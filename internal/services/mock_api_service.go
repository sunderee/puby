package services

import (
	"github.com/sunderee/puby/internal/models"
)

// MockAPIService is a mock implementation of APIServiceInterface
type MockAPIService struct {
	GetSDKReleaseFunc func() (*models.SDKReleaseWrapper, error)
	GetPackageFunc    func(packageName string) (*models.PackageWrapper, error)
}

// GetSDKRelease implements the APIServiceInterface
func (m *MockAPIService) GetSDKRelease() (*models.SDKReleaseWrapper, error) {
	return m.GetSDKReleaseFunc()
}

// GetPackage implements the APIServiceInterface
func (m *MockAPIService) GetPackage(packageName string) (*models.PackageWrapper, error) {
	return m.GetPackageFunc(packageName)
}
