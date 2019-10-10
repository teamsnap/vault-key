package vault

import (
	"context"
	"github.com/GoogleCloudPlatform/berglas/pkg/berglas"
	log "github.com/sirupsen/logrus"
	"go.opencensus.io/trace"
	"os"
	"strconv"
	"strings"
)

func getConfigFromEnv(ctx context.Context) {
	if traceEnabled {
		_, span := trace.StartSpan(ctx, tracePrefix+"/getConfigFromEnv")
		defer span.End()
	}

	traceEnabledString := getEnv(ctx, "TRACE_ENABLED", "false")
	traceEnabled, _ = strconv.ParseBool(traceEnabledString)
	log.Info("TRACE_ENABLED=" + traceEnabledString)

	tracePrefix = getEnv(ctx, "TRACE_PREFIX", "vault")
	log.Info("TRACE_PREFIX=" + tracePrefix)
	if tracePrefix == "" {
		log.Fatal("Error occurred getting TRACE_PREFIX variable from environment.")
	}

	vaultRole = getEnv(ctx, "VAULT_ROLE", "")
	log.Info("VAULT_ROLE=" + vaultRole)
	if vaultRole == "" {
		log.Fatal("You need to set the VAULT_ROLE environment variable")
	}

	// google injects this env var automatically in gcp environments
	project = getEnv(ctx, "GCLOUD_PROJECT", "")
	log.Info("GCLOUD_PROJECT=" + project)
	if project == "" {
		log.Fatal("You need to set the GCLOUD_PROJECT environment variable")
	}

	// google injects this env var automatically in gcp environments
	serviceAccount = getEnv(ctx, "FUNCTION_IDENTITY", "")
	log.Info("FUNCTION_IDENTITY=" + serviceAccount)
	if serviceAccount == "" {
		log.Fatal("You need to set the FUNCTION_IDENTITY environment variable")
	}
}

func getEnv(ctx context.Context, varName, defaultVal string) string {
	if traceEnabled {
		_, span := trace.StartSpan(ctx, tracePrefix+"/getEnv")
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
func getEncrEnvVar(ctx context.Context, n string) (string, error) {
	if traceEnabled {
		_, span := trace.StartSpan(ctx, tracePrefix+"/getEncrEnvVar")
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
