package services

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

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
	service := NewAPIService()
	if service == nil {
		t.Fatal("expected service to be non-nil, got nil")
	}
	if service.Client == nil {
		t.Fatal("expected client to be non-nil, got nil")
	}
}

func TestGetSDKRelease(t *testing.T) {
	tests := []struct {
		name           string
		mockResp       string
		mockStatusCode int
		mockErr        error
		want           *models.SDKReleaseWrapper
		wantErr        bool
	}{
		{
			name: "successful response",
			mockResp: `{
				"current_release": {
					"beta": "beta-hash",
					"stable": "stable-hash"
				},
				"releases": [
					{
						"hash": "test-hash",
						"channel": "stable",
						"version": "2.15.0",
						"dart_sdk_version": "3.0.0"
					}
				]
			}`,
			mockStatusCode: http.StatusOK,
			mockErr:        nil,
			want: &models.SDKReleaseWrapper{
				CurrentRelease: models.SDKReleaseHashes{
					Beta:   "beta-hash",
					Stable: "stable-hash",
				},
				Releases: []models.SDKRelease{
					{
						Hash:              "test-hash",
						Channel:           "stable",
						FlutterSDKVersion: "2.15.0",
						DartSDKVersion:    "3.0.0",
					},
				},
			},
			wantErr: false,
		},
		{
			name:           "empty response",
			mockResp:       `{}`,
			mockStatusCode: http.StatusOK,
			mockErr:        nil,
			want: &models.SDKReleaseWrapper{
				CurrentRelease: models.SDKReleaseHashes{},
				Releases:       nil,
			},
			wantErr: false,
		},
		{
			name:           "malformed json",
			mockResp:       `{malformed`,
			mockStatusCode: http.StatusOK,
			mockErr:        nil,
			want:           nil,
			wantErr:        true,
		},
		{
			name:           "http request error",
			mockResp:       "",
			mockStatusCode: 0,
			mockErr:        errors.New("network error"),
			want:           nil,
			wantErr:        true,
		},
		{
			name:           "server error",
			mockResp:       `{"error": "Internal Server Error"}`,
			mockStatusCode: http.StatusInternalServerError,
			mockErr:        nil,
			want: &models.SDKReleaseWrapper{
				CurrentRelease: models.SDKReleaseHashes{},
				Releases:       nil,
			},
			wantErr: false, // API service doesn't check status code
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock client
			mockClient := NewTestClient(func(req *http.Request) (*http.Response, error) {
				if tt.mockErr != nil {
					return nil, tt.mockErr
				}

				// Check request URL and method
				if req.URL.String() != SDK_RELEASE_URL {
					t.Errorf("URL = %v, want %v", req.URL.String(), SDK_RELEASE_URL)
				}
				if req.Method != HTTP_METHOD {
					t.Errorf("Method = %v, want %v", req.Method, HTTP_METHOD)
				}

				return &http.Response{
					StatusCode: tt.mockStatusCode,
					Body:       io.NopCloser(bytes.NewBufferString(tt.mockResp)),
				}, nil
			})

			// Create service with mock client
			service := &APIService{
				Client: mockClient,
			}

			got, err := service.GetSDKRelease()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetSDKRelease() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetSDKRelease() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetPackage(t *testing.T) {
	tests := []struct {
		name        string
		packageName string
		mockResp    string
		mockStatus  int
		mockErr     error
		want        *models.PackageWrapper
		wantErr     bool
		expectedURL string
	}{
		{
			name:        "successful response",
			packageName: "flutter",
			mockResp: `{
				"name": "flutter",
				"latest": {
					"version": "1.0.0",
					"dependencies": {"dep1": "^1.0.0"},
					"dev_dependencies": {"dev_dep1": "^1.0.0"}
				}
			}`,
			mockStatus: http.StatusOK,
			mockErr:    nil,
			want: &models.PackageWrapper{
				Name: "flutter",
				LatestVersion: models.Package{
					Version:         "1.0.0",
					Dependencies:    map[string]any{"dep1": "^1.0.0"},
					DevDependencies: map[string]any{"dev_dep1": "^1.0.0"},
				},
			},
			wantErr:     false,
			expectedURL: "https://pub.dev/api/packages/flutter",
		},
		{
			name:        "empty package name",
			packageName: "",
			mockResp:    `{}`,
			mockStatus:  http.StatusOK,
			mockErr:     nil,
			want: &models.PackageWrapper{
				Name: "",
				LatestVersion: models.Package{
					Version:         "",
					Dependencies:    nil,
					DevDependencies: nil,
				},
			},
			wantErr:     false,
			expectedURL: "https://pub.dev/api/packages/",
		},
		{
			name:        "malformed json",
			packageName: "test",
			mockResp:    `{malformed`,
			mockStatus:  http.StatusOK,
			mockErr:     nil,
			want:        nil,
			wantErr:     true,
			expectedURL: "https://pub.dev/api/packages/test",
		},
		{
			name:        "http request error",
			packageName: "test",
			mockResp:    "",
			mockStatus:  0,
			mockErr:     errors.New("network error"),
			want:        nil,
			wantErr:     true,
			expectedURL: "https://pub.dev/api/packages/test",
		},
		{
			name:        "package not found",
			packageName: "nonexistent",
			mockResp:    `{"error": "Not found"}`,
			mockStatus:  http.StatusNotFound,
			mockErr:     nil,
			want: &models.PackageWrapper{
				Name: "",
				LatestVersion: models.Package{
					Version:         "",
					Dependencies:    nil,
					DevDependencies: nil,
				},
			},
			wantErr:     false, // API service doesn't check status code
			expectedURL: "https://pub.dev/api/packages/nonexistent",
		},
		{
			name:        "package with special characters",
			packageName: "flutter_bloc/bloc",
			mockResp:    `{"name": "flutter_bloc/bloc"}`,
			mockStatus:  http.StatusOK,
			mockErr:     nil,
			want:        &models.PackageWrapper{Name: "flutter_bloc/bloc"},
			wantErr:     false,
			expectedURL: "https://pub.dev/api/packages/flutter_bloc/bloc",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock client
			mockClient := NewTestClient(func(req *http.Request) (*http.Response, error) {
				if tt.mockErr != nil {
					return nil, tt.mockErr
				}

				// Check request URL and method
				if req.URL.String() != tt.expectedURL {
					t.Errorf("URL = %v, want %v", req.URL.String(), tt.expectedURL)
				}
				if req.Method != HTTP_METHOD {
					t.Errorf("Method = %v, want %v", req.Method, HTTP_METHOD)
				}

				return &http.Response{
					StatusCode: tt.mockStatus,
					Body:       io.NopCloser(bytes.NewBufferString(tt.mockResp)),
				}, nil
			})

			// Create service with mock client
			service := &APIService{
				Client: mockClient,
			}

			got, err := service.GetPackage(tt.packageName)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPackage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetPackage() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Test for error in creating HTTP request
func TestAPIService_RequestCreationError(t *testing.T) {
	// Create a test server that will return 200 OK but should never be called
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("Server should not be called with invalid request")
	}))
	defer server.Close()

	// Test GetSDKRelease with invalid URL
	t.Run("Invalid SDK URL", func(t *testing.T) {
		// Create a client with a transport that cannot connect
		client := &http.Client{
			Transport: &http.Transport{
				// Force a dial error by using an invalid dial context function
				DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
					return nil, errors.New("forced dial error")
				},
			},
		}

		service := &APIService{Client: client}
		_, err := service.GetSDKRelease()
		if err == nil {
			t.Error("Expected error when creating/sending request, got nil")
		}
	})

	// Test GetPackage with invalid URL
	t.Run("Invalid Package URL", func(t *testing.T) {
		// Create a client with a transport that cannot connect
		client := &http.Client{
			Transport: &http.Transport{
				// Force a dial error by using an invalid dial context function
				DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
					return nil, errors.New("forced dial error")
				},
			},
		}

		service := &APIService{Client: client}
		_, err := service.GetPackage("test")
		if err == nil {
			t.Error("Expected error when creating/sending request, got nil")
		}
	})
}

// Test for error in reading response body
func TestAPIService_ResponseBodyReadError(t *testing.T) {
	// Create a mock client that returns a reader that errors on read
	mockClient := NewTestClient(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(errorReader{}),
		}, nil
	})

	// Test GetSDKRelease with read error
	service := &APIService{Client: mockClient}

	_, err := service.GetSDKRelease()
	if err == nil {
		t.Error("Expected error when reading response body, got nil")
	}

	// Test GetPackage with read error
	_, err = service.GetPackage("test")
	if err == nil {
		t.Error("Expected error when reading response body, got nil")
	}
}
