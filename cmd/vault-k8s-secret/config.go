package main

import (
	"os"
	"strings"

	"github.com/teamsnap/vault-key/pkg/k8s"
	"go.uber.org/zap"
)

type config struct {
	engine            string
	defaultSecretPath string
	k8sNamespace      string
	k8sSecretName     string
}

func newConfig(lgr *zap.Logger) *config {

	// load required config from environment
	defaultSecret := getEnv("VAULT_SECRET", "")
	k8sNamespace := getEnv("K8S_NAMESPACE", "")
	k8sSecretName := getEnv("K8S_SECRET_NAME", k8s.DefaultSecretName)

	req := []string{}
	if len(defaultSecret) < 1 {
		req = append(req, "VAULT_SECRET")
	}
	if len(k8sNamespace) < 1 {
		req = append(req, "K8S_NAMESPACE")
	}
	if len(req) > 0 {
		lgr.Fatal("bad configuration", zap.Strings("missing environment variables", req))
	}
	lgr.Debug("vault-path to the default engine", zap.String("VAULT_SECRET", defaultSecret))
	lgr.Debug("kubernetes namespace to apply secret", zap.String("K8S_NAMESPACE", k8sNamespace))
	lgr.Debug("kubernetes secret name", zap.String("K8S_SECRET_NAME", k8sSecretName))

	return &config{
		engine:            translatePath(defaultSecret),
		defaultSecretPath: defaultSecret,
		k8sNamespace:      k8sNamespace,
		k8sSecretName:     k8sSecretName,
	}
}

func getEnv(varName, defaultVal string) string {
	if value, isPresent := os.LookupEnv(varName); isPresent {
		return value
	}
	return defaultVal
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
