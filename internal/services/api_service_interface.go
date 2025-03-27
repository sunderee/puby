package services

import (
	"github.com/sunderee/puby/internal/models"
)

// APIServiceInterface defines the interface for API service operations
type APIServiceInterface interface {
	GetSDKRelease() (*models.SDKReleaseWrapper, error)
	GetPackage(packageName string) (*models.PackageWrapper, error)
}
