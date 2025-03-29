package services

import (
	"errors"
	"strings"

	"github.com/sunderee/puby/internal/config"
	"github.com/sunderee/puby/internal/models"
	"github.com/sunderee/puby/internal/parsers"
)

type UpdateService struct {
	Config        *config.CLIConfig
	PubspecParser parsers.PubspecParserInterface
	APIService    APIServiceInterface
}

func NewUpdateService(pubspecParser parsers.PubspecParserInterface, apiService APIServiceInterface) *UpdateService {
	return &UpdateService{
		PubspecParser: pubspecParser,
		APIService:    apiService,
	}
}

func (s *UpdateService) CheckForUpdates() (*models.Update, error) {
	// Open the pubspec.yaml file and parse it
	pubspec, err := s.PubspecParser.Parse()
	if err != nil {
		return nil, err
	}

	// Get the latest SDK release
	sdkRelease, err := s.APIService.GetSDKRelease()
	if err != nil {
		return nil, err
	}

	// Check if there's an update needed for the Dart and Flutter SDKs
	isDartSDKUpdateNeeded := s.isDartSDKUpdateNeeded(sdkRelease, pubspec)
	isFlutterSDKUpdateNeeded := s.isFlutterSDKUpdateNeeded(sdkRelease, pubspec)

	var environmentUpdate *models.EnvironmentUpdate = nil
	if isDartSDKUpdateNeeded || isFlutterSDKUpdateNeeded {
		// Update is needed for either one of them...
		var dartSDKVersion *string = nil
		var flutterSDKVersion *string = nil

		if isDartSDKUpdateNeeded {
			// Get the latest Dart SDK version (if different from the current one)
			dartSDKVersion = s.dartSDKToUpdateTo(sdkRelease, pubspec)
		}

		if isFlutterSDKUpdateNeeded {
			// Get the latest Flutter SDK version
			flutterSDKVersion = s.flutterSDKToUpdateTo(sdkRelease, pubspec)
		}

		environmentUpdate = &models.EnvironmentUpdate{
			DartSDKVersion:    dartSDKVersion,
			FlutterSDKVersion: flutterSDKVersion,
		}
	}

	// Check if there's a conflict between included and excluded packages
	if s.isThereAConflictBetweenIncludedAndExcludedPackages(pubspec) {
		return nil, errors.New("there's a conflict between included and excluded packages")
	}

	// Produce a slice of dependencies to update
	dependenciesToUpdate := s.produceSliceOfDependenciesToUpdate(pubspec)

	// Fetch latest dependency data from API for each dependency
	var dependencyDataFromAPI []*models.PackageWrapper
	for _, dependency := range dependenciesToUpdate {
		data, err := s.APIService.GetPackage(dependency)
		if err != nil {
			return nil, err
		}

		dependencyDataFromAPI = append(dependencyDataFromAPI, data)
	}

	// Produce a slice of dependency updates
	dependencyUpdates := s.produceSliceOfDependencyUpdates(dependenciesToUpdate, dependencyDataFromAPI)

	// Return the update object
	return &models.Update{
		EnvironmentUpdate: environmentUpdate,
		DependencyUpdates: dependencyUpdates,
	}, nil
}

func (s *UpdateService) isDartSDKUpdateNeeded(sdkRelease *models.SDKReleaseWrapper, pubspec *models.Pubspec) bool {
	var latestVersionHash string
	if s.Config.UseBetaSDKVersions != nil && *s.Config.UseBetaSDKVersions {
		latestVersionHash = sdkRelease.CurrentRelease.Beta
	} else {
		latestVersionHash = sdkRelease.CurrentRelease.Stable
	}

	var latestStableVersionDartSDK string
	for _, version := range sdkRelease.Releases {
		if version.Hash == latestVersionHash {
			latestStableVersionDartSDK = version.DartSDKVersion
			break
		}
	}

	cleanedCurrentVersion := cleanupVersionString(*pubspec.Environment.DartSDKVersion)
	cleanedLatestVersion := cleanupVersionString(latestStableVersionDartSDK)

	return cleanedCurrentVersion != cleanedLatestVersion
}

func (s *UpdateService) dartSDKToUpdateTo(sdkRelease *models.SDKReleaseWrapper, pubspec *models.Pubspec) *string {
	var latestVersionHash string
	if s.Config.UseBetaSDKVersions != nil && *s.Config.UseBetaSDKVersions {
		latestVersionHash = sdkRelease.CurrentRelease.Beta
	} else {
		latestVersionHash = sdkRelease.CurrentRelease.Stable
	}

	for _, version := range sdkRelease.Releases {
		if version.Hash == latestVersionHash {
			return &version.DartSDKVersion
		}
	}

	return nil
}

func (s *UpdateService) isFlutterSDKUpdateNeeded(sdkRelease *models.SDKReleaseWrapper, pubspec *models.Pubspec) bool {
	if s.Config.CheckFlutterSDKVersion != nil && *s.Config.CheckFlutterSDKVersion {
		var latestVersionHash string
		if s.Config.UseBetaSDKVersions != nil && *s.Config.UseBetaSDKVersions {
			latestVersionHash = sdkRelease.CurrentRelease.Beta
		} else {
			latestVersionHash = sdkRelease.CurrentRelease.Stable
		}

		var latestStableVersionFlutterSDK string
		for _, version := range sdkRelease.Releases {
			if version.Hash == latestVersionHash {
				latestStableVersionFlutterSDK = version.FlutterSDKVersion
				break
			}
		}

		cleanedCurrentVersion := cleanupVersionString(*pubspec.Environment.FlutterSDKVersion)
		cleanedLatestVersion := cleanupVersionString(latestStableVersionFlutterSDK)

		return cleanedCurrentVersion != cleanedLatestVersion
	}

	return false
}

func (s *UpdateService) flutterSDKToUpdateTo(sdkRelease *models.SDKReleaseWrapper, pubspec *models.Pubspec) *string {
	var latestVersionHash string
	if s.Config.UseBetaSDKVersions != nil && *s.Config.UseBetaSDKVersions {
		latestVersionHash = sdkRelease.CurrentRelease.Beta
	} else {
		latestVersionHash = sdkRelease.CurrentRelease.Stable
	}

	for _, version := range sdkRelease.Releases {
		if version.Hash == latestVersionHash {
			return &version.FlutterSDKVersion
		}
	}

	return nil
}

func (s *UpdateService) isThereAConflictBetweenIncludedAndExcludedPackages(pubspec *models.Pubspec) bool {
	var includedPackages []string
	if s.Config.IncludePackages != nil {
		includedPackages = *s.Config.IncludePackages
	}

	var excludedPackages []string
	if s.Config.ExcludePackages != nil {
		excludedPackages = *s.Config.ExcludePackages
	}

	for _, includedPackage := range includedPackages {
		for _, excludedPackage := range excludedPackages {
			if strings.Contains(includedPackage, excludedPackage) {
				return true
			}
		}
	}

	return false
}

func (s *UpdateService) produceSliceOfDependenciesToUpdate(pubspec *models.Pubspec) []string {
	var dependenciesToUpdate []string
	var includedPackages, excludedPackages []string

	// Get include and exclude packages if they exist
	if s.Config.IncludePackages != nil {
		includedPackages = *s.Config.IncludePackages
	}
	if s.Config.ExcludePackages != nil {
		excludedPackages = *s.Config.ExcludePackages
	}

	// Process all dependencies
	for dependencyName, dependencyVersion := range pubspec.Dependencies {
		// Skip non-string versions (like SDK references)
		if _, ok := dependencyVersion.(string); !ok {
			continue
		}

		// If no includes or excludes are specified, add all dependencies
		if len(includedPackages) == 0 && len(excludedPackages) == 0 {
			dependenciesToUpdate = append(dependenciesToUpdate, dependencyName)
			continue
		}

		// If includes are specified, only add if in the includes list
		if len(includedPackages) > 0 {
			for _, includedPackage := range includedPackages {
				if strings.Contains(dependencyName, includedPackage) {
					dependenciesToUpdate = append(dependenciesToUpdate, dependencyName)
					break
				}
			}
			continue
		}

		// If excludes are specified, add unless in the excludes list
		if len(excludedPackages) > 0 {
			isExcluded := false
			for _, excludedPackage := range excludedPackages {
				if strings.Contains(dependencyName, excludedPackage) {
					isExcluded = true
					break
				}
			}
			if !isExcluded {
				dependenciesToUpdate = append(dependenciesToUpdate, dependencyName)
			}
		}
	}

	return dependenciesToUpdate
}

func (s *UpdateService) produceSliceOfDependencyUpdates(dependenciesToUpdate []string, dependencyDataFromAPI []*models.PackageWrapper) []models.DependencyUpdate {
	var dependencyUpdates []models.DependencyUpdate

	// Get the pubspec to extract current versions
	pubspec, err := s.PubspecParser.Parse()
	if err != nil {
		return []models.DependencyUpdate{}
	}

	// Map package names to their API data for easier lookup
	packageDataMap := make(map[string]*models.PackageWrapper)
	for _, packageData := range dependencyDataFromAPI {
		packageDataMap[packageData.Name] = packageData
	}

	// Create dependency updates
	for _, dependencyName := range dependenciesToUpdate {
		// Get current version from pubspec
		if currentVersion, ok := pubspec.Dependencies[dependencyName]; ok {
			// Only handle string versions
			if currentVersionStr, ok := currentVersion.(string); ok {
				// Clean up the version string (remove ^, ~, >=, etc.)
				cleanedCurrentVersion := cleanupVersionString(currentVersionStr)

				// Get the latest version from API data
				if packageData, exists := packageDataMap[dependencyName]; exists {
					latestVersion := packageData.LatestVersion.Version

					// Only add to updates if the versions are different
					if cleanedCurrentVersion != latestVersion {
						dependencyUpdates = append(dependencyUpdates, models.DependencyUpdate{
							Name:           dependencyName,
							CurrentVersion: cleanedCurrentVersion,
							LatestVersion:  latestVersion,
						})
					}
				}
			}
		}
	}

	return dependencyUpdates
}

func cleanupVersionString(input string) string {
	input = strings.ReplaceAll(input, ">", "")
	input = strings.ReplaceAll(input, "<", "")
	input = strings.ReplaceAll(input, "=", "")
	input = strings.ReplaceAll(input, "^", "")
	input = strings.ReplaceAll(input, "~", "")

	return input
}
