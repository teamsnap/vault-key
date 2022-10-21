package vault

import (
	"context"
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
var secretEngine string

func TestGetSecretFromVault(t *testing.T) {
	secretKey, secretValue, secretEngine = "myKey", "myValue", "kv/data/get/foo"

	cluster := createTestVault(t)
	defer cluster.Cleanup()

	rootVaultClient := cluster.Cores[0].Client
	vc := &vaultClient{
		config: &config{},
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
		secrets, err := vc.SecretFromVault(secretEngine)
		if err != nil {
			t.Errorf("get secret from vault, %s", err)
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
		path := "kv/metadata/get/foo"

		version, err := vc.SecretVersionFromVault(path)
		if err != nil {
			t.Errorf("get versioned secret from vault, %s", err)
		}

		is.Equal(version, int64(1))
	}
}

func tesetInvalidPath(vc *vaultClient) func(*testing.T) {
	return func(t *testing.T) {
		is := is.New(t)
		path := "foo"
		_, err := vc.SecretFromVault(path)

		is.True(err != nil)
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

	core := cluster.Cores[0].Core
	hashivault.TestWaitActive(t, core)

	// Setup required secrets, policies, etc.
	_, err := cluster.Cores[0].Client.Logical().Write(secretEngine, secrets)
	if err != nil {
		t.Fatal(err)
	}

	return cluster
}
