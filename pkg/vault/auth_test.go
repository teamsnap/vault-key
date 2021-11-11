package vault

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/matryer/is"
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

	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "test/default_credentials.json")

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
	is := is.New(t)
	loginServer := vaultLoginServer()
	defer loginServer.Close()

	os.Setenv("VAULT_ADDR", loginServer.URL)

	_, err := NewVaultClient(context.Background(), githubConfig(t))
	is.NoErr(err)
}

func vaultLoginServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/login") {
			json.NewEncoder(w).Encode(api.Secret{
				Auth: &api.SecretAuth{ClientToken: "vault-test-token"},
			})
			return
		}

		switch r.Method {
		case http.MethodGet:
			json.NewEncoder(w).Encode(
				&api.Secret{
					Data: map[string]interface{}{
						"data":            map[string]interface{}{"my-key": "bar"},
						"current_version": 1,
					},
				},
			)
		case http.MethodPost, http.MethodPut:
			var incoming map[string]interface{}
			json.NewDecoder(r.Body).Decode(&incoming)
		}
	}))
}
