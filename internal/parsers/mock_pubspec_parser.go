package parsers

import (
	"github.com/sunderee/puby/internal/models"
)

// MockPubspecParser is a mock implementation of PubspecParserInterface
type MockPubspecParser struct {
	ParseFunc func() (*models.Pubspec, error)
}

// Parse implements the PubspecParserInterface
func (m *MockPubspecParser) Parse() (*models.Pubspec, error) {
	return m.ParseFunc()
}
