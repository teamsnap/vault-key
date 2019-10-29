// Package key provides access to the vault.
package key

import (
	"C"
	"context"
	"encoding/json"
	"fmt"

	"github.com/teamsnap/vault-key/pkg/vault"
)

var env = map[string]map[string]string{}

// Loot grabs secrets from the vault.
func Loot(secretNames string) (string, error) {
	ctx := context.Background()

	var envArr []string
	err := json.Unmarshal([]byte(secretNames), &envArr)
	if err != nil {
		return "", fmt.Errorf("Failed to unmarshall secrets as json.  Error: %v", err)
	}

	vault.GetSecrets(ctx, &env, envArr)

	secrets, err := json.Marshal(env)
	if err != nil {
		return "", fmt.Errorf("Failed to marshall secrets as json.  Error: %v", err)
	}

	return string(secrets), nil
}
