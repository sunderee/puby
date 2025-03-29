package services

import (
	"errors"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/sunderee/puby/internal/config"
	"github.com/sunderee/puby/internal/models"
	"github.com/sunderee/puby/internal/parsers"
)

func TestUpdateService_CheckForUpdates(t *testing.T) {
	// Setup test cases
	tests := []struct {
		name           string
		pubspecParser  func() *parsers.MockPubspecParser
		apiService     func() *MockAPIService
		config         *config.CLIConfig
		expectedUpdate *models.Update
		expectedError  error
	}{
		{
			name: "PubspecParser error",
			pubspecParser: func() *parsers.MockPubspecParser {
				return &parsers.MockPubspecParser{
					ParseFunc: func() (*models.Pubspec, error) {
						return nil, errors.New("failed to parse pubspec")
					},
				}
			},
			apiService: func() *MockAPIService {
				return &MockAPIService{}
			},
			expectedUpdate: nil,
			expectedError:  errors.New("failed to parse pubspec"),
		},
		{
			name: "GetSDKRelease error",
			pubspecParser: func() *parsers.MockPubspecParser {
				return &parsers.MockPubspecParser{
					ParseFunc: func() (*models.Pubspec, error) {
						sdkVersion := ">=2.12.0 <3.0.0"
						flutterVersion := ">=2.5.0 <3.0.0"
						return &models.Pubspec{
							Environment: &models.PubspecEnvironment{
								DartSDKVersion:    &sdkVersion,
								FlutterSDKVersion: &flutterVersion,
							},
							Dependencies: map[string]any{
								"flutter": map[string]any{
									"sdk": "flutter",
								},
								"http": "^0.13.3",
							},
						}, nil
					},
				}
			},
			apiService: func() *MockAPIService {
				return &MockAPIService{
					GetSDKReleaseFunc: func() (*models.SDKReleaseWrapper, error) {
						return nil, errors.New("failed to get SDK release")
					},
				}
			},
			expectedUpdate: nil,
			expectedError:  errors.New("failed to get SDK release"),
		},
		{
			name: "No updates needed",
			pubspecParser: func() *parsers.MockPubspecParser {
				return &parsers.MockPubspecParser{
					ParseFunc: func() (*models.Pubspec, error) {
						sdkVersion := "3.0.0"
						flutterVersion := "3.19.0"
						return &models.Pubspec{
							Environment: &models.PubspecEnvironment{
								DartSDKVersion:    &sdkVersion,
								FlutterSDKVersion: &flutterVersion,
							},
							Dependencies: map[string]any{},
						}, nil
					},
				}
			},
			apiService: func() *MockAPIService {
				return &MockAPIService{
					GetSDKReleaseFunc: func() (*models.SDKReleaseWrapper, error) {
						return &models.SDKReleaseWrapper{
							CurrentRelease: models.SDKReleaseHashes{
								Stable: "abc123",
							},
							Releases: []models.SDKRelease{
								{
									Hash:              "abc123",
									DartSDKVersion:    "3.0.0",
									FlutterSDKVersion: "3.19.0",
								},
							},
						}, nil
					},
				}
			},
			config: &config.CLIConfig{
				CheckFlutterSDKVersion: boolPtr(true),
			},
			expectedUpdate: &models.Update{
				EnvironmentUpdate: nil,
				DependencyUpdates: nil,
			},
			expectedError: nil,
		},
		{
			name: "Dart SDK update needed",
			pubspecParser: func() *parsers.MockPubspecParser {
				return &parsers.MockPubspecParser{
					ParseFunc: func() (*models.Pubspec, error) {
						sdkVersion := "2.12.0"
						flutterVersion := "3.19.0"
						return &models.Pubspec{
							Environment: &models.PubspecEnvironment{
								DartSDKVersion:    &sdkVersion,
								FlutterSDKVersion: &flutterVersion,
							},
							Dependencies: map[string]any{},
						}, nil
					},
				}
			},
			apiService: func() *MockAPIService {
				return &MockAPIService{
					GetSDKReleaseFunc: func() (*models.SDKReleaseWrapper, error) {
						return &models.SDKReleaseWrapper{
							CurrentRelease: models.SDKReleaseHashes{
								Stable: "abc123",
							},
							Releases: []models.SDKRelease{
								{
									Hash:              "abc123",
									DartSDKVersion:    "3.0.0",
									FlutterSDKVersion: "3.19.0",
								},
							},
						}, nil
					},
				}
			},
			config: &config.CLIConfig{
				CheckFlutterSDKVersion: boolPtr(true),
			},
			expectedUpdate: &models.Update{
				EnvironmentUpdate: &models.EnvironmentUpdate{
					DartSDKVersion:    stringPtr("3.0.0"),
					FlutterSDKVersion: nil,
				},
				DependencyUpdates: nil,
			},
			expectedError: nil,
		},
		{
			name: "Flutter SDK update needed",
			pubspecParser: func() *parsers.MockPubspecParser {
				return &parsers.MockPubspecParser{
					ParseFunc: func() (*models.Pubspec, error) {
						sdkVersion := "3.0.0"
						flutterVersion := "2.5.0"
						return &models.Pubspec{
							Environment: &models.PubspecEnvironment{
								DartSDKVersion:    &sdkVersion,
								FlutterSDKVersion: &flutterVersion,
							},
							Dependencies: map[string]any{},
						}, nil
					},
				}
			},
			apiService: func() *MockAPIService {
				return &MockAPIService{
					GetSDKReleaseFunc: func() (*models.SDKReleaseWrapper, error) {
						return &models.SDKReleaseWrapper{
							CurrentRelease: models.SDKReleaseHashes{
								Stable: "abc123",
							},
							Releases: []models.SDKRelease{
								{
									Hash:              "abc123",
									DartSDKVersion:    "3.0.0",
									FlutterSDKVersion: "3.19.0",
								},
							},
						}, nil
					},
				}
			},
			config: &config.CLIConfig{
				CheckFlutterSDKVersion: boolPtr(true),
			},
			expectedUpdate: &models.Update{
				EnvironmentUpdate: &models.EnvironmentUpdate{
					DartSDKVersion:    nil,
					FlutterSDKVersion: stringPtr("3.19.0"),
				},
				DependencyUpdates: nil,
			},
			expectedError: nil,
		},
		{
			name: "Conflict between included and excluded packages",
			pubspecParser: func() *parsers.MockPubspecParser {
				return &parsers.MockPubspecParser{
					ParseFunc: func() (*models.Pubspec, error) {
						sdkVersion := "3.0.0"
						flutterVersion := "3.19.0"
						return &models.Pubspec{
							Environment: &models.PubspecEnvironment{
								DartSDKVersion:    &sdkVersion,
								FlutterSDKVersion: &flutterVersion,
							},
							Dependencies: map[string]any{
								"http": "^0.13.3",
							},
						}, nil
					},
				}
			},
			apiService: func() *MockAPIService {
				return &MockAPIService{
					GetSDKReleaseFunc: func() (*models.SDKReleaseWrapper, error) {
						return &models.SDKReleaseWrapper{
							CurrentRelease: models.SDKReleaseHashes{
								Stable: "abc123",
							},
							Releases: []models.SDKRelease{
								{
									Hash:              "abc123",
									DartSDKVersion:    "3.0.0",
									FlutterSDKVersion: "3.19.0",
								},
							},
						}, nil
					},
				}
			},
			config: &config.CLIConfig{
				IncludePackages: &[]string{"http"},
				ExcludePackages: &[]string{"http"},
			},
			expectedUpdate: nil,
			expectedError:  errors.New("there's a conflict between included and excluded packages"),
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockPubspecParser := tt.pubspecParser()
			mockAPIService := tt.apiService()

			service := NewUpdateService(mockPubspecParser, mockAPIService)
			service.Config = tt.config

			// Act
			update, err := service.CheckForUpdates()

			// Assert
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.expectedUpdate, update)
		})
	}
}

// Helper function to create a pointer to a bool
func boolPtr(b bool) *bool {
	return &b
}

// Helper function to create a pointer to a string
func stringPtr(s string) *string {
	return &s
}

func TestUpdateService_IsDartSDKUpdateNeeded(t *testing.T) {
	tests := []struct {
		name           string
		config         *config.CLIConfig
		sdkRelease     *models.SDKReleaseWrapper
		pubspec        *models.Pubspec
		expectedResult bool
	}{
		{
			name: "needs update - stable channel",
			config: &config.CLIConfig{
				UseBetaSDKVersions: boolPtr(false),
			},
			sdkRelease: &models.SDKReleaseWrapper{
				CurrentRelease: models.SDKReleaseHashes{
					Stable: "stable-hash",
				},
				Releases: []models.SDKRelease{
					{
						Hash:           "stable-hash",
						DartSDKVersion: "2.19.0",
					},
				},
			},
			pubspec: &models.Pubspec{
				Environment: &models.PubspecEnvironment{
					DartSDKVersion: stringPtr("2.18.0"),
				},
			},
			expectedResult: true,
		},
		{
			name: "no update needed - same version",
			config: &config.CLIConfig{
				UseBetaSDKVersions: boolPtr(false),
			},
			sdkRelease: &models.SDKReleaseWrapper{
				CurrentRelease: models.SDKReleaseHashes{
					Stable: "stable-hash",
				},
				Releases: []models.SDKRelease{
					{
						Hash:           "stable-hash",
						DartSDKVersion: "2.19.0",
					},
				},
			},
			pubspec: &models.Pubspec{
				Environment: &models.PubspecEnvironment{
					DartSDKVersion: stringPtr("2.19.0"),
				},
			},
			expectedResult: false,
		},
		{
			name: "needs update - beta channel",
			config: &config.CLIConfig{
				UseBetaSDKVersions: boolPtr(true),
			},
			sdkRelease: &models.SDKReleaseWrapper{
				CurrentRelease: models.SDKReleaseHashes{
					Beta: "beta-hash",
				},
				Releases: []models.SDKRelease{
					{
						Hash:           "beta-hash",
						DartSDKVersion: "2.20.0-beta",
					},
				},
			},
			pubspec: &models.Pubspec{
				Environment: &models.PubspecEnvironment{
					DartSDKVersion: stringPtr("2.19.0"),
				},
			},
			expectedResult: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &UpdateService{
				Config: tt.config,
			}
			result := service.isDartSDKUpdateNeeded(tt.sdkRelease, tt.pubspec)
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}

func TestUpdateService_IsFlutterSDKUpdateNeeded(t *testing.T) {
	tests := []struct {
		name           string
		config         *config.CLIConfig
		sdkRelease     *models.SDKReleaseWrapper
		pubspec        *models.Pubspec
		expectedResult bool
	}{
		{
			name: "check disabled",
			config: &config.CLIConfig{
				CheckFlutterSDKVersion: boolPtr(false),
			},
			sdkRelease:     &models.SDKReleaseWrapper{},
			pubspec:        &models.Pubspec{},
			expectedResult: false,
		},
		{
			name: "needs update - stable channel",
			config: &config.CLIConfig{
				CheckFlutterSDKVersion: boolPtr(true),
				UseBetaSDKVersions:     boolPtr(false),
			},
			sdkRelease: &models.SDKReleaseWrapper{
				CurrentRelease: models.SDKReleaseHashes{
					Stable: "stable-hash",
				},
				Releases: []models.SDKRelease{
					{
						Hash:              "stable-hash",
						FlutterSDKVersion: "3.0.0",
					},
				},
			},
			pubspec: &models.Pubspec{
				Environment: &models.PubspecEnvironment{
					FlutterSDKVersion: stringPtr("2.0.0"),
				},
			},
			expectedResult: true,
		},
		{
			name: "no update needed - same version",
			config: &config.CLIConfig{
				CheckFlutterSDKVersion: boolPtr(true),
				UseBetaSDKVersions:     boolPtr(false),
			},
			sdkRelease: &models.SDKReleaseWrapper{
				CurrentRelease: models.SDKReleaseHashes{
					Stable: "stable-hash",
				},
				Releases: []models.SDKRelease{
					{
						Hash:              "stable-hash",
						FlutterSDKVersion: "3.0.0",
					},
				},
			},
			pubspec: &models.Pubspec{
				Environment: &models.PubspecEnvironment{
					FlutterSDKVersion: stringPtr("3.0.0"),
				},
			},
			expectedResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &UpdateService{
				Config: tt.config,
			}
			result := service.isFlutterSDKUpdateNeeded(tt.sdkRelease, tt.pubspec)
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}

func TestUpdateService_ProduceSliceOfDependenciesToUpdate(t *testing.T) {
	tests := []struct {
		name                 string
		config               *config.CLIConfig
		pubspec              *models.Pubspec
		expectedDependencies []string
	}{
		{
			name: "include packages only",
			config: &config.CLIConfig{
				IncludePackages: &[]string{"http", "path"},
			},
			pubspec: &models.Pubspec{
				Dependencies: map[string]any{
					"http":        "^0.13.3",
					"path":        "^1.8.0",
					"flutter_svg": "^1.0.0",
				},
			},
			expectedDependencies: []string{"http", "path"},
		},
		{
			name: "exclude packages only",
			config: &config.CLIConfig{
				ExcludePackages: &[]string{"flutter_svg"},
			},
			pubspec: &models.Pubspec{
				Dependencies: map[string]any{
					"http":        "^0.13.3",
					"path":        "^1.8.0",
					"flutter_svg": "^1.0.0",
				},
			},
			expectedDependencies: []string{"http", "path"},
		},
		{
			name:   "non-string version ignored",
			config: &config.CLIConfig{},
			pubspec: &models.Pubspec{
				Dependencies: map[string]any{
					"http": "^0.13.3",
					"flutter": map[string]string{
						"sdk": "flutter",
					},
				},
			},
			expectedDependencies: []string{"http"},
		},
		{
			name:   "default case - no include or exclude filters",
			config: &config.CLIConfig{},
			pubspec: &models.Pubspec{
				Dependencies: map[string]any{
					"http":        "^0.13.3",
					"path":        "^1.8.0",
					"flutter_svg": "^1.0.0",
					"flutter": map[string]string{
						"sdk": "flutter",
					},
				},
			},
			expectedDependencies: []string{"http", "path", "flutter_svg"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &UpdateService{
				Config: tt.config,
			}
			result := service.produceSliceOfDependenciesToUpdate(tt.pubspec)

			// Since the ordering is not guaranteed in maps
			assert.ElementsMatch(t, tt.expectedDependencies, result)
		})
	}
}

func TestUpdateService_IsThereAConflictBetweenIncludedAndExcludedPackages(t *testing.T) {
	tests := []struct {
		name             string
		config           *config.CLIConfig
		pubspec          *models.Pubspec
		expectedConflict bool
	}{
		{
			name: "no conflict",
			config: &config.CLIConfig{
				IncludePackages: &[]string{"http"},
				ExcludePackages: &[]string{"flutter_svg"},
			},
			pubspec:          &models.Pubspec{},
			expectedConflict: false,
		},
		{
			name: "conflict exists",
			config: &config.CLIConfig{
				IncludePackages: &[]string{"http", "flutter_svg"},
				ExcludePackages: &[]string{"flutter_svg"},
			},
			pubspec:          &models.Pubspec{},
			expectedConflict: true,
		},
		{
			name:             "no include or exclude packages",
			config:           &config.CLIConfig{},
			pubspec:          &models.Pubspec{},
			expectedConflict: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &UpdateService{
				Config: tt.config,
			}
			result := service.isThereAConflictBetweenIncludedAndExcludedPackages(tt.pubspec)
			assert.Equal(t, tt.expectedConflict, result)
		})
	}
}

func TestCleanupVersionString(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"^1.2.3", "1.2.3"},
		{"~1.2.3", "1.2.3"},
		{">=1.2.3", "1.2.3"},
		{"<2.0.0", "2.0.0"},
		{"1.2.3", "1.2.3"},
		{"^1.2.3 <2.0.0", "1.2.3 2.0.0"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := cleanupVersionString(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestUpdateService_ProduceSliceOfDependencyUpdates(t *testing.T) {
	tests := []struct {
		name                  string
		dependenciesToUpdate  []string
		dependencyDataFromAPI []*models.PackageWrapper
		mockPubspec           *models.Pubspec
		expected              []models.DependencyUpdate
	}{
		{
			name:                 "update needed",
			dependenciesToUpdate: []string{"http"},
			dependencyDataFromAPI: []*models.PackageWrapper{
				{
					Name: "http",
					LatestVersion: models.Package{
						Version: "0.13.5",
					},
				},
			},
			mockPubspec: &models.Pubspec{
				Dependencies: map[string]any{
					"http": "^0.13.3",
				},
			},
			expected: []models.DependencyUpdate{
				{
					Name:           "http",
					CurrentVersion: "0.13.3",
					LatestVersion:  "0.13.5",
				},
			},
		},
		{
			name:                 "no update needed - same version",
			dependenciesToUpdate: []string{"http"},
			dependencyDataFromAPI: []*models.PackageWrapper{
				{
					Name: "http",
					LatestVersion: models.Package{
						Version: "0.13.3",
					},
				},
			},
			mockPubspec: &models.Pubspec{
				Dependencies: map[string]any{
					"http": "^0.13.3",
				},
			},
			expected: nil,
		},
		{
			name:                 "multiple dependencies",
			dependenciesToUpdate: []string{"http", "path"},
			dependencyDataFromAPI: []*models.PackageWrapper{
				{
					Name: "http",
					LatestVersion: models.Package{
						Version: "0.13.5",
					},
				},
				{
					Name: "path",
					LatestVersion: models.Package{
						Version: "1.8.3",
					},
				},
			},
			mockPubspec: &models.Pubspec{
				Dependencies: map[string]any{
					"http": "^0.13.3",
					"path": "^1.8.0",
				},
			},
			expected: []models.DependencyUpdate{
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockPubspecParser := &parsers.MockPubspecParser{
				ParseFunc: func() (*models.Pubspec, error) {
					return tt.mockPubspec, nil
				},
			}

			service := &UpdateService{
				PubspecParser: mockPubspecParser,
			}

			// Act
			result := service.produceSliceOfDependencyUpdates(tt.dependenciesToUpdate, tt.dependencyDataFromAPI)

			// Sort both slices to ensure consistent comparison
			sortDependencyUpdates := func(updates []models.DependencyUpdate) {
				sort.Slice(updates, func(i, j int) bool {
					return updates[i].Name < updates[j].Name
				})
			}

			sortDependencyUpdates(result)
			sortDependencyUpdates(tt.expected)

			// Assert
			assert.Equal(t, tt.expected, result)
		})
	}
}
