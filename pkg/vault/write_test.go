package vault

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/matryer/is"
)

var secrets map[string]string

func TestWriteSecret(t *testing.T) {
	secrets = map[string]string{"existing-key": "foo"}
	secretEngine = "kv/data/update/foo"

	cluster := createTestVault(t)
	defer cluster.Cleanup()

	rootVaultClient := cluster.Cores[0].Client
	vc := &vaultClient{
		config: githubConfig(t),
		ctx:    context.Background(),
		client: rootVaultClient,
	}
	vc.tracer = vc

	t.Run("write new secret when secret exists", write_existing(vc))
	t.Run("write new secret", write_new(vc))
	t.Run("write new secret with a missing path", write_missingPath(vc))
}

func write_new(vc *vaultClient) func(*testing.T) {
	return func(t *testing.T) {
		is := is.New(t)
		secrets["new-key"] = "new new"
		secrets["fresh-key"] = "freshy new"
		_, err := vc.write(secretEngine, secrets)
		is.NoErr(err)

		datum, err := vc.SecretFromVault(secretEngine)
		is.NoErr(err)
		is.Equal(len(datum), len(secrets))
	}
}

func write_existing(vc *vaultClient) func(*testing.T) {
	return func(t *testing.T) {
		is := is.New(t)

		version, err := vc.SecretVersionFromVault("kv/metadata/update/foo")
		is.NoErr(err)

		secret, err := vc.write(secretEngine, secrets)
		is.NoErr(err)

		currentVersion, err := secret.Data["version"].(json.Number).Int64()
		is.NoErr(err)

		is.Equal(version+1, currentVersion)
	}
}

func write_missingPath(vc *vaultClient) func(*testing.T) {
	return func(t *testing.T) {
		is := is.New(t)
		engine := "kv/data/updawrite/missingte/foo"

		_, err := vc.write(engine, secrets)
		is.NoErr(err)
	}
}
