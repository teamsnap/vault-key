package vault

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/GoogleCloudPlatform/berglas/pkg/berglas"
	log "github.com/sirupsen/logrus"
	"go.opencensus.io/trace"
)

// NewVaultToken uses a github token or service account to get a vault auth token
func NewVaultToken(vc *vaultClient) (string, error) {
	if vc.config.traceEnabled {
		var span *trace.Span
		vc.ctx, span = trace.StartSpan(vc.ctx, fmt.Sprintf("%s/NewVaultToken", vc.config.tracePrefix))
		defer span.End()
	}

	return NewAuthClient(vc.config).GetVaultToken(vc)
}

// GetSecrets fills a map with the values of secrets pulled from Vault.
func GetSecrets(ctx context.Context, secretValues *map[string]map[string]string, secretNames []string) error {
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
		return fmt.Errorf("load client environment: %v", err)
	}

	vc, err := NewVaultClient(ctx, config)

	if err != nil {
		return fmt.Errorf("error initializing vault client: %v", err)
	}

	if config.traceEnabled {
		var span *trace.Span
		ctx, span = trace.StartSpan(ctx, fmt.Sprintf("%s/GetSecrets", config.tracePrefix))
		defer span.End()
	}

	for _, secretName := range secretNames {
		log.Debug(fmt.Sprintf("secret= %s", secretNames))

		secret, err := vc.SecretFromVault(secretName)
		if err != nil {
			return fmt.Errorf("getting secret: %v", err)
		}

		(*secretValues)[secretName] = secret
	}

	return nil
}

// GetSecretVersions fills a map with the versions of secrets pulled from Vault.
func GetSecretVersions(ctx context.Context, secretVersions *map[string]int64, secretNames []string) error {
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
		return fmt.Errorf("load client environment: %v", err)
	}

	vc, err := NewVaultClient(ctx, config)

	if config.traceEnabled {
		var span *trace.Span
		ctx, span = trace.StartSpan(ctx, fmt.Sprintf("%s/GetSecrets", config.tracePrefix))
		defer span.End()
	}

	for _, secretName := range secretNames {
		log.Debug(fmt.Sprintf("secret= %s", secretNames))

		secretVersion, err := vc.SecretVersionFromVault(secretName)
		if err != nil {
			return fmt.Errorf("getting secret version: %v", err)
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
