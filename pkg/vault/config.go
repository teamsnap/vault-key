package vault

import (
	"errors"
	"fmt"
	"strconv"

	log "github.com/sirupsen/logrus"
)

type config struct {
	googleAuth     bool
	project        string
	serviceAccount string
	githubToken    string
	githubAuth     bool
	traceEnabled   bool
	tracePrefix    string
	vaultRole      string
	gcpAuthPath    string
}

func loadVaultEnvironment() (*config, error) {
	c := &config{}
	c.githubAuth = false
	c.googleAuth = true

	traceEnabledString := getEnv("TRACE_ENABLED", "false")
	c.traceEnabled, _ = strconv.ParseBool(traceEnabledString)

	c.tracePrefix = getEnv("TRACE_PREFIX", "vault")
	if c.tracePrefix == "" {
		return nil, errors.New("set the TRACE_PREFIX variable from environment")
	}

	// If we have an oauth token set, override google authentication
	c.githubToken = getEnv("GITHUB_OAUTH_TOKEN", "")
	if c.githubToken != "" {
		c.githubAuth = true
		c.googleAuth = false
		log.Info(fmt.Sprintf("TRACE_PREFIX=%s, GITHUB_OAUTH_TOKEN=%s, GITHUB_AUTH=%t", c.tracePrefix, c.githubToken, c.githubAuth))
	}

	if c.googleAuth {
		// google injects this env var automatically in gcp environments
		c.project = getEnv("GCLOUD_PROJECT", "")
		if c.project == "" {
			return nil, errors.New("set the GCLOUD_PROJECT environment variable")
		}

		// google injects this env var automatically in gcp environments
		c.project = getEnv("GCLOUD_PROJECT", "")
		if c.project == "" {
			return nil, errors.New("set the GCLOUD_PROJECT environment variable")
		}

		// google injects this env var automatically in gcp environments
		c.serviceAccount = getEnv("FUNCTION_IDENTITY", "")
		if c.serviceAccount == "" {
			return nil, errors.New("set the FUNCTION_IDENTITY environment variable")
		}

		c.vaultRole = getEnv("VAULT_ROLE", "")
		if c.vaultRole == "" {
			return nil, errors.New("set the VAULT_ROLE environment variable")
		}
		log.Info(fmt.Sprintf("TRACE_PREFIX=%s, VAULT_ROLE=%s, GCLOUD_PROJECT=%s, FUNCTION_IDENTITY=%s, GOOGLE_AUTH=%t", c.tracePrefix, c.vaultRole, c.project, c.serviceAccount, c.googleAuth))
	}
	return c, nil
}
