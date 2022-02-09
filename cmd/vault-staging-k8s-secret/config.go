package main

import (
	"os"
	"strings"

	"go.uber.org/zap"
)

type config struct {
	engine             string
	defaultSecretPath  string
	overrideSecretPath string
	k8sNamespace       string
}

func newConfig(lgr *zap.Logger) *config {

	// load required config from environment
	defaultSecret := os.Getenv("VAULT_SECRET")
	k8sNamespace := os.Getenv("K8S_NAMESPACE")

	req := []string{}
	if len(defaultSecret) < 1 {
		req = append(req, ("VAULT_SECRET"))
	}
	if len(k8sNamespace) < 1 {
		req = append(req, "K8S_NAMESPACE")
	}
	if len(req) > 0 {
		lgr.Fatal("bad configuration", zap.Strings("missing environment variables", req))
	}
	lgr.Debug("vault-path to the default engine", zap.String("VAULT_SECRET", defaultSecret))
	lgr.Debug("kubernetes namespace to apply secret", zap.String("K8S_NAMESPACE", k8sNamespace))

	// load optional config from environment
	override := os.Getenv("VAULT_SECRET_OVERRIDE")
	if len(override) < 1 {
		lgr.Info("No override secret provided.")
	} else {
		lgr.Debug("vault-path to an overriding secret", zap.String("VAULT_SECRET_OVERRIDE", override))
	}

	return &config{
		engine:             translatePath(defaultSecret),
		defaultSecretPath:  defaultSecret,
		overrideSecretPath: override,
		k8sNamespace:       k8sNamespace,
	}
}

// translatePath is a helper that converts a path to a vault secret
// to a truncated path needed for the vault api engines list function
//
// ie staging/applications/data/foo/dotenv -> staging/applications/metadata/foo
func translatePath(path string) string {
	strs := strings.Split(path, "/")

	for i, s := range strs {
		if s == "data" {
			strs[i] = "metadata"
		}
	}

	return strings.Join(strs[:len(strs)-1], "/")
}

// getSecret is a helper that takes a path to a vault secret
// and returns the secret
//
// staging/applications/data/foo/dotenv -> foo
func getSecret(path string) string {
	strs := strings.Split(path, "/")
	return strs[len(strs)-1]
}
