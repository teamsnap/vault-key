package vault

import (
	"context"
	"testing"

	"github.com/matryer/is"
)

func TestDelete(t *testing.T) {
	secretKey, secretValue = "deletable-key", "foo"

	cluster := createTestVault(t)
	defer cluster.Cleanup()

	rootVaultClient := cluster.Cores[0].Client
	vc := &vaultClient{
		config: githubConfig(t),
		ctx:    context.Background(),
		client: rootVaultClient,
	}

	t.Run("delete secret that does not exist", delete_new(vc))
	t.Run("delete secret", delete_existing(vc))
	t.Run("delete missing path", delete_missingPath(vc))
}

func delete_new(vc *vaultClient) func(*testing.T) {
	return func(t *testing.T) {
		is := is.New(t)

		engine, k := "kv/data/foo", "new-key"
		_, err := vc.Delete(engine, k)
		is.Equal(err != nil, true)
	}
}

func delete_existing(vc *vaultClient) func(*testing.T) {
	return func(t *testing.T) {
		is := is.New(t)

		engine, k := "kv/data/foo", secretKey
		_, err := vc.Delete(engine, k)
		is.NoErr(err)

		secret, err := vc.SecretFromVault(engine + "/" + k)
		is.Equal(err != nil, true)
		_, ok := secret[k]
		is.Equal(false, ok)
	}
}

func delete_missingPath(vc *vaultClient) func(*testing.T) {
	return func(t *testing.T) {
		is := is.New(t)

		engine, k := "kv/data/delete/foo", secretKey
		_, err := vc.Delete(engine, k)
		is.Equal(err != nil, true)
	}
}