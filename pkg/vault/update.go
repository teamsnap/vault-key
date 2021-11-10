package vault

import (
	"fmt"

	"github.com/hashicorp/vault/api"
)

// Update takes a given key for an engine, and modifies its value in vault.
func (vc *vaultClient) Update(engine, key, value string) (*api.Secret, error) {
	vc.tracer.trace(fmt.Sprintf("%s/Update", vc.config.tracePrefix))

	data, err := vc.SecretFromVault(engine)
	if err != nil {
		return nil, fmt.Errorf("failed to verify engine at %s: %w", engine, err)
	}

	if _, ok := data[key]; !ok {
		return nil, fmt.Errorf("missing key: %s for secret at %s", key, engine)
	} else {
		data[key] = value
	}

	secret, err := vc.write(engine, data)
	if err != nil {
		return secret, fmt.Errorf("failed to update secret: %w", err)
	}

	return secret, nil
}
