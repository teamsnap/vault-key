package vault

import (
	"context"
	"fmt"

	"encoding/json"

	"github.com/hashicorp/vault/api"
	log "github.com/sirupsen/logrus"
	"go.opencensus.io/trace"
)

type vaultClient struct {
	authClient AuthClient
	client     *api.Client
	config     *config
	ctx        context.Context
}

func newVaultClient(c *config) *vaultClient {
	client := &vaultClient{
		config: c,
	}

	return client
}

// initClient takes context and a vault role and returns an initialized Vault
// client using the value in the "VAULT_ADDR" env var.
// It will exit the process if it fails to initialize.
func (vc *vaultClient) initClient() (err error) {
	if vc.config.traceEnabled {
		var span *trace.Span
		vc.ctx, span = trace.StartSpan(vc.ctx, fmt.Sprintf("%s/initClient", vc.config.tracePrefix))
		defer span.End()
	}

	vaultAddr, err := getEncrEnvVar(vc.ctx, "VAULT_ADDR")
	if err != nil {
		return fmt.Errorf("vault address: %v", err)
	}

	vc.client, err = api.NewClient(&api.Config{
		Address: vaultAddr,
	})
	if err != nil {
		return fmt.Errorf("initializing new vault api client: %v", err)
	}

	return err
}

// getSecretFromVault takes a vault client, key name, and data name, and returns the
// value returned from vault as a string.
func (vc *vaultClient) getSecretFromVault(secretName string) (map[string]string, error) {
	if vc.config.traceEnabled {
		var span *trace.Span
		vc.ctx, span = trace.StartSpan(vc.ctx, fmt.Sprintf("%s/getSecretFromVault", vc.config.tracePrefix))
		defer span.End()
	}

	secretMap := map[string]string{}

	secretValues, err := vc.client.Logical().Read(secretName)
	if err != nil {
		log.Error(fmt.Sprintf("Error reading secret from Vault: %v", err))
		return secretMap, fmt.Errorf("error reading secret from Vault for %s", secretName)
	}

	if secretValues == nil {
		log.Error("Secret values returned from Vault are <nil> for " + secretName)
		return secretMap, fmt.Errorf("secret values returned from Vault are <nil> for %s", secretName)
	}

	// https://stackoverflow.com/questions/26975880/convert-mapinterface-interface-to-mapstringstring
	m, ok := secretValues.Data["data"].(map[string]interface{})
	if ok {
		for key, value := range m {
			secretMap[key] = value.(string)
		}

		return secretMap, nil
	}

	log.Errorf("%T %#v\n", secretValues.Data["data"], secretValues.Data["data"])
	return secretMap, fmt.Errorf("failed to convert secret data from Vault to a string for %s", secretName)
}

// getSecretVersionFromVault takes a vault client, key name, and data name, and returns the
// version of the Vault secret as an int.
func (vc *vaultClient) getSecretVersionFromVault(secretName string) (int64, error) {
	if vc.config.traceEnabled {
		var span *trace.Span
		vc.ctx, span = trace.StartSpan(vc.ctx, fmt.Sprintf("%s/getSecretFromVault", vc.config.tracePrefix))
		defer span.End()
	}

	version := int64(0)

	secretValues, err := vc.client.Logical().Read(secretName)
	if err != nil {
		log.Error(fmt.Sprintf("Error reading secret from Vault: %v", err))
		return version, fmt.Errorf("error reading secret from Vault for %s", secretName)
	}

	version, err = secretValues.Data["metadata"].(map[string]interface{})["version"].(json.Number).Int64()
	if err != nil {
		log.Error(fmt.Sprintf("Error converting secret version to integer for %s: %v", secretName, err))
		return version, fmt.Errorf("Error converting secret version to integer for %s: %v", secretName, err)
	}

	return version, nil
}
