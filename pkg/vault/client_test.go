package vault

import (
	"context"
	"errors"
	"testing"

	vaulthttp "github.com/hashicorp/vault/http"
	hashivault "github.com/hashicorp/vault/vault"
)

type GitHubMockAuthClient struct{}
type GoogleMockAuthClient struct{}

func (m *GoogleMockAuthClient) GetVaultToken(gc *vaultClient) (string, error)  { return "token", nil }
func (m *GitHubMockAuthClient) GetVaultToken(ghc *vaultClient) (string, error) { return "token", nil }

var googleAuthClient AuthClient
var githubAuthClient AuthClient
var gc *config
var ghc *config
var ctx context.Context
var secretKey string
var secretValue string

func init() {
	gc = &config{
		project:        "test",
		googleAuth:     true,
		serviceAccount: "none",
		traceEnabled:   false,
		tracePrefix:    "test",
		vaultRole:      "read",
	}

	ghc = &config{
		project:        "test",
		googleAuth:     false,
		githubAuth:     true,
		githubToken:    "1234",
		serviceAccount: "none",
		traceEnabled:   false,
		tracePrefix:    "test",
		vaultRole:      "read",
	}

	secretKey = "myKey"
	secretValue = "myValue"

	googleAuthClient = &GoogleMockAuthClient{}
	githubAuthClient = &GitHubMockAuthClient{}
	ctx = context.Background()
}

func TestGoogleVaultClient(t *testing.T) {
	tests := []struct {
		name     string
		auth     AuthClient
		expected error
	}{
		{"valid auth client", googleAuthClient, nil},
		{"invalid auth client", NewAuthClient(), errors.New("initialze client: getting vault api token from client: generating signed jwt, sigining jwt: Post")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewVaultClient(ctx, tt.auth, gc)
			if err != nil && tt.expected != nil {
				if err.Error() != tt.expected.Error() {
					// This error happens in travis ci, not the iam error.
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
		auth     AuthClient
		expected error
	}{
		{"valid auth client", githubAuthClient, nil},
		{"invalid auth client", NewAuthClient(), errors.New("initialze client: initialze client: getting vault api token from client: logging into vault with github:Put /v1/auth/github/login: unsupported protocol scheme \"\"")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewVaultClient(ctx, tt.auth, ghc)
			if err != nil && tt.expected != nil {
				if err.Error() != tt.expected.Error() {
					// This error happens in travis ci, not the iam error.
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

func TestGoogleGetSecretFromVault(t *testing.T) {
	cluster := createTestVault(t)
	defer cluster.Cleanup()

	rootVaultClient := cluster.Cores[0].Client
	vc := &vaultClient{
		authClient: googleAuthClient,
		config:     gc,
		ctx:        ctx,
		client:     rootVaultClient,
	}

	t.Run("valid client", testValidClient(vc))
	t.Run("invalid path", tesetInvalidPath(vc))
	t.Run("versioned secrets", testVersionedSecrets(vc))
}

func TestGitHubGetSecretFromVault(t *testing.T) {
	cluster := createTestVault(t)
	defer cluster.Cleanup()

	rootVaultClient := cluster.Cores[0].Client
	vc := &vaultClient{
		authClient: githubAuthClient,
		config:     ghc,
		ctx:        ctx,
		client:     rootVaultClient,
	}

	t.Run("valid client", testValidClient(vc))
	t.Run("invalid path", tesetInvalidPath(vc))
	t.Run("versioned secrets", testVersionedSecrets(vc))
}

func testValidClient(vc *vaultClient) func(*testing.T) {
	return func(t *testing.T) {
		path := "secret/test/data/"
		secrets, err := vc.GetSecretFromVault(path)
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
		path := "secret/test/data/"
		version, err := vc.GetSecretVersionFromVault(path)
		if err != nil {
			t.Errorf("get versioned secret from vault, %v", err)
		}

		if version != 2 {
			t.Errorf("Actual: %v, Expected: %v", version, 2)
		}

	}
}

func tesetInvalidPath(vc *vaultClient) func(*testing.T) {
	return func(t *testing.T) {
		path := "foo"
		_, err := vc.GetSecretFromVault(path)

		expected := errors.New("secret values returned from Vault are <nil> for foo")
		if err.Error() != expected.Error() {
			t.Errorf("Expected invalid path to raise error: %v", err)
		}
	}
}

func createTestVault(t *testing.T) *hashivault.TestCluster {
	t.Helper()

	coreConfig := &hashivault.CoreConfig{}
	cluster := hashivault.NewTestCluster(t, coreConfig, &hashivault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})
	cluster.Start()

	secrets := map[string]interface{}{
		"data":     map[string]interface{}{secretKey: secretValue},
		"metadata": map[string]interface{}{"version": 2},
	}
	// Setup required secrets, policies, etc.
	_, err := cluster.Cores[0].Client.Logical().Write("secret/test/data", secrets)
	if err != nil {
		t.Fatal(err)
	}

	return cluster
}
