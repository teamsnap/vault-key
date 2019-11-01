package vault

import (
	"C"
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/hashicorp/vault/api"
	log "github.com/sirupsen/logrus"
	"go.opencensus.io/trace"
)

// Loot returns a map encoded in json with the values of secrets pulled from Vault.
func Loot(secretNames string) (string, error) {
	var envArr []string
	env := map[string]map[string]string{}

	err := json.Unmarshal([]byte(secretNames), &envArr)
	if err != nil {
		return "", fmt.Errorf("Failed to unmarshall secrets as json.  Error: %v", err)
	}

	err = GetSecrets(context.Background(), &env, envArr)
	if err != nil {
		return "", err
	}

	secrets, err := json.Marshal(env)
	if err != nil {
		return "", fmt.Errorf("Failed to marshall secrets as json.  Error: %v", err)
	}

	return string(secrets), nil
}

// GetSecrets fills a map with the values of secrets pulled from Vault.
func GetSecrets(ctx context.Context, secretValues *map[string]map[string]string, secretNames []string) error {
	var err error
	a, err := initVault(ctx)

	if a.traceEnabled {
		var span *trace.Span
		a.ctx, span = trace.StartSpan(a.ctx, fmt.Sprintf("%s/GetSecrets", a.tracePrefix))
		defer span.End()
	}

	a.client, err = a.initVaultClient()
	if err != nil {
		return err
	}

	for _, secretName := range secretNames {
		log.Debug("secret, err := GetSecret(c, " + secretName + ")")

		secret, err := a.getSecret(ctx, a.client, secretName)
		if err != nil {
			return fmt.Errorf("Error getting secret: %v", err)
		}

		(*secretValues)[secretName] = secret
	}

	return nil
}

func initVault(ctx context.Context) (*App, error) {
	a := newApp()
	a.ctx = ctx
	a.environment = os.Getenv("ENVIRONMENT")

	if a.environment == "production" {
		log.SetFormatter(&log.JSONFormatter{})
		log.SetLevel(log.WarnLevel)
	} else {
		log.SetFormatter(&log.TextFormatter{})
		log.SetLevel(log.TraceLevel)
	}

	err := a.getConfigFromEnv()
	if err != nil {
		return nil, fmt.Errorf("Unable to laod config from environment.  Error: %v ", err)
	}

	return a, nil
}

// initVaultClient takes context and a vault role and returns an initialized Vault
// client using the value in the "VAULT_ADDR" env var.
// It will exit the process if it fails to initialize.
func (a *App) initVaultClient() (*api.Client, error) {
	if a.traceEnabled {
		var span *trace.Span
		a.ctx, span = trace.StartSpan(a.ctx, fmt.Sprintf("%s/InitVaultClient", a.tracePrefix))
		defer span.End()
	}

	vaultAddr, err := a.getEncrEnvVar(a.ctx, "VAULT_ADDR")
	if err != nil {
		return nil, fmt.Errorf("Error getting vault address: %v", err)
	}

	vaultToken, err := a.getVaultToken()
	if err != nil {
		return nil, err
	}

	a.client, err = a.newVaultClient(vaultAddr, vaultToken)
	if err != nil {
		return nil, err
	}

	return a.client, nil
}

// newVaultClient returns a configured vault api client
func (a *App) newVaultClient(addr, token string) (*api.Client, error) {
	var err error
	if a.traceEnabled {
		var span *trace.Span
		a.ctx, span = trace.StartSpan(a.ctx, fmt.Sprintf("%s/NewVaultClient", a.tracePrefix))
		defer span.End()
	}

	a.client, err = api.NewClient(&api.Config{
		Address: addr,
	})
	if err != nil {
		return nil, fmt.Errorf("Error initializing vault client: %v", err)
	}

	a.client.SetToken(token)

	return a.client, nil
}

// getSecret takes a vault client, key name, and data name, and returns the
// value returned from vault as a string.
func (a *App) getSecret(ctx context.Context, c *api.Client, secretName string) (map[string]string, error) {
	if a.traceEnabled {
		var span *trace.Span
		a.ctx, span = trace.StartSpan(a.ctx, fmt.Sprintf("%s/getSecret", a.tracePrefix))
		defer span.End()
	}

	secretMap := map[string]string{}

	secretValues, err := c.Logical().Read(secretName)
	if err != nil {
		return nil, fmt.Errorf("error reading secret from Vault for %s: Error %v", secretName, err)
	}

	if secretValues == nil {
		return nil, fmt.Errorf("secret values returned from Vault are <nil> for %s", secretName)
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
	return nil, fmt.Errorf("failed to convert secret data from Vault to a string for %s", secretName)
}
