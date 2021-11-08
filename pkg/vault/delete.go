package vault

import (
	"fmt"

	"github.com/hashicorp/vault/api"
	"go.opencensus.io/trace"
)

// Delete takes a given key for an engine, and removes the key/value pair from vault.
func (vc *vaultClient) Delete(engine, key string) (*api.Secret, error) {
	if vc.config.traceEnabled {
		var span *trace.Span
		vc.ctx, span = trace.StartSpan(vc.ctx, fmt.Sprintf("%s/Delete", vc.config.tracePrefix))
		defer span.End()
	}

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
