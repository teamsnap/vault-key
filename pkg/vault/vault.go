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

func newVault() *vault {
	v := &vault{}
	return v
}

type vault struct {
	traceEnabled   bool
	tracePrefix    string
	project        string
	serviceAccount string
	vaultRole      string
	environment    string
	env            map[string]map[string]string
	client         *api.Client
	ctx            context.Context
}

// Loot returns a map encoded in json with the values of secrets pulled from Vault.
func Loot(secretNames string) (string, error) {
	v, err := initVault(context.Background())

	var envArr []string

	err = json.Unmarshal([]byte(secretNames), &envArr)
	if err != nil {
		return "", fmt.Errorf("Failed to unmarshall secrets as json.  Error: %v", err)
	}

	v.getSecrets(v.ctx, &v.env, envArr)

	secrets, err := json.Marshal(v.env)
	if err != nil {
		return "", fmt.Errorf("Failed to marshall secrets as json.  Error: %v", err)
	}

	return string(secrets), nil

}

// GetSecrets fills a map with the values of secrets pulled from Vault.
func (v *vault) GetSecrets(ctx context.Context, secretValues *map[string]map[string]string, secretNames []string) error {
	var err error

	if v.traceEnabled {
		_, span := trace.StartSpan(v.ctx, fmt.Sprintf("%s/GetSecrets", v.tracePrefix))
		defer span.End()
	}

	v.client, err = v.initVaultClient()
	if err != nil {
		return err
	}

	for _, secretName := range secretNames {
		log.Debug("secret, err := GetSecret(c, " + secretName + ")")

		secret, err := v.getSecret(ctx, v.client, secretName)
		if err != nil {
			return fmt.Errorf("Error getting secret: %v", err)
		}

		(*secretValues)[secretName] = secret
	}

	return nil
}

func initVault(ctx context.Context) (*vault, error) {
	v := newVault()
	v.environment = os.Getenv("ENVIRONMENT")

	if v.environment == "production" {
		log.SetFormatter(&log.JSONFormatter{})
		log.SetLevel(log.WarnLevel)
	} else {
		log.SetFormatter(&log.TextFormatter{})
		log.SetLevel(log.TraceLevel)
	}

	err := v.getConfigFromEnv(ctx)
	if err != nil {
		return nil, fmt.Errorf("Unable to laod config from environment.  Error: %v ", err)
	}

	return v, nil
}

// initVaultClient takes context and a vault role and returns an initialized Vault
// client using the value in the "VAULT_ADDR" env var.
// It will exit the process if it fails to initialize.
func (v *vault) initVaultClient() (*api.Client, error) {
	if v.traceEnabled {
		_, span := trace.StartSpan(v.ctx, fmt.Sprintf("%s/InitVaultClient", v.tracePrefix))
		defer span.End()
	}

	vaultAddr, err := v.getEncrEnvVar(v.ctx, "VAULT_ADDR")
	if err != nil {
		return nil, fmt.Errorf("Error getting vault address: %v", err)
	}

	vaultToken, err := v.getVaultToken()
	if err != nil {
		return nil, err
	}

	v.client, err = v.newVaultClient(vaultAddr, vaultToken)
	if err != nil {
		return nil, err
	}

	return v.client, nil
}

// newVaultClient returns a configured vault api client
func (v *vault) newVaultClient(addr, token string) (*api.Client, error) {
	var err error
	if v.traceEnabled {
		_, span := trace.StartSpan(v.ctx, fmt.Sprintf("%s/NewVaultClient", v.tracePrefix))
		defer span.End()
	}

	v.client, err = api.NewClient(&api.Config{
		Address: addr,
	})
	if err != nil {
		return nil, fmt.Errorf("Error initializing vault client: %v", err)
	}

	v.client.SetToken(token)

	return v.client, nil
}

// getSecret takes a vault client, key name, and data name, and returns the
// value returned from vault as a string.
func (v *vault) getSecret(ctx context.Context, c *api.Client, secretName string) (map[string]string, error) {
	if v.traceEnabled {
		_, span := trace.StartSpan(v.ctx, fmt.Sprintf("%s/getSecret", v.tracePrefix))
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
