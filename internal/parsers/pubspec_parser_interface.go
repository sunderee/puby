package parsers

import (
	"github.com/sunderee/puby/internal/models"
)

// PubspecParserInterface defines the interface for pubspec file parsing
type PubspecParserInterface interface {
	Parse() (*models.Pubspec, error)
}
