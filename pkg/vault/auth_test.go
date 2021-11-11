package vault

import (
	"context"
	"errors"
	"os"
	"testing"
)

func gcpConfig(t *testing.T) *config {
	os.Setenv("GITHUB_OAUTH_TOKEN", "")
	os.Setenv("GCLOUD_PROJECT", "value")
	os.Setenv("FUNCTION_IDENTITY", "value")
	os.Setenv("GCP_AUTH_PATH", "value")
	os.Setenv("VAULT_ROLE", "role")

	cfg, err := loadVaultEnvironment()
	if err != nil {
		t.Fatal(err)
	}

	return cfg
}

func githubConfig(t *testing.T) *config {
	os.Setenv("GITHUB_OAUTH_TOKEN", "token")
	cfg, err := loadVaultEnvironment()
	if err != nil {
		t.Fatal(err)
	}

	return cfg
}

func TestGoogleVaultClient(t *testing.T) {
	tests := []struct {
		name     string
		expected interface{}
	}{
		{
			name:     "valid auth client",
			expected: nil,
		}, {
			name:     "invalid auth client",
			expected: errors.New("initialze client: getting vault api token from client: getting new iam service: google: could not find default credentials"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewVaultClient(context.Background(), gcpConfig(t))
			_, typ_err := tt.expected.(error)
			if tt.expected != nil && err != nil && !typ_err {
				t.Errorf("Actual: %q. Expected: %q", err, tt.expected)
			}

			if err == nil {
				if tt.expected != nil {
					t.Errorf("Actual: %q. Expected: %q", err, tt.expected)
				}
			}
		})
	}
}

func TestGitHubVaultClient(t *testing.T) {
	tests := []struct {
		name     string
		expected interface{}
	}{
		{
			name:     "valid auth client",
			expected: nil,
		}, {
			name:     "invalid auth client",
			expected: errors.New("super error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewVaultClient(context.Background(), githubConfig(t))
			if err != nil && tt.expected != nil {
				_, typ_err := tt.expected.(error)
				if err != nil && !typ_err {
					t.Errorf("Actual: %q. Expected: %q", err, tt.expected)
				}
			}

			if err == nil {
				if tt.expected != nil {
					t.Errorf("Actual: %q. Expected: %q", err, tt.expected)
				}
			}
		})
	}
}
