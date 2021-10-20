package vault

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/vault/api"
	"go.opencensus.io/trace"
)

type vaultClient struct {
	client *api.Client
	config *config
	ctx    context.Context
	tracer
}

// NewVaultClient configures and returns an initialized vault client.
func NewVaultClient(ctx context.Context, c *config) (*vaultClient, error) {
	client := &vaultClient{
		config: c,
		ctx:    ctx,
	}
	client.tracer = client

	err := initClient(client)
	if err != nil {
		return nil, fmt.Errorf("initialze client: %w", err)
	}

	return client, nil
}

// initClient takes context and a vault role and returns an initialized Vault
// client using the value in the "VAULT_ADDR" env var.
// It will exit the process if it fails to initialize.
func initClient(vc *vaultClient) error {
	vc.tracer.trace(fmt.Sprintf("%s/initClient", vc.config.tracePrefix))

	vaultAddr, err := getEncrEnvVar(vc.ctx, "VAULT_ADDR")
	if err != nil {
		return fmt.Errorf("vault address: %w", err)
	}

	vc.client, err = api.NewClient(&api.Config{
		Address: vaultAddr,
	})
	if err != nil {
		return fmt.Errorf("initializing new vault api client: %w", err)
	}

	token, err := NewVaultToken(vc)
	if err != nil {
		return fmt.Errorf("getting vault api token from client: %w", err)
	}

	vc.client.SetToken(token)
	return err
}

// SecretFromVault takes a secret name and returns the value returned from vault as a string.
func (vc *vaultClient) SecretFromVault(secretName string) (map[string]string, error) {
	vc.tracer.trace(fmt.Sprintf("%s/SecretFromVault", vc.config.tracePrefix))

	secretMap := map[string]string{}

	secretValues, err := vc.client.Logical().Read(secretName)
	if err != nil {
		return secretMap, fmt.Errorf("reading secret from Vault for %s", secretName)
	}

	if secretValues == nil {
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

	return secretMap, fmt.Errorf("converting secret data from Vault to a string for %s", secretName)
}

// SecretVersionFromVault takes a secret name and returns the version of the Vault secret as an int.
func (vc *vaultClient) SecretVersionFromVault(secretName string) (int64, error) {
	vc.tracer.trace(fmt.Sprintf("%s/SecretVersionFromVault", vc.config.tracePrefix))

	version := int64(0)

	secretValues, err := vc.client.Logical().Read(secretName)
	if err != nil {
		return version, fmt.Errorf("reading secret from Vault for %s", secretName)
	}

	version, err = secretValues.Data["current_version"].(json.Number).Int64()
	if err != nil {
		return version, fmt.Errorf("converting secret version to integer for %s: %v", secretName, err)
	}

	return version, nil
}

type tracer interface {
	trace(string) func()
}

func (vc *vaultClient) trace(name string) func() {
	if !vc.config.traceEnabled {
		return func() {}
	}

	var span *trace.Span
	vc.ctx, span = trace.StartSpan(
		vc.ctx,
		name,
	)

	return func() { defer span.End() }

// EnginesFromVault takes a path and returns a list of engines from vault.
func (vc *vaultClient) EnginesFromVault(path string) ([]string, error) {
	if vc.config.traceEnabled {
		var span *trace.Span
		vc.ctx, span = trace.StartSpan(vc.ctx, fmt.Sprintf("%s/EnginesFromVault", vc.config.tracePrefix))
		defer span.End()
	}

	engines, err := vc.client.Logical().List(path)
	if err != nil {
		log.Error(fmt.Sprintf("listing engines from Vault: %v", err))
		return nil, fmt.Errorf("listing engines from Vault for %s", path)
	}

	if engines == nil {
		log.Error("engines returned from Vault are <nil> for " + path)
		return nil, fmt.Errorf("engines returned from Vault are <nil> for %s", path)
	}

	engineData, _ := extractListData(engines)

	result := []string{}

	for _, value := range engineData {
		switch v := value.(type) {
		case string:
			result = append(result, v)
		default:
			return nil, fmt.Errorf("unexpected type, expected string, got: %T, value: %v", v, result)
		}
	}
	return result, nil
}

func extractListData(secret *api.Secret) ([]interface{}, bool) {
	if secret == nil || secret.Data == nil {
		return nil, false
	}

	k, ok := secret.Data["keys"]
	if !ok || k == nil {
		return nil, false
	}

	i, ok := k.([]interface{})
	return i, ok
}
