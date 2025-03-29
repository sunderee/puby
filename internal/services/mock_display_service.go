package services

import (
	"github.com/sunderee/puby/internal/models"
)

// MockDisplayService is a mock implementation of DisplayServiceInterface
type MockDisplayService struct {
	PrintUpdateFunc func(update *models.Update)
}

// PrintUpdate implements the DisplayServiceInterface
func (m *MockDisplayService) PrintUpdate(update *models.Update) {
	if m.PrintUpdateFunc != nil {
		m.PrintUpdateFunc(update)
	}
}
