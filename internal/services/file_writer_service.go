package services

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/sunderee/puby/internal/models"
)

// FileWriterService is responsible for writing updates to the pubspec.yaml file
type FileWriterService struct {
	PubspecFilePath string
}

// NewFileWriterService creates a new instance of FileWriterService
func NewFileWriterService(pubspecFilePath string) FileWriterInterface {
	return &FileWriterService{
		PubspecFilePath: pubspecFilePath,
	}
}

// WriteUpdates writes the updates to the pubspec.yaml file
func (s *FileWriterService) WriteUpdates(update *models.Update) error {
	if update == nil {
		return fmt.Errorf("no updates to write")
	}

	// Read the file
	fileContent, err := os.ReadFile(s.PubspecFilePath)
	if err != nil {
		return fmt.Errorf("failed to read pubspec.yaml: %v", err)
	}

	// Convert to string for easier manipulation
	content := string(fileContent)

	// Apply SDK updates if needed
	if update.EnvironmentUpdate != nil {
		if update.EnvironmentUpdate.DartSDKVersion != nil {
			content = s.updateDartSDKVersion(content, *update.EnvironmentUpdate.DartSDKVersion)
		}
		if update.EnvironmentUpdate.FlutterSDKVersion != nil {
			content = s.updateFlutterSDKVersion(content, *update.EnvironmentUpdate.FlutterSDKVersion)
		}
	}

	// Apply dependency updates
	for _, dep := range update.DependencyUpdates {
		content = s.updateDependencyVersion(content, dep.Name, dep.LatestVersion)
	}

	// Write the updated content back to the file
	return os.WriteFile(s.PubspecFilePath, []byte(content), 0644)
}

// updateDartSDKVersion updates the Dart SDK version in the pubspec.yaml file
func (s *FileWriterService) updateDartSDKVersion(content, newVersion string) string {
	// Look for sdk: "..." or sdk: '...' pattern in the environment section
	sdkPattern := regexp.MustCompile(`(environment:[\s\S]*?sdk:[\s]*['"])([^'"]*?)(['"])`)
	return sdkPattern.ReplaceAllString(content, "${1}"+newVersion+"${3}")
}

// updateFlutterSDKVersion updates the Flutter SDK version in the pubspec.yaml file
func (s *FileWriterService) updateFlutterSDKVersion(content, newVersion string) string {
	// Look for flutter: "..." or flutter: '...' pattern in the environment section
	flutterPattern := regexp.MustCompile(`(environment:[\s\S]*?flutter:[\s]*['"])([^'"]*?)(['"])`)
	return flutterPattern.ReplaceAllString(content, "${1}"+newVersion+"${3}")
}

// updateDependencyVersion updates a specific dependency version in the pubspec.yaml file
func (s *FileWriterService) updateDependencyVersion(content, dependencyName, newVersion string) string {
	// This pattern matches:
	//   dependencyName: "any-version"
	//   dependencyName: ^0.13.3
	//   dependencyName: '~1.2.3'
	// with any amount of whitespace and quotes
	dep := regexp.QuoteMeta(dependencyName)
	pattern := regexp.MustCompile(`(\n\s*` + dep + `:\s*\^?~?)([\d\.]+)`)

	// Find all matches
	matches := pattern.FindStringSubmatch(content)
	if len(matches) > 2 {
		// Get the prefix (includes whitespace, name, and any version constraint like ^, ~)
		prefix := matches[1]
		return pattern.ReplaceAllString(content, prefix+newVersion)
	}

	return content
}

// extractVersionPrefix extracts the version constraint prefix (^, ~, >=, etc.) from a version string
func extractVersionPrefix(version string) string {
	// Check for common version constraints
	if strings.HasPrefix(version, "^") {
		return "^"
	} else if strings.HasPrefix(version, "~") {
		return "~"
	} else if strings.HasPrefix(version, ">=") {
		return ">="
	} else if strings.HasPrefix(version, ">") {
		return ">"
	}

	// No constraint found
	return ""
}
