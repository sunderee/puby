package services

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/sunderee/puby/internal/models"
)

func TestDisplayService_PrintUpdate(t *testing.T) {
	// Setup
	displayService := NewDisplayService()

	// Prepare for capturing stdout
	originalStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Helper function to get the captured output
	getCapturedOutput := func() string {
		_ = w.Close()
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, r)
		os.Stdout = originalStdout
		return buf.String()
	}

	t.Run("No updates", func(t *testing.T) {
		// Call the method
		displayService.PrintUpdate(nil)
		output := getCapturedOutput()

		// Reset output capture
		r, w, _ = os.Pipe()
		os.Stdout = w

		// Assert
		assert.Contains(t, output, "No updates available.")
	})

	t.Run("Everything up to date", func(t *testing.T) {
		// Create an update with no actual updates
		update := &models.Update{
			EnvironmentUpdate: nil,
			DependencyUpdates: []models.DependencyUpdate{},
		}

		// Call the method
		displayService.PrintUpdate(update)
		output := getCapturedOutput()

		// Reset output capture
		r, w, _ = os.Pipe()
		os.Stdout = w

		// Assert
		assert.Contains(t, output, "Everything is up to date!")
	})

	t.Run("SDK updates", func(t *testing.T) {
		// Create test data
		dartVersion := "3.0.0"
		flutterVersion := "3.19.0"
		update := &models.Update{
			EnvironmentUpdate: &models.EnvironmentUpdate{
				DartSDKVersion:    &dartVersion,
				FlutterSDKVersion: &flutterVersion,
			},
			DependencyUpdates: []models.DependencyUpdate{},
		}

		// Call the method
		displayService.PrintUpdate(update)
		output := getCapturedOutput()

		// Reset output capture
		r, w, _ = os.Pipe()
		os.Stdout = w

		// Assert
		assert.Contains(t, output, "SDK Updates")
		assert.Contains(t, output, "Dart SDK")
		assert.Contains(t, output, "3.0.0")
		assert.Contains(t, output, "Flutter SDK")
		assert.Contains(t, output, "3.19.0")
	})

	t.Run("Dependency updates", func(t *testing.T) {
		// Create test data
		update := &models.Update{
			EnvironmentUpdate: nil,
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
		}

		// Call the method
		displayService.PrintUpdate(update)
		output := getCapturedOutput()

		// Assert
		assert.Contains(t, output, "Dependency Updates")
		assert.Contains(t, output, "http")
		assert.Contains(t, output, "0.13.3")
		assert.Contains(t, output, "0.13.5")
		assert.Contains(t, output, "path")
		assert.Contains(t, output, "1.8.0")
		assert.Contains(t, output, "1.8.3")
	})

	// Restore stdout
	os.Stdout = originalStdout
}
