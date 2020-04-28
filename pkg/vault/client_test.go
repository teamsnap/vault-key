package vault

import (
	"context"
	"errors"
	"net"
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/vault"
)

type MockAuthClient struct{}

func (m *MockAuthClient) GetVaultToken(c *vaultClient) (string, error) { return "token", nil }

var authClient AuthClient
var c *config
var ctx context.Context
var secretKey string
var secretValue string

func init() {
	c = &config{
		project:        "test",
		serviceAccount: "none",
		traceEnabled:   false,
		tracePrefix:    "test",
		vaultRole:      "read",
	}

	secretKey = "myKey"
	secretValue = "myValue"

	authClient = &MockAuthClient{}
	ctx = context.Background()
}

func TestNewVaultClient(t *testing.T) {
	tests := []struct {
		name     string
		auth     AuthClient
		expected error
	}{
		{"valid auth client", authClient, nil},
		{"invalid auth client", NewAuthClient(), errors.New("initialze client: getting new iam service: google: could not find default credentials")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewVaultClient(ctx, tt.auth, c)
			if err != nil && tt.expected != nil {
				if err.Error() != tt.expected.Error() {
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

func TestGetSecretFromVault(t *testing.T) {
	ln, client := createTestVault(t)
	defer ln.Close()
	vc := &vaultClient{
		authClient: authClient,
		config:     c,
		ctx:        ctx,
		client:     client,
	}

	t.Run("valid client", testValidClient(vc))
	t.Run("invalid path", tesetInvalidPath(vc))
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

func createTestVault(t *testing.T) (net.Listener, *api.Client) {
	t.Helper()

	// Create an in-memory, unsealed core (the "backend", if you will).
	core, keyShares, rootToken := vault.TestCoreUnsealed(t)
	_ = keyShares

	// Start an HTTP server for the core.
	ln, addr := http.TestServer(t, core)

	// Create a client that talks to the server, initially authenticating with
	// the root token.
	conf := api.DefaultConfig()
	conf.Address = addr

	client, err := api.NewClient(conf)
	if err != nil {
		t.Fatal(err)
	}
	client.SetToken(rootToken)

	secrets := map[string]interface{}{
		"data": map[string]interface{}{secretKey: secretValue},
	}
	// Setup required secrets, policies, etc.
	_, err = client.Logical().Write("secret/test/data", secrets)
	if err != nil {
		t.Fatal(err)
	}

	return ln, client
}
