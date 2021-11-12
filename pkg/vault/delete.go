package vault

import (
	"fmt"

	"github.com/hashicorp/vault/api"
)

// Delete takes a given key for an engine, and removes the key/value pair from vault.
func (vc *vaultClient) delete(engine, key string) (*api.Secret, error) {
	vc.tracer.trace(fmt.Sprintf("%s/delete", vc.config.tracePrefix))

	res, err := vc.SecretFromVault(engine)
	if err != nil {
		return nil, fmt.Errorf("failed to verify engine %s: %w", engine, err)
	}

	if _, ok := res[key]; !ok {
		return nil, fmt.Errorf("key: %s does not exist for engine at %s", key, engine)
	}

	secret, err := vc.client.Logical().Delete(engine + "/" + key)
	if err != nil {
		return secret, fmt.Errorf("failed to delete key %s at %s:%w", key, engine, err)
	}

	return secret, nil
}
