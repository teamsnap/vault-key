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

func (vc *vaultClient) loadVaultEnvironment() error {

	traceEnabledString := getEnv("TRACE_ENABLED", "false")
	vc.traceEnabled, _ = strconv.ParseBool(traceEnabledString)

	if vc.traceEnabled {
		var span *trace.Span
		vc.ctx, span = trace.StartSpan(vc.ctx, fmt.Sprintf("%s/loadVaultEnvironment", vc.tracePrefix))
		defer span.End()
	}

	vc.tracePrefix = getEnv("TRACE_PREFIX", "vault")
	if vc.tracePrefix == "" {
		return errors.New("Error occurred getting TRACE_PREFIX variable from environment")
	}

	vc.vaultRole = getEnv("VAULT_ROLE", "")
	if vc.vaultRole == "" {
		return errors.New("You need to set the VAULT_ROLE environment variable")
	}

	// google injects this env var automatically in gcp environments
	vc.project = getEnv("GCLOUD_PROJECT", "")
	if vc.project == "" {
		return errors.New("You need to set the GCLOUD_PROJECT environment variable")
	}

	// google injects this env var automatically in gcp environments
	vc.serviceAccount = getEnv("FUNCTION_IDENTITY", "")
	if vc.serviceAccount == "" {
		return errors.New("You need to set the FUNCTION_IDENTITY environment variable")
	}

	log.Info(fmt.Sprintf("TRACE_PREFIX=%s, VAULT_ROLE=%s, GCLOUD_PROJECT=%s, FUNCTION_IDENTITY=%s", vc.tracePrefix, vc.vaultRole, vc.project, vc.serviceAccount))

	return nil
}

// initClient takes context and a vault role and returns an initialized Vault
// client using the value in the "VAULT_ADDR" env var.
// It will exit the process if it fails to initialize.
func (vc *vaultClient) initClient() (err error) {
	if vc.traceEnabled {
		var span *trace.Span
		vc.ctx, span = trace.StartSpan(vc.ctx, fmt.Sprintf("%s/initClient", vc.tracePrefix))
		defer span.End()
	}

	vaultAddr, err := getEncrEnvVar(vc.ctx, "VAULT_ADDR")
	if err != nil {
		return fmt.Errorf("Error getting vault address: %v", err)
	}

	vc.client, err = api.NewClient(&api.Config{
		Address: vaultAddr,
	})
	if err != nil {
		return fmt.Errorf("Error initializing new vault api client: %v", err)
	}
	token, err := getVaultToken(vc)
	if err != nil {
		return err
	}

	vc.client.SetToken(token)
	if err != nil {
		return fmt.Errorf("Error setting vault token: %v", err)
	}

	return err
}

// getSecretFromVault takes a vault client, key name, and data name, and returns the
// value returned from vault as a string.
func (vc *vaultClient) getSecretFromVault(secretName string) (map[string]string, error) {
	if vc.traceEnabled {
		var span *trace.Span
		vc.ctx, span = trace.StartSpan(vc.ctx, fmt.Sprintf("%s/getSecretFromVault", vc.tracePrefix))
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
