package vault

import (
	"context"
	"testing"

	"github.com/matryer/is"
)

func TestCreateSecret(t *testing.T) {
	secretKey, secretValue, secretEngine = "existing-key", "foo", "kv/data/create/foo"

	cluster := createTestVault(t)
	defer cluster.Cleanup()

	rootVaultClient := cluster.Cores[0].Client
	vc := &vaultClient{
		config: githubConfig(t),
		ctx:    context.Background(),
		client: rootVaultClient,
	}
	vc.tracer = vc

	t.Run("create new secret", create_new(vc))
	t.Run("create new secret when secret exists", create_existing(vc))
	t.Run("create new secret with a missing path", create_missingPath(vc))
}

func create_new(vc *vaultClient) func(*testing.T) {
	return func(t *testing.T) {
		is := is.New(t)

		engine, k, v := secretEngine, "new-key", secretValue
		_, err := vc.create(engine, k, v)

		is.NoErr(err)
	}
}

func create_existing(vc *vaultClient) func(*testing.T) {
	return func(t *testing.T) {
		is := is.New(t)

		engine, k, v := secretEngine, secretKey, secretValue
		version, err := vc.SecretVersionFromVault("kv/metadata/create/foo")
		is.NoErr(err)

		_, err = vc.create(engine, k, v)
		is.True(err != nil)

		currentVersion, err := vc.SecretVersionFromVault("kv/metadata/create/foo")
		is.NoErr(err)

		is.Equal(version, currentVersion)
	}
}

func create_missingPath(vc *vaultClient) func(*testing.T) {
	return func(t *testing.T) {
		is := is.New(t)

		engine, k, v := "kv/data/create/missing/foo", secretKey, secretValue
		_, err := vc.create(engine, k, v)

		is.True(err != nil)
	}
}
