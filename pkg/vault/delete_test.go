package vault

import (
	"context"
	"testing"

	"github.com/matryer/is"
)

func TestDelete(t *testing.T) {
	secretKey, secretValue, secretEngine = "deletable-key", "value", "kv/data/delete/foo"

	cluster := createTestVault(t)
	defer cluster.Cleanup()

	rootVaultClient := cluster.Cores[0].Client
	vc := &vaultClient{
		config: githubConfig(t),
		ctx:    context.Background(),
		client: rootVaultClient,
	}
	vc.tracer = vc

	t.Run("delete secret that does not exist", delete_new(vc))
	t.Run("delete secret", delete_existing(vc))
	t.Run("delete missing path", delete_missingPath(vc))
}

func delete_new(vc *vaultClient) func(*testing.T) {
	return func(t *testing.T) {
		is := is.New(t)

		engine, k := secretEngine, "new-key"
		_, err := vc.delete(engine, k)
		is.True(err != nil)
	}
}

func delete_existing(vc *vaultClient) func(*testing.T) {
	return func(t *testing.T) {
		is := is.New(t)
		engine, k := secretEngine, secretKey

		_, err := vc.delete(engine, k)
		is.NoErr(err)

		secret, err := vc.SecretFromVault(engine)
		is.NoErr(err)
		_, present := secret[k]
		is.Equal(false, present)
	}
}

func delete_missingPath(vc *vaultClient) func(*testing.T) {
	return func(t *testing.T) {
		is := is.New(t)

		engine, k := "kv/data/delete/missing/foo", secretKey
		_, err := vc.delete(engine, k)
		is.True(err != nil)
	}
}
