package vault

import (
	"fmt"

	"github.com/hashicorp/vault/api"
)

// create takes a given key for an engine, and adds a new key/value pair in vault.
func (vc *vaultClient) create(engine, key, value string) (*api.Secret, error) {
	vc.tracer.trace(fmt.Sprintf("%s/create", vc.config.tracePrefix))

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

// createMulti takes a map of keys and values for an engine, and adds new key/value pairs in vault.
func (vc *vaultClient) createMulti(engine string, secrets map[string]string) (*api.Secret, error) {
	vc.tracer.trace(fmt.Sprintf("%s/createMulti", vc.config.tracePrefix))

	data, err := vc.SecretFromVault(engine)
	if err != nil {
		return nil, fmt.Errorf("failed to verify engine at %s: %w", engine, err)
	}

	for key, value := range secrets {
		if _, ok := data[key]; ok {
			return nil, fmt.Errorf("key: %s for secret at %s already exists", key, engine)
		} else {
			data[key] = value
		}
	}

	secret, err := vc.write(engine, data)
	if err != nil {
		return secret, fmt.Errorf("failed to create secret for %s: %w", engine, err)
	}

	return secret, nil
}

// createPath takes a path, and adds a new path to a KV v2 engine
func (vc *vaultClient) createPath(path string) error {
	vc.tracer.trace(fmt.Sprintf("%s/createPath", vc.config.tracePrefix))

	_, err := vc.write(path, nil)
	if err != nil {
		return fmt.Errorf("failed to create new path at %s: %w", path, err)
	}

	return nil
}
