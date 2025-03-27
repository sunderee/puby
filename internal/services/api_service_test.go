package services

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/sunderee/puby/internal/models"
)

// RoundTripFunc is a function type that implements the RoundTripper interface
type RoundTripFunc func(req *http.Request) (*http.Response, error)

// RoundTrip executes the mock round trip
func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

// NewTestClient returns a new http.Client with Transport replaced with a mock
func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: fn,
	}
}

// Mock reader that always errors
type errorReader struct{}

func (e errorReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("read error")
}

func TestNewAPIService(t *testing.T) {
	apiService := NewAPIService()
	assert.NotNil(t, apiService)
	assert.NotNil(t, apiService.Client)
	assert.Equal(t, DEFAULT_SDK_RELEASE_URL, apiService.SDKReleaseURL)
	assert.Equal(t, DEFAULT_PACKAGE_URL, apiService.PackageURL)
}

func TestGetSDKRelease(t *testing.T) {
	testCases := []struct {
		name           string
		serverResponse string
		statusCode     int
		expectError    bool
		expectedResult *models.SDKReleaseWrapper
	}{
		{
			name: "successful response",
			serverResponse: `{
				"current_release": {
					"stable": "abc123",
					"beta": "def456"
				},
				"releases": [
					{
						"hash": "abc123",
						"channel": "stable",
						"version": "3.19.0",
						"dart_sdk_version": "3.0.0"
					}
				]
			}`,
			statusCode:  http.StatusOK,
			expectError: false,
			expectedResult: &models.SDKReleaseWrapper{
				CurrentRelease: models.SDKReleaseHashes{
					Stable: "abc123",
					Beta:   "def456",
				},
				Releases: []models.SDKRelease{
					{
						Hash:              "abc123",
						Channel:           "stable",
						FlutterSDKVersion: "3.19.0",
						DartSDKVersion:    "3.0.0",
					},
				},
			},
		},
		{
			name:           "empty response",
			serverResponse: `{}`,
			statusCode:     http.StatusOK,
			expectError:    false,
			expectedResult: &models.SDKReleaseWrapper{},
		},
		{
			name:           "malformed json",
			serverResponse: `{"current_release": {`,
			statusCode:     http.StatusOK,
			expectError:    true,
			expectedResult: nil,
		},
		{
			name:           "http request error",
			serverResponse: ``,
			statusCode:     http.StatusInternalServerError,
			expectError:    true,
			expectedResult: nil,
		},
		{
			name:           "server error",
			serverResponse: ``,
			statusCode:     http.StatusBadRequest,
			expectError:    true,
			expectedResult: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup a test HTTP server
			server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				// Test the request
				assert.Equal(t, HTTP_METHOD, req.Method)
				assert.Equal(t, "GET", req.Method)

				// Send response
				rw.WriteHeader(tc.statusCode)
				fmt.Fprintln(rw, tc.serverResponse)
			}))
			defer server.Close()

			// Create a service with the test server URL
			apiService := NewAPIService()
			apiService.SDKReleaseURL = server.URL

			result, err := apiService.GetSDKRelease()

			// Check expectations
			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedResult, result)
			}
		})
	}
}

func TestGetPackage(t *testing.T) {
	testCases := []struct {
		name           string
		packageName    string
		serverResponse string
		statusCode     int
		expectError    bool
		expectedResult *models.PackageWrapper
	}{
		{
			name:        "successful response",
			packageName: "http",
			serverResponse: `{
				"name": "http",
				"latest": {
					"version": "0.13.5",
					"dependencies": {
						"http_parser": "^4.0.0"
					},
					"dev_dependencies": {
						"test": "^1.0.0"
					}
				}
			}`,
			statusCode:  http.StatusOK,
			expectError: false,
			expectedResult: &models.PackageWrapper{
				Name: "http",
				LatestVersion: models.Package{
					Version: "0.13.5",
					Dependencies: map[string]any{
						"http_parser": "^4.0.0",
					},
					DevDependencies: map[string]any{
						"test": "^1.0.0",
					},
				},
			},
		},
		{
			name:           "empty package name",
			packageName:    "",
			serverResponse: `{}`,
			statusCode:     http.StatusOK,
			expectError:    true,
			expectedResult: nil,
		},
		{
			name:           "malformed json",
			packageName:    "http",
			serverResponse: `{"name": "http", "latest": {`,
			statusCode:     http.StatusOK,
			expectError:    true,
			expectedResult: nil,
		},
		{
			name:           "http request error",
			packageName:    "http",
			serverResponse: ``,
			statusCode:     http.StatusInternalServerError,
			expectError:    true,
			expectedResult: nil,
		},
		{
			name:           "package not found",
			packageName:    "nonexistent_package",
			serverResponse: `{"error": {"message": "Package not found"}}`,
			statusCode:     http.StatusNotFound,
			expectError:    true,
			expectedResult: nil,
		},
		{
			name:           "package with special characters",
			packageName:    "package$with@special#chars",
			serverResponse: `{}`,
			statusCode:     http.StatusOK,
			expectError:    false,
			expectedResult: &models.PackageWrapper{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup a test HTTP server
			server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				// Test the request
				assert.Equal(t, HTTP_METHOD, req.Method)

				// Send response
				rw.WriteHeader(tc.statusCode)
				fmt.Fprintln(rw, tc.serverResponse)
			}))
			defer server.Close()

			// Create a service with the test server URL
			apiService := NewAPIService()
			apiService.PackageURL = server.URL + "/%s"

			result, err := apiService.GetPackage(tc.packageName)

			// Check expectations
			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedResult, result)
			}
		})
	}
}

func TestAPIService_RequestCreationError(t *testing.T) {
	testCases := []struct {
		name      string
		setupFunc func(s *APIService)
		testFunc  func(s *APIService) error
	}{
		{
			name: "Invalid SDK URL",
			setupFunc: func(s *APIService) {
				s.SDKReleaseURL = "\\invalid-url\\"
			},
			testFunc: func(s *APIService) error {
				_, err := s.GetSDKRelease()
				return err
			},
		},
		{
			name: "Invalid Package URL",
			setupFunc: func(s *APIService) {
				s.PackageURL = "\\invalid-url\\%s"
			},
			testFunc: func(s *APIService) error {
				_, err := s.GetPackage("test")
				return err
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			apiService := NewAPIService()
			tc.setupFunc(apiService)

			// Test
			err := tc.testFunc(apiService)

			// Verify an error was returned
			assert.Error(t, err)
		})
	}
}

func TestAPIService_ResponseBodyReadError(t *testing.T) {
	// Setup a test HTTP server that returns an invalid response
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Close the connection without sending a response
		hj, ok := rw.(http.Hijacker)
		if !ok {
			t.Fatal("ResponseWriter does not implement http.Hijacker")
		}
		conn, _, err := hj.Hijack()
		if err != nil {
			t.Fatal("Failed to hijack connection")
		}
		conn.Close()
	}))
	defer server.Close()

	// Test the service
	apiService := NewAPIService()
	apiService.SDKReleaseURL = server.URL

	_, err := apiService.GetSDKRelease()

	// Verify an error was returned
	assert.Error(t, err)
}
