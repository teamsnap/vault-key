package vault

import (
	"fmt"

	"github.com/hashicorp/vault/api"
)

func (vc *vaultClient) write(engine string, m map[string]string) (*api.Secret, error) {
	vc.tracer.trace(fmt.Sprintf("%s/write", vc.config.tracePrefix))

	data := make(map[string]interface{}, len(m))
	for k, v := range m {
		data[k] = v
	}

	secrets := map[string]interface{}{
		"data": data,
	}

	secret, err := vc.client.Logical().Write(engine, secrets)
	if err != nil {
		return secret, fmt.Errorf("failed to write data to %s: %w", engine, err)
	}

	return secret, nil
}
