package services

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/sunderee/puby/internal/models"
)

func TestFileWriterService_WriteUpdates(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "puby-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Test cases
	tests := []struct {
		name            string
		initialContent  string
		update          *models.Update
		expectedContent string
		expectError     bool
	}{
		{
			name: "update dart SDK version",
			initialContent: `name: test_app
description: A test Flutter application.
environment:
  sdk: "2.18.0"
  flutter: "3.0.0"
dependencies:
  flutter:
    sdk: flutter
  http: ^0.13.3
`,
			update: &models.Update{
				EnvironmentUpdate: &models.EnvironmentUpdate{
					DartSDKVersion: createStringPtr("2.19.0"),
				},
			},
			expectedContent: `name: test_app
description: A test Flutter application.
environment:
  sdk: "2.19.0"
  flutter: "3.0.0"
dependencies:
  flutter:
    sdk: flutter
  http: ^0.13.3
`,
			expectError: false,
		},
		{
			name: "update flutter SDK version",
			initialContent: `name: test_app
description: A test Flutter application.
environment:
  sdk: "2.18.0"
  flutter: "3.0.0"
dependencies:
  flutter:
    sdk: flutter
  http: ^0.13.3
`,
			update: &models.Update{
				EnvironmentUpdate: &models.EnvironmentUpdate{
					FlutterSDKVersion: createStringPtr("3.19.0"),
				},
			},
			expectedContent: `name: test_app
description: A test Flutter application.
environment:
  sdk: "2.18.0"
  flutter: "3.19.0"
dependencies:
  flutter:
    sdk: flutter
  http: ^0.13.3
`,
			expectError: false,
		},
		{
			name: "update dependency version",
			initialContent: `name: test_app
description: A test Flutter application.
environment:
  sdk: "2.18.0"
  flutter: "3.0.0"
dependencies:
  flutter:
    sdk: flutter
  http: ^0.13.3
`,
			update: &models.Update{
				DependencyUpdates: []models.DependencyUpdate{
					{
						Name:           "http",
						CurrentVersion: "0.13.3",
						LatestVersion:  "0.13.5",
					},
				},
			},
			expectedContent: `name: test_app
description: A test Flutter application.
environment:
  sdk: "2.18.0"
  flutter: "3.0.0"
dependencies:
  flutter:
    sdk: flutter
  http: ^0.13.5
`,
			expectError: false,
		},
		{
			name: "update multiple items",
			initialContent: `name: test_app
description: A test Flutter application.
environment:
  sdk: "2.18.0"
  flutter: "3.0.0"
dependencies:
  flutter:
    sdk: flutter
  http: ^0.13.3
  path: ~1.8.0
`,
			update: &models.Update{
				EnvironmentUpdate: &models.EnvironmentUpdate{
					DartSDKVersion:    createStringPtr("2.19.0"),
					FlutterSDKVersion: createStringPtr("3.19.0"),
				},
				DependencyUpdates: []models.DependencyUpdate{
					{
						Name:           "http",
						CurrentVersion: "0.13.3",
						LatestVersion:  "0.13.5",
					},
					{
						Name:           "path",
						CurrentVersion: "1.8.0",
						LatestVersion:  "1.8.3",
					},
				},
			},
			expectedContent: `name: test_app
description: A test Flutter application.
environment:
  sdk: "2.19.0"
  flutter: "3.19.0"
dependencies:
  flutter:
    sdk: flutter
  http: ^0.13.5
  path: ~1.8.3
`,
			expectError: false,
		},
		{
			name:           "nil update",
			initialContent: `name: test_app`,
			update:         nil,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a temporary pubspec.yaml file
			pubspecPath := filepath.Join(tempDir, "pubspec.yaml")
			err := os.WriteFile(pubspecPath, []byte(tt.initialContent), 0644)
			if err != nil {
				t.Fatalf("Failed to write test file: %v", err)
			}

			// Create a file writer service
			writer := NewFileWriterService(pubspecPath)

			// Call the WriteUpdates method
			err = writer.WriteUpdates(tt.update)

			// Check error expectation
			if tt.expectError {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)

			// Read the updated file
			updatedContent, err := os.ReadFile(pubspecPath)
			if err != nil {
				t.Fatalf("Failed to read updated file: %v", err)
			}

			// Compare the updated content with expected
			assert.Equal(t, tt.expectedContent, string(updatedContent))
		})
	}
}

// Helper function to create a pointer to a string
func createStringPtr(s string) *string {
	return &s
}
