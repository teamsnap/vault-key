package vault

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/matryer/is"
)

func TestWriteSecret(t *testing.T) {
	secretKey, secretValue = "existing-key", "existing-value"

	cluster := createTestVault(t)
	defer cluster.Cleanup()

	rootVaultClient := cluster.Cores[0].Client
	vc := &vaultClient{
		config: githubConfig(t),
		ctx:    context.Background(),
		client: rootVaultClient,
	}

	t.Run("write new secret", write_new(vc))
	t.Run("write new secret when secret exits", write_existing(vc))
	t.Run("write new secret with a missing path", write_missingPath(vc))
}

func write_new(vc *vaultClient) func(*testing.T) {
	return func(t *testing.T) {
		is := is.New(t)

		engine, k, v := "kv/data/create/foo", "new-key", "new-value"
		_, err := vc.write(engine, k, v)

		is.NoErr(err)
	}
}

func write_existing(vc *vaultClient) func(*testing.T) {
	return func(t *testing.T) {
		is := is.New(t)

		engine, k, v := "kv/data/create/foo", secretKey, secretValue
		version, err := vc.SecretVersionFromVault("kv/metadata/create/foo")
		is.NoErr(err)

		secret, err := vc.write(engine, k, v)
		is.NoErr(err)

		currentVersion, err := secret.Data["version"].(json.Number).Int64()
		is.NoErr(err)

		is.Equal(version+1, currentVersion)
	}
}

func write_missingPath(vc *vaultClient) func(*testing.T) {
	return func(t *testing.T) {
		is := is.New(t)

		engine, k, v := "", "new-key", "new-value"
		_, err := vc.write(engine, k, v)

		is.Equal(err != nil, true)
	}
}
