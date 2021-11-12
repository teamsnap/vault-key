package vault

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/matryer/is"
)

func TestUpdate(t *testing.T) {
	secretKey, secretValue, secretEngine = "existing-key", "foo", "kv/data/update/foo"

	cluster := createTestVault(t)
	defer cluster.Cleanup()

	rootVaultClient := cluster.Cores[0].Client
	vc := &vaultClient{
		config: githubConfig(t),
		ctx:    context.Background(),
		client: rootVaultClient,
	}
	vc.tracer = vc

	t.Run("update secret that does not exist", update_new(vc))
	t.Run("update secret", update_existing(vc))
	t.Run("update missing path", update_missingPath(vc))
}

func update_new(vc *vaultClient) func(*testing.T) {
	return func(t *testing.T) {
		is := is.New(t)

		engine, k, v := secretEngine, "update-new-key", "new-value"
		_, err := vc.update(engine, k, v)

		is.True(err != nil)
	}
}

func update_existing(vc *vaultClient) func(*testing.T) {
	return func(t *testing.T) {
		is := is.New(t)
		engine, k, v := secretEngine, secretKey, secretValue

		version, err := vc.SecretVersionFromVault("kv/metadata/update/foo")
		is.NoErr(err)

		secret, err := vc.update(engine, k, v)
		is.NoErr(err)

		currentVersion, err := secret.Data["version"].(json.Number).Int64()
		is.NoErr(err)

		is.Equal(version+1, currentVersion)
	}
}

func update_missingPath(vc *vaultClient) func(*testing.T) {
	return func(t *testing.T) {
		is := is.New(t)

		engine, k, v := "kv/data/update/missing", secretKey, secretValue
		_, err := vc.update(engine, k, v)

		is.True(err != nil)
	}
}
