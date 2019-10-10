package vault

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/GoogleCloudPlatform/berglas/pkg/berglas"
	log "github.com/sirupsen/logrus"
	"go.opencensus.io/trace"
)

func (v *vault) getConfigFromEnv(ctx context.Context) error {
	if v.traceEnabled {
		_, span := trace.StartSpan(ctx, fmt.Sprintf("%s/getConfigFromEnv", v.tracePrefix))
		defer span.End()
	}

	traceEnabledString := v.getEnv(ctx, "TRACE_ENABLED", "false")
	v.traceEnabled, _ = strconv.ParseBool(traceEnabledString)
	log.Info(fmt.Sprintf("TRACE_ENABLED=%s", traceEnabledString))

	v.tracePrefix = v.getEnv(ctx, "TRACE_PREFIX", "vault")
	log.Info(fmt.Sprintf("TRACE_PREFIX=%s", v.tracePrefix))
	if v.tracePrefix == "" {
		return errors.New("Error occurred getting TRACE_PREFIX variable from environment")
	}

	v.vaultRole = v.getEnv(ctx, "VAULT_ROLE", "")
	log.Info(fmt.Sprintf("VAULT_ROLE=%s", v.vaultRole))
	if v.vaultRole == "" {
		return errors.New("You need to set the VAULT_ROLE environment variable")
	}

	// google injects this env var automatically in gcp environments
	v.project = v.getEnv(ctx, "GCLOUD_PROJECT", "")
	log.Info(fmt.Sprintf("GCLOUD_PROJECT=%s", v.project))
	if v.project == "" {
		return errors.New("You need to set the GCLOUD_PROJECT environment variable")
	}

	// google injects this env var automatically in gcp environments
	v.serviceAccount = v.getEnv(ctx, "FUNCTION_IDENTITY", "")
	log.Info(fmt.Sprintf("FUNCTION_IDENTITY=%s", v.serviceAccount))
	if v.serviceAccount == "" {
		return errors.New("You need to set the FUNCTION_IDENTITY environment variable")
	}

	return nil
}

func (v *vault) getEnv(ctx context.Context, varName, defaultVal string) string {
	if v.traceEnabled {
		_, span := trace.StartSpan(ctx, fmt.Sprintf("%s/getEnv", v.tracePrefix))
		defer span.End()
	}

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
func (v *vault) getEncrEnvVar(ctx context.Context, n string) (string, error) {
	if v.traceEnabled {
		_, span := trace.StartSpan(ctx, fmt.Sprintf("%s/getEncrEnvVar", v.tracePrefix))
		defer span.End()
	}

	val := os.Getenv(n)
	if strings.HasPrefix(val, "berglas://") {
		if err := berglas.Replace(ctx, n); err != nil {
			return "", err
		}
	}

	return os.Getenv(n), nil
}
