package main_test

import (
	"os"
	"testing"
)

func MainTest(t *testing.T) {
	os.Setenv("VAULT_SECRET", "secret/test/data/dotenv")
	os.Setenv("K8S_NAMESPACE", "test")

	// cluster := createTestVault(t)
	// defer cluster.Cleanup()

	// rootVaultClient := cluster.Cores[0].Client
	// vc := &vaultClient{
	// 	config: githubConfig(t),
	// 	ctx:    context.Background(),
	// 	client: rootVaultClient,
	// }

	// t.Run("valid client", testValidClient(vc))
}

// func testValidClient(vc *vaultClient) func(*testing.T) {
// 	return func(t *testing.T) {
// 		path := "secret/test/data/"
// 		secrets, err := vc.SecretFromVault(path)
// 		if err != nil {
// 			t.Errorf("get secret from vault, %v", err)
// 		}

// 		for k, v := range secrets {
// 			if k != secretKey {
// 				t.Errorf("Actual: %q, Expected: %q", k, secretKey)
// 			}
// 			if v != secretValue {
// 				t.Errorf("Actual: %q, Expected: %q", v, secretValue)
// 			}
// 		}
// 	}
// }

// func createTestVault(t *testing.T) *hashivault.TestCluster {
// 	t.Helper()

// 	coreConfig := &hashivault.CoreConfig{}
// 	cluster := hashivault.NewTestCluster(t, coreConfig, &hashivault.TestClusterOptions{
// 		HandlerFunc: vaulthttp.Handler,
// 	})
// 	cluster.Start()

// 	secrets := map[string]interface{}{
// 		"data":     map[string]interface{}{secretKey: secretValue},
// 		"metadata": map[string]interface{}{"version": 2},
// 	}
// 	// Setup required secrets, policies, etc.
// 	_, err := cluster.Cores[0].Client.Logical().Write("secret/test/data", secrets)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	return cluster
// }
