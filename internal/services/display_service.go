package services

import (
	"fmt"
	"strings"

	"github.com/sunderee/puby/internal/models"
)

// DisplayService is responsible for displaying updates to the console
type DisplayService struct{}

// NewDisplayService creates a new instance of DisplayService
func NewDisplayService() DisplayServiceInterface {
	return &DisplayService{}
}

// PrintUpdate prints update information to the console with colorful output
func (s *DisplayService) PrintUpdate(update *models.Update) {
	if update == nil {
		fmt.Println("No updates available.")
		return
	}

	// Print environment updates
	if update.EnvironmentUpdate != nil {
		printEnvironmentUpdate(update.EnvironmentUpdate)
	}

	// Print dependency updates
	if len(update.DependencyUpdates) > 0 {
		printDependencyUpdates(update.DependencyUpdates)
	}

	// If no updates were printed, show a message
	if update.EnvironmentUpdate == nil && len(update.DependencyUpdates) == 0 {
		fmt.Println("Everything is up to date!")
	}
}

// printEnvironmentUpdate prints information about environment updates
func printEnvironmentUpdate(env *models.EnvironmentUpdate) {
	fmt.Println("\033[1;36m=== SDK Updates ===\033[0m")

	if env.DartSDKVersion != nil {
		fmt.Printf("\033[1;33mDart SDK:\033[0m \033[0;32m%s\033[0m\n", *env.DartSDKVersion)
	}

	if env.FlutterSDKVersion != nil {
		fmt.Printf("\033[1;33mFlutter SDK:\033[0m \033[0;32m%s\033[0m\n", *env.FlutterSDKVersion)
	}

	fmt.Println()
}

// printDependencyUpdates prints information about dependency updates
func printDependencyUpdates(deps []models.DependencyUpdate) {
	fmt.Println("\033[1;36m=== Dependency Updates ===\033[0m")

	// Find the maximum length of dependency names for proper alignment
	maxNameLength := 0
	for _, dep := range deps {
		if len(dep.Name) > maxNameLength {
			maxNameLength = len(dep.Name)
		}
	}

	// Print each dependency with proper alignment
	for _, dep := range deps {
		namePadding := strings.Repeat(" ", maxNameLength-len(dep.Name))
		fmt.Printf("\033[1;33m%s\033[0m%s: \033[0;31m%s\033[0m → \033[0;32m%s\033[0m\n",
			dep.Name,
			namePadding,
			dep.CurrentVersion,
			dep.LatestVersion)
	}

	fmt.Println()
}
