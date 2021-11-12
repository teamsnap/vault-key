package vault

import (
	"errors"
	"os"
	"strconv"
)

type config struct {
	project        string
	serviceAccount string
	githubToken    string
	traceEnabled   bool
	tracePrefix    string
	vaultRole      string
	gcpAuthPath    string
}

func loadVaultEnvironment() (*config, error) {
	c := &config{}

	traceEnabledString := getEnv("TRACE_ENABLED", "false")
	c.traceEnabled, _ = strconv.ParseBool(traceEnabledString)

	c.tracePrefix = getEnv("TRACE_PREFIX", "vault")
	if c.tracePrefix == "" {
		return nil, errors.New("set the TRACE_PREFIX variable from environment")
	}

	// Prefer github oauth token if available
	if token := getEnv("GITHUB_OAUTH_TOKEN", ""); len(token) > 0 {
		c.githubToken = token

		return c, nil
	}

	c.project = getEnv("GCLOUD_PROJECT", "")
	if c.project == "" {
		return nil, errors.New("set the GCLOUD_PROJECT environment variable")
	}

	// google injects this env var automatically in gcp environments
	c.serviceAccount = getEnv("FUNCTION_IDENTITY", "")
	if c.serviceAccount == "" {
		return nil, errors.New("set the FUNCTION_IDENTITY environment variable")
	}

	c.gcpAuthPath = getEnv("GCP_AUTH_PATH", "gcp")
	c.vaultRole = getEnv("VAULT_ROLE", "")
	if c.vaultRole == "" {
		return nil, errors.New("set the VAULT_ROLE environment variable")
	}

	return c, nil
}

func getEnv(varName, defaultVal string) string {
	if value, isPresent := os.LookupEnv(varName); isPresent {
		return value
	}

	return defaultVal
}
