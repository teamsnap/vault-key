package vault

import (
	"fmt"

	"github.com/hashicorp/vault/api"
)

// Create takes a given key for an engine, and adds a new key/value pair in vault.
func (vc *vaultClient) Create(engine, key, value string) (*api.Secret, error) {
	vc.tracer.trace(fmt.Sprintf("%s/Create", vc.config.tracePrefix))

	data, err := vc.SecretFromVault(engine)
	if err != nil {
		return nil, fmt.Errorf("failed to verify engine at %s: %w", engine, err)
	}

	if _, ok := data[key]; ok {
		return nil, fmt.Errorf("key: %s for secret at %s already exists", key, engine)
	} else {
		data[key] = value
	}

	secret, err := vc.write(engine, data)
	if err != nil {
		return secret, fmt.Errorf("failed to create secret for %s: %w", key, err)
	}

	return secret, nil
}
