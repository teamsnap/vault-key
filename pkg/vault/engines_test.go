package vault

import (
	"context"
	"testing"

	"github.com/matryer/is"
)

func TestEngines(t *testing.T) {
	secretKey, secretValue, secretEngine = "myKey", "myValue", "kv/data/delete/foo"

	secrets := map[string]interface{}{
		"data": map[string]interface{}{secretKey: secretValue},
	}

	cluster := createTestVault(t)
	defer cluster.Cleanup()

	_, err := cluster.Cores[0].Client.Logical().Write("kv/data/delete/bar", secrets)
	if err != nil {
		t.Fatal(err)
	}

	rootVaultClient := cluster.Cores[0].Client
	vc := &vaultClient{
		config: githubConfig(t),
		ctx:    context.Background(),
		client: rootVaultClient,
	}
	vc.tracer = vc

	t.Run("engine listings", testEngineListing(vc))
}

func testEngineListing(vc *vaultClient) func(*testing.T) {
	return func(t *testing.T) {
		is := is.New(t)
		expected := []string{"bar", "foo"}
		path := "kv/metadata/delete"
		engines, err := vc.enginesFromVault(path)
		if err != nil {
			t.Errorf("list engines from vault, %v", err)
		}

		if len(engines) < 1 {
			t.Errorf("no engines returned")
		}

		is.Equal(engines, expected)

	}
}
