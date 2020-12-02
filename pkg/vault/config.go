package vault

import (
	"errors"
	"fmt"
	"strconv"

	log "github.com/sirupsen/logrus"
)

type config struct {
	project        string
	serviceAccount string
	githubOAuthToken string
	traceEnabled   bool
	tracePrefix    string
	vaultRole      string
}

func loadVaultEnvironment() (*config, error) {
	c := &config{}
	traceEnabledString := getEnv("TRACE_ENABLED", "false")
	c.traceEnabled, _ = strconv.ParseBool(traceEnabledString)

	c.tracePrefix = getEnv("TRACE_PREFIX", "vault")
	if c.tracePrefix == "" {
		return nil, errors.New("set the TRACE_PREFIX variable from environment")
	}

	c.vaultRole = getEnv("VAULT_ROLE", "")
	if c.vaultRole == "" {
		return nil, errors.New("set the VAULT_ROLE environment variable")
	}

	// google injects this env var automatically in gcp environments
	c.githubOAuthToken = getEnv("GITHUB_OAUTH_TOKEN", "")

	// google injects this env var automatically in gcp environments
	c.project = getEnv("GCLOUD_PROJECT", "")

	// google injects this env var automatically in gcp environments
	c.serviceAccount = getEnv("FUNCTION_IDENTITY", "")

	if c.githubOAuthToken == "" && (c.project == "" || c.serviceAccount == "") {
			return nil, errors.New("if using github auth, set the GITHUB_OAUTH_TOKEN. if using Google servive account auth, set the GCLOUD_PROJECT and FUNCTION_IDENTITY environment variables")
	}

	log.Info(fmt.Sprintf("TRACE_PREFIX=%s, VAULT_ROLE=%s, GCLOUD_PROJECT=%s, FUNCTION_IDENTITY=%s", c.tracePrefix, c.vaultRole, c.project, c.serviceAccount))

	return c, nil
}
