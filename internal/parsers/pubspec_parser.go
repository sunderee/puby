package parsers

import (
	"os"

	"github.com/sunderee/puby/internal/models"
	"gopkg.in/yaml.v3"
)

type PubspecParser struct {
	PubspecFilePath string
}

func NewPubspecParser(pubspecFilePath string) *PubspecParser {
	return &PubspecParser{
		PubspecFilePath: pubspecFilePath,
	}
}

// Open the pubspec.yaml file and parse it into a Pubspec struct
func (p *PubspecParser) Parse() (*models.Pubspec, error) {
	yamlFile, err := os.ReadFile(p.PubspecFilePath)
	if err != nil {
		return nil, err
	}

	var pubspec models.Pubspec
	if err := yaml.Unmarshal(yamlFile, &pubspec); err != nil {
		return nil, err
	}

	return &pubspec, nil
}
