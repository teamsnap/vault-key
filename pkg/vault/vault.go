package vault

import (
	"context"
	"errors"
	"github.com/hashicorp/vault/api"
	log "github.com/sirupsen/logrus"
	"go.opencensus.io/trace"
)

// InitVaultClient takes context and a vault role and returns an initialized Vault
// client using the value in the "VAULT_ADDR" env var.
// It will exit the process if it fails to initialize.
func InitVaultClient(ctx context.Context, vaultRole string) *api.Client {
	if traceEnabled {
		_, span := trace.StartSpan(ctx, tracePrefix+"/InitVaultClient")
		defer span.End()
	}

	vaultAddr, errVaultAddr := getEncrEnvVar(ctx, "VAULT_ADDR")

	if errVaultAddr != nil {
		log.Error("Error getting vault address:", errVaultAddr)
		return nil
	}

	vaultToken := getVaultToken(ctx, project, serviceAccount, vaultRole)

	vaultClient, err := NewVaultClient(ctx, vaultAddr, vaultToken)

	if err != nil {
		log.Error("Error initializing vault client:", err)
		return nil
	}

	return vaultClient
}

// NewVaultClient returns a configured vault api client
func NewVaultClient(ctx context.Context, addr, token string) (*api.Client, error) {
	if traceEnabled {
		_, span := trace.StartSpan(ctx, tracePrefix+"/NewVaultClient")
		defer span.End()
	}

	client, err := api.NewClient(&api.Config{
		Address: addr,
	})

	if err != nil {
		return nil, err
	}

	client.SetToken(token)

	return client, nil
}

// GetSecret takes a vault client, key name, and data name, and returns the
// value returned from vault as a string.
func GetSecret(ctx context.Context, c *api.Client, secretName string) (map[string]string, error) {
	if traceEnabled {
		_, span := trace.StartSpan(ctx, tracePrefix+"/GetSecret")
		defer span.End()
	}

	secretMap := map[string]string{}

	secretValues, err := c.Logical().Read(secretName)
	if err != nil {
		log.Error("Error reading secret from Vault:", err)
		return secretMap, errors.New("error reading secret from Vault for " + secretName + ": " + err.Error())
	}

	if secretValues == nil {
		log.Error("Secret values returned from Vault are <nil> for " + secretName)
		return secretMap, errors.New("secret values returned from Vault are <nil> for " + secretName)
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
	return secretMap, errors.New("failed to convert secret data from Vault to a string for " + secretName)
}

// GetSecrets fills a map with the values of secrets pulled from Vault.
func GetSecrets(ctx context.Context, secretValues *map[string]map[string]string, secretNames []string) {
	if traceEnabled {
		_, span := trace.StartSpan(ctx, tracePrefix+"/GetSecrets")
		defer span.End()
	}

	initialize(ctx)

	c := InitVaultClient(ctx, vaultRole)

	for _, secretName := range secretNames {
		log.Debug("secret, err := GetSecret(c, " + secretName + ")")

		secret, err := GetSecret(ctx, c, secretName)
		if err != nil {
			log.Error("Error getting secret:", err)
		} else {
			(*secretValues)[secretName] = secret
		}
	}
}
