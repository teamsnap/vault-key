package vault

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/api"
	"github.com/matryer/is"

	kv "github.com/hashicorp/vault-plugin-secrets-kv"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/sdk/logical"
	hashivault "github.com/hashicorp/vault/vault"
)

var secretKey string
var secretValue string

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
		expected error
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
			if err != nil && tt.expected != nil {
				if err.Error() != tt.expected.Error() {
					// This error happens in ci, not the iam error.
					travisError := errors.New("initialze client: generating signed jwt, sigining jwt: Post")
					if err.Error() != travisError.Error() {
						t.Errorf("Actual: %q. Expected: %q", err, tt.expected)
					}
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

func TestGitHubVaultClient(t *testing.T) {
	tests := []struct {
		name     string
		expected error
	}{
		{
			name:     "valid auth client",
			expected: nil,
		}, {
			name:     "invalid auth client",
			expected: errors.New("initialze client: getting vault api token from client: logging into vault with github:Put \"/v1/auth/github/login\": unsupported protocol scheme \"\""),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewVaultClient(context.Background(), githubConfig(t))
			if err != nil && tt.expected != nil {
				if err.Error() != tt.expected.Error() {
					// This error happens in ci, not the iam error.
					ciError := errors.New("initialze client: getting vault api token from client: logging into vault with github:Put /v1/auth/github/login: unsupported protocol scheme \"\"")
					if err.Error() != ciError.Error() {
						t.Errorf("Actual: %q. Expected: %q", err, tt.expected)
					}
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

func TestGoogleGetSecretFromVault(t *testing.T) {
	secretKey = "myGcpKey"
	secretValue = "myGcpValue"

	cluster := createTestVault(t)
	defer cluster.Cleanup()

	rootVaultClient := cluster.Cores[0].Client
	vc := &vaultClient{
		config: gcpConfig(t),
		ctx:    context.Background(),
		client: rootVaultClient,
	}
	vc.tracer = vc

	t.Run("valid client", testValidClient(vc))
	t.Run("invalid path", tesetInvalidPath(vc))
	t.Run("versioned secrets", testVersionedSecrets(vc))
}

func TestGitHubGetSecretFromVault(t *testing.T) {
	secretKey = "myGithubKey"
	secretValue = "myGithubValue"

	cluster := createTestVault(t)
	defer cluster.Cleanup()

	rootVaultClient := cluster.Cores[0].Client
	vc := &vaultClient{
		config: githubConfig(t),
		ctx:    context.Background(),
		client: rootVaultClient,
	}
	vc.tracer = vc

	t.Run("valid client", testValidClient(vc))
	t.Run("invalid path", tesetInvalidPath(vc))
	t.Run("versioned secrets", testVersionedSecrets(vc))
}

func testValidClient(vc *vaultClient) func(*testing.T) {
	return func(t *testing.T) {
		path := "kv/data/foo"
		secrets, err := vc.SecretFromVault(path)
		if err != nil {
			t.Errorf("get secret from vault, %v", err)
		}

		for k, v := range secrets {
			if k != secretKey {
				t.Errorf("Actual: %q, Expected: %q", k, secretKey)
			}
			if v != secretValue {
				t.Errorf("Actual: %q, Expected: %q", v, secretValue)
			}
		}
	}
}

func testVersionedSecrets(vc *vaultClient) func(*testing.T) {
	return func(t *testing.T) {
		is := is.New(t)
		path := "kv/metadata/foo"

		version, err := vc.SecretVersionFromVault(path)
		if err != nil {
			t.Errorf("get versioned secret from vault, %v", err)
		}

		is.Equal(version, int64(1))
	}
}

func tesetInvalidPath(vc *vaultClient) func(*testing.T) {
	return func(t *testing.T) {
		is := is.New(t)
		path := "foo"
		_, err := vc.SecretFromVault(path)

		is.Equal(err != nil, true)
	}
}

func createTestVault(t *testing.T) *hashivault.TestCluster {
	t.Helper()

	coreConfig := &hashivault.CoreConfig{
		LogicalBackends: map[string]logical.Factory{
			"kv": kv.Factory,
		},
	}

	cluster := hashivault.NewTestCluster(t, coreConfig, &hashivault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
		Logger:      hclog.NewNullLogger(),
	})

	cluster.Start()

	secrets := map[string]interface{}{
		"data": map[string]interface{}{secretKey: secretValue},
	}

	// Create KV V2 mount
	if err := cluster.Cores[0].Client.Sys().Mount("kv", &api.MountInput{
		Type: "kv",
		Options: map[string]string{
			"version": "2",
		},
	}); err != nil {
		t.Fatal(err)
	}

	// Setup required secrets, policies, etc.
	_, err := cluster.Cores[0].Client.Logical().Write("kv/data/foo", secrets)
	if err != nil {
		t.Fatal(err)
	}

	return cluster
}
