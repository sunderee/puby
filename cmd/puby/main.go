package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/sunderee/puby/internal/config"
	"github.com/sunderee/puby/internal/models"
	"github.com/sunderee/puby/internal/parsers"
	"github.com/sunderee/puby/internal/services"
)

const (
	appName    = "puby"
	appVersion = "2.0.0"
)

func main() {
	// Define command-line flags
	useBetaSDKs := flag.Bool("beta", false, "Use beta versions for SDK updates")
	checkFlutterSDK := flag.Bool("flutter", false, "Check Flutter SDK version")
	writeChanges := flag.Bool("write", false, "Write changes to pubspec.yaml (otherwise run in dry-run mode)")
	pubspecPath := flag.String("path", "pubspec.yaml", "Path to the pubspec.yaml file")
	showHelp := flag.Bool("help", false, "Show help message")
	showVersion := flag.Bool("version", false, "Show version information")

	// Define include/exclude packages flags
	var includePackages, excludePackages string
	flag.StringVar(&includePackages, "include", "", "Comma-separated list of packages to include in update check")
	flag.StringVar(&excludePackages, "exclude", "", "Comma-separated list of packages to exclude from update check")

	// Parse flags
	flag.Parse()

	// Show help if requested
	if *showHelp {
		printHelp()
		return
	}

	// Show version if requested
	if *showVersion {
		fmt.Printf("%s version %s\n", appName, appVersion)
		return
	}

	// Resolve the absolute path to pubspec.yaml
	absPath, err := resolveAbsolutePath(*pubspecPath)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// Check if the pubspec.yaml file exists
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		fmt.Printf("Error: pubspec.yaml not found at %s\n", absPath)
		os.Exit(1)
	}

	// Parse include/exclude packages
	var includeSlice, excludeSlice *[]string
	if includePackages != "" {
		includes := splitCommaSeparatedList(includePackages)
		includeSlice = &includes
	}
	if excludePackages != "" {
		excludes := splitCommaSeparatedList(excludePackages)
		excludeSlice = &excludes
	}

	// Create CLI config
	cliConfig := &config.CLIConfig{
		UseBetaSDKVersions:     useBetaSDKs,
		CheckFlutterSDKVersion: checkFlutterSDK,
		IncludePackages:        includeSlice,
		ExcludePackages:        excludeSlice,
		WriteChangesToFile:     writeChanges,
	}

	// Create services
	pubspecParser := parsers.NewPubspecParser(absPath)
	apiService := services.NewAPIService()
	updateService := services.NewUpdateService(pubspecParser, apiService)
	displayService := services.NewDisplayService()

	// Set the config in the update service
	updateService.Config = cliConfig

	// Check for updates
	fmt.Printf("Checking for updates in %s...\n", absPath)
	update, err := updateService.CheckForUpdates()
	if err != nil {
		fmt.Printf("Error checking for updates: %v\n", err)
		os.Exit(1)
	}

	// Display the updates
	displayService.PrintUpdate(update)

	// Write changes if needed
	if *writeChanges && hasUpdates(update) {
		fileWriter := services.NewFileWriterService(absPath)
		if err := fileWriter.WriteUpdates(update); err != nil {
			fmt.Printf("Error writing updates: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("\nUpdates have been written to pubspec.yaml")
	} else if hasUpdates(update) {
		fmt.Println("\nRunning in dry-run mode. Use --write flag to apply changes.")
	}
}

// printHelp prints the help message
func printHelp() {
	fmt.Printf("%s - A utility for managing Dart/Flutter package dependencies\n\n", appName)
	fmt.Println("Usage:")
	fmt.Printf("  %s [options]\n\n", appName)
	fmt.Println("Options:")
	flag.PrintDefaults()
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Printf("  %s                            # Check for updates in current directory\n", appName)
	fmt.Printf("  %s --path=/path/to/project    # Check for updates in a specific directory\n", appName)
	fmt.Printf("  %s --write                    # Apply updates to pubspec.yaml\n", appName)
	fmt.Printf("  %s --include=http,path        # Only check specific packages\n", appName)
	fmt.Printf("  %s --exclude=flutter_svg      # Check all packages except flutter_svg\n", appName)
}

// resolveAbsolutePath resolves the absolute path to pubspec.yaml
func resolveAbsolutePath(path string) (string, error) {
	// If the path is already absolute, return it
	if filepath.IsAbs(path) {
		return path, nil
	}

	// Get the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current working directory: %v", err)
	}

	// Join the current working directory with the relative path
	absPath := filepath.Join(cwd, path)
	return absPath, nil
}

// splitCommaSeparatedList splits a comma-separated list into a slice of strings
func splitCommaSeparatedList(list string) []string {
	if list == "" {
		return nil
	}

	// Split by comma
	items := strings.Split(list, ",")

	// Trim each item
	for i := range items {
		items[i] = strings.TrimSpace(items[i])
	}

	return items
}

// hasUpdates checks if there are any updates
func hasUpdates(update *models.Update) bool {
	if update == nil {
		return false
	}

	// Check for environment updates
	if update.EnvironmentUpdate != nil {
		if update.EnvironmentUpdate.DartSDKVersion != nil || update.EnvironmentUpdate.FlutterSDKVersion != nil {
			return true
		}
	}

	// Check for dependency updates
	if len(update.DependencyUpdates) > 0 {
		return true
	}

	return false
}
