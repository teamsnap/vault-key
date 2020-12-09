package vault

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/GoogleCloudPlatform/berglas/pkg/berglas"
	log "github.com/sirupsen/logrus"
	"go.opencensus.io/trace"
	"google.golang.org/api/iam/v1"
)

// AuthClient is a type that satifies the necesary authorization layer for a vault client.
type AuthClient interface {
	GetVaultToken(vc *vaultClient) (string, error)
}

// Client is a vault api client that is authorized to get secrets out of a vault.
type Client interface {
	GetSecretFromVault(secret string) (map[string]string, error)
	GetSecretVersionFromVault(secret string) (int64, error)
}

// GetVaultToken uses a service account to get a vault auth token
func (a *gcpAuthClient) GetVaultToken(vc *vaultClient) (string, error) {
	if vc.config.traceEnabled {
		var span *trace.Span
		vc.ctx, span = trace.StartSpan(vc.ctx, fmt.Sprintf("%s/getVaultToken", vc.config.tracePrefix))
		defer span.End()
	}

	var err error

	if vc.config.githubAuth {
		vaultResp, err := githubVaultAuth(vc)
		if err != nil {
			return "", err
		}

		log.Debug("Successfully logged into Vault with auth/github/login")

		return vaultResp.Auth.ClientToken, nil
	}

	if vc.config.googleAuth {
		a.iamService, err = iam.NewService(vc.ctx)
		if err != nil {
			log.Debugf("initialze client: getting new iam service: %v", err)
			return "", fmt.Errorf("getting new iam service: google: could not find default credentials")
		}

		log.Debug("Successfully created IAM client")

		err = a.generateSignedJWT(vc.ctx, vc.config)
		if err != nil {
			log.Debugf("generating signed jwt: %v", err)
			return "", fmt.Errorf("generating signed jwt, sigining jwt: Post")
		}

		log.Debug("Successfully generated signed JWT")

		vaultResp, err := a.gcpSaAuth(vc)
		if err != nil {
			return "", err
		}

		log.Debug("Successfully logged into Vault with auth/gcp/login")

		return vaultResp.Auth.ClientToken, nil
	}

	return "", errors.New("GetVaultToken: configuration error, one of [githubAuth, googleAuth] must be set to true")
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

	auth := NewAuthClient()
	vc, err := NewVaultClient(ctx, auth, config)

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

		secret, err := vc.GetSecretFromVault(secretName)
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

	auth := NewAuthClient()
	vc, err := NewVaultClient(ctx, auth, config)

	if config.traceEnabled {
		var span *trace.Span
		ctx, span = trace.StartSpan(ctx, fmt.Sprintf("%s/GetSecrets", config.tracePrefix))
		defer span.End()
	}

	for _, secretName := range secretNames {
		log.Debug(fmt.Sprintf("secret= %s", secretNames))

		secretVersion, err := vc.GetSecretVersionFromVault(secretName)
		if err != nil {
			return fmt.Errorf("getting secret version: %v", err)
		}

		(*secretVersions)[secretName] = secretVersion
	}

	return nil
}

func getEnv(varName, defaultVal string) string {

	if value, isPresent := os.LookupEnv(varName); isPresent {
		return value
	}

	return defaultVal
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
