package services

import (
	"github.com/sunderee/puby/internal/models"
)

// FileWriterInterface defines the interface for writing updates to files
type FileWriterInterface interface {
	WriteUpdates(update *models.Update) error
}
