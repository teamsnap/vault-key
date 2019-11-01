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

func (a *App) getConfigFromEnv() error {
	if a.traceEnabled {
		var span *trace.Span
		a.ctx, span = trace.StartSpan(a.ctx, fmt.Sprintf("%s/getConfigFromEnv", a.tracePrefix))
		defer span.End()
	}

	traceEnabledString := a.getEnv("TRACE_ENABLED", "false")
	a.traceEnabled, _ = strconv.ParseBool(traceEnabledString)
	log.Info(fmt.Sprintf("TRACE_ENABLED=%s", traceEnabledString))

	a.tracePrefix = a.getEnv("TRACE_PREFIX", "vault")
	log.Info(fmt.Sprintf("TRACE_PREFIX=%s", a.tracePrefix))
	if a.tracePrefix == "" {
		return errors.New("Error occurred getting TRACE_PREFIX variable from environment")
	}

	a.vaultRole = a.getEnv("VAULT_ROLE", "")
	log.Info(fmt.Sprintf("VAULT_ROLE=%s", a.vaultRole))
	if a.vaultRole == "" {
		return errors.New("You need to set the VAULT_ROLE environment variable")
	}

	// google injects this env var automatically in gcp environments
	a.project = a.getEnv("GCLOUD_PROJECT", "")
	log.Info(fmt.Sprintf("GCLOUD_PROJECT=%s", a.project))
	if a.project == "" {
		return errors.New("You need to set the GCLOUD_PROJECT environment variable")
	}

	// google injects this env var automatically in gcp environments
	a.serviceAccount = a.getEnv("FUNCTION_IDENTITY", "")
	log.Info(fmt.Sprintf("FUNCTION_IDENTITY=%s", a.serviceAccount))
	if a.serviceAccount == "" {
		return errors.New("You need to set the FUNCTION_IDENTITY environment variable")
	}

	return nil
}

func (a *App) getEnv(varName, defaultVal string) string {
	if a.traceEnabled {
		var span *trace.Span
		a.ctx, span = trace.StartSpan(a.ctx, fmt.Sprintf("%s/getEnv", a.tracePrefix))
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
func (a *App) getEncrEnvVar(ctx context.Context, n string) (string, error) {
	if a.traceEnabled {
		var span *trace.Span
		a.ctx, span = trace.StartSpan(ctx, fmt.Sprintf("%s/getEncrEnvVar", a.tracePrefix))
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
