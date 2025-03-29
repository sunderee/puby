package services

import (
	"github.com/sunderee/puby/internal/models"
)

// DisplayServiceInterface defines the interface for display service operations
type DisplayServiceInterface interface {
	PrintUpdate(update *models.Update)
}
