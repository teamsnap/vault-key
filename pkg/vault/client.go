package vault

import (
	"context"
	"fmt"

	"encoding/json"

	"github.com/hashicorp/vault/api"
	log "github.com/sirupsen/logrus"
	"go.opencensus.io/trace"
)

// Client is a vault api client that is authorized to get secrets out of a vault.
type Client interface {
	GetSecretFromVault(secret string) (map[string]string, error)
	GetSecretVersionFromVault(secret string) (int64, error)
}

type vaultClient struct {
	client *api.Client
	config *config
	ctx    context.Context
}

// NewVaultClient configures and returns an initialized vault client.
func NewVaultClient(ctx context.Context, c *config) (Client, error) {
	client := &vaultClient{
		config: c,
		ctx:    ctx,
	}

	err := initClient(client)
	if err != nil {
		return nil, fmt.Errorf("initialze client: %v", err)
	}

	return client, nil
}

// initClient takes context and a vault role and returns an initialized Vault
// client using the value in the "VAULT_ADDR" env var.
// It will exit the process if it fails to initialize.
func initClient(vc *vaultClient) error {
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

	token, err := NewVaultToken(vc)
	if err != nil {
		return fmt.Errorf("getting vault api token from client: %v", err)
	}

	vc.client.SetToken(token)
	return err
}

// GetSecretFromVault takes a vault client, key name, and data name, and returns the
// value returned from vault as a string.
func (vc *vaultClient) GetSecretFromVault(secretName string) (map[string]string, error) {
	if vc.config.traceEnabled {
		var span *trace.Span
		vc.ctx, span = trace.StartSpan(vc.ctx, fmt.Sprintf("%s/getSecretFromVault", vc.config.tracePrefix))
		defer span.End()
	}

	secretMap := map[string]string{}

	secretValues, err := vc.client.Logical().Read(secretName)
	if err != nil {
		log.Error(fmt.Sprintf("reading secret from Vault: %v", err))
		return secretMap, fmt.Errorf("reading secret from Vault for %s", secretName)
	}

	if secretValues == nil {
		log.Error("secret values returned from Vault are <nil> for " + secretName)
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
	return secretMap, fmt.Errorf("converting secret data from Vault to a string for %s", secretName)
}

// GetSecretVersionFromVault takes a vault client, key name, and data name, and returns the
// version of the Vault secret as an int.
func (vc *vaultClient) GetSecretVersionFromVault(secretName string) (int64, error) {
	if vc.config.traceEnabled {
		var span *trace.Span
		vc.ctx, span = trace.StartSpan(vc.ctx, fmt.Sprintf("%s/getSecretFromVault", vc.config.tracePrefix))
		defer span.End()
	}

	version := int64(0)

	secretValues, err := vc.client.Logical().Read(secretName)
	if err != nil {
		log.Error(fmt.Sprintf("Error reading secret from Vault: %v", err))
		return version, fmt.Errorf("reading secret from Vault for %s", secretName)
	}

	version, err = secretValues.Data["metadata"].(map[string]interface{})["version"].(json.Number).Int64()
	if err != nil {
		log.Error(fmt.Sprintf("Error converting secret version to integer for %s: %v", secretName, err))
		return version, fmt.Errorf("converting secret version to integer for %s: %v", secretName, err)
	}

	return version, nil
}
