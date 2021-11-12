package vault

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/GoogleCloudPlatform/berglas/pkg/berglas"
	log "github.com/sirupsen/logrus"
)

// NewVaultToken uses a github token or service account to get a vault auth token
func NewVaultToken(vc *vaultClient) (string, error) {
	vc.tracer.trace(fmt.Sprintf("%s/NewVaultToken", vc.config.tracePrefix))

	return NewAuthClient(vc.config).GetVaultToken(vc)
}

// GetSecrets fills a map with the values of secrets pulled from Vault.
func GetSecrets(ctx context.Context, secretValues *map[string]map[string]string, secretNames []string) error {
	config, err := getConfig()
	if err != nil {
		return err
	}

	vc, err := NewVaultClient(ctx, config)
	if err != nil {
		return fmt.Errorf("error initializing vault client: %w", err)
	}

	vc.tracer.trace(fmt.Sprintf("%s/GetSecrets", vc.config.tracePrefix))

	for _, secretName := range secretNames {
		secret, err := vc.SecretFromVault(secretName)
		if err != nil {
			return fmt.Errorf("getting secret: %w", err)
		}

		(*secretValues)[secretName] = secret
	}

	return nil
}

// CreateSecret takes a given key for an engine, and adds a new key/value pair in vault.
func CreateSecret(ctx context.Context, engine, key, value string) error {
	config, err := getConfig()
	if err != nil {
		return err
	}

	vc, err := NewVaultClient(ctx, config)
	if err != nil {
		return fmt.Errorf("error initializing vault client: %w", err)
	}

	vc.tracer.trace(fmt.Sprintf("%s/CreateSecret", vc.config.tracePrefix))

	if _, err := vc.create(engine, key, value); err != nil {
		return err
	}

	return nil
}

// UpdateSecret takes a given key for an engine, and modifies its value in vault.
func UpdateSecret(ctx context.Context, engine, key, value string) error {
	config, err := getConfig()
	if err != nil {
		return err
	}

	vc, err := NewVaultClient(ctx, config)
	if err != nil {
		return fmt.Errorf("error initializing vault client: %w", err)
	}

	vc.tracer.trace(fmt.Sprintf("%s/UpdateSecret", vc.config.tracePrefix))

	if _, err := vc.update(engine, key, value); err != nil {
		return err
	}

	return nil
}

func DeleteSecret(ctx context.Context, engine, key string) error {
	config, err := getConfig()
	if err != nil {
		return err
	}

	vc, err := NewVaultClient(ctx, config)
	if err != nil {
		return fmt.Errorf("error initializing vault client: %w", err)
	}

	vc.tracer.trace(fmt.Sprintf("%s/DeleteSecret", vc.config.tracePrefix))

	if _, err := vc.delete(engine, key); err != nil {
		return err
	}

	return nil
}

// GetSecretVersions fills a map with the versions of secrets pulled from Vault.
func GetSecretVersions(ctx context.Context, secretVersions *map[string]int64, secretNames []string) error {
	config, err := getConfig()
	if err != nil {
		return err
	}

	vc, err := NewVaultClient(ctx, config)
	if err != nil {
		return fmt.Errorf("error initializing vault client: %w", err)
	}

	vc.tracer.trace(fmt.Sprintf("%s/GetSecretVersions", vc.config.tracePrefix))

	for _, secretName := range secretNames {
		secretVersion, err := vc.SecretVersionFromVault(secretName)
		if err != nil {
			return fmt.Errorf("getting secret version: %w", err)
		}

		(*secretVersions)[secretName] = secretVersion
	}

	return nil
}

// getEncrEnvVar takes the name of an environment variable that's value begins
// with "berglas://", decrypts the value from a Google Storage Bucket with KMS,
// replaces the original environment variable value with the decrypted value,
// and returns the value as a string. If there's an error fetching the value, it
// will return an empty string along with the error message.
func getEncrEnvVar(ctx context.Context, n string) (string, error) {
	val := os.Getenv(n)
	if strings.HasPrefix(val, "berglas://") {
		if err := berglas.Replace(ctx, n); err != nil {
			return "", err
		}
	}

	return os.Getenv(n), nil
}

func getConfig() (*config, error) {
	environment := os.Getenv("ENVIRONMENT")

	if environment == "production" {
		log.SetFormatter(&log.JSONFormatter{})
		log.SetLevel(log.WarnLevel)
	} else {
		log.SetFormatter(&log.TextFormatter{})
		log.SetLevel(log.TraceLevel)
	}

	config, err := loadVaultEnvironment()
	if err != nil {
		return nil, fmt.Errorf("load client environment: %w", err)
	}

	return config, nil

}
