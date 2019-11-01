package vault

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/hashicorp/vault/api"
	log "github.com/sirupsen/logrus"
	"go.opencensus.io/trace"
)

type vaultClient struct {
	client         *api.Client
	traceEnabled   bool
	tracePrefix    string
	project        string
	serviceAccount string
	vaultRole      string

	ctx context.Context
}

func newVaultClient() *vaultClient {
	c := &vaultClient{}

	return c
}

func (c *vaultClient) loadVaultEnvironment() error {

	traceEnabledString := getEnv("TRACE_ENABLED", "false")
	c.traceEnabled, _ = strconv.ParseBool(traceEnabledString)

	if c.traceEnabled {
		_, span := trace.StartSpan(c.ctx, c.tracePrefix+"/getConfigFromEnv")
		defer span.End()
	}

	c.tracePrefix = getEnv("TRACE_PREFIX", "vault")
	if c.tracePrefix == "" {
		return errors.New("Error occurred getting TRACE_PREFIX variable from environment")
	}

	c.vaultRole = getEnv("VAULT_ROLE", "")
	if c.vaultRole == "" {
		return errors.New("You need to set the VAULT_ROLE environment variable")
	}

	// google injects this env var automatically in gcp environments
	c.project = getEnv("GCLOUD_PROJECT", "")
	if c.project == "" {
		return errors.New("You need to set the GCLOUD_PROJECT environment variable")
	}

	// google injects this env var automatically in gcp environments
	c.serviceAccount = getEnv("FUNCTION_IDENTITY", "")
	if c.serviceAccount == "" {
		return errors.New("You need to set the FUNCTION_IDENTITY environment variable")
	}

	log.Info(fmt.Sprintf("TRACE_PREFIX=%s, VAULT_ROLE=%s, GCLOUD_PROJECT=%s, FUNCTION_IDENTITY=%s", c.tracePrefix, c.vaultRole, c.project, c.serviceAccount))

	return nil
}

// initClient takes context and a vault role and returns an initialized Vault
// client using the value in the "VAULT_ADDR" env var.
// It will exit the process if it fails to initialize.
func (c *vaultClient) initClient() (err error) {
	if c.traceEnabled {
		_, span := trace.StartSpan(c.ctx, c.tracePrefix+"/InitVaultClient")
		defer span.End()
	}

	vaultAddr, err := getEncrEnvVar(c.ctx, "VAULT_ADDR")
	if err != nil {
		return fmt.Errorf("Error getting vault address: %v", err)
	}

	c.client, err = api.NewClient(&api.Config{
		Address: vaultAddr,
	})
	if err != nil {
		return fmt.Errorf("Error initializing new vault api client: %v", err)
	}
	token, err := getVaultToken(c)
	if err != nil {
		return err
	}

	c.client.SetToken(token)
	if err != nil {
		return fmt.Errorf("Error setting vault token: %v", err)
	}

	return err
}

// getSecretFromVault takes a vault client, key name, and data name, and returns the
// value returned from vault as a string.
func (c *vaultClient) getSecretFromVault(secretName string) (map[string]string, error) {
	if c.traceEnabled {
		_, span := trace.StartSpan(c.ctx, c.tracePrefix+"/GetSecret")
		defer span.End()
	}

	secretMap := map[string]string{}

	secretValues, err := c.client.Logical().Read(secretName)
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
