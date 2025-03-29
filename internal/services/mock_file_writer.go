package services

import (
	"github.com/sunderee/puby/internal/models"
)

// MockFileWriter is a mock implementation of the FileWriterInterface
type MockFileWriter struct {
	WriteUpdatesFunc func(update *models.Update) error
}

// WriteUpdates implements the FileWriterInterface
func (m *MockFileWriter) WriteUpdates(update *models.Update) error {
	if m.WriteUpdatesFunc != nil {
		return m.WriteUpdatesFunc(update)
	}
	return nil
}
