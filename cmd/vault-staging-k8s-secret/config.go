package main

import (
	"os"

	"go.uber.org/zap"
)

type config struct {
	defaultEngine string
	override      string
	k8sNamespace  string
}

func newConfig(lgr *zap.Logger) *config {

	// load required config from environment
	defaultEngine := os.Getenv("VAULT_SECRET")
	k8sNamespace := os.Getenv("K8S_NAMESPACE")

	req := []string{}
	if len(defaultEngine) < 1 {
		req = append(req, ("VAULT_SECRET"))
	}
	if len(k8sNamespace) < 1 {
		req = append(req, "K8S_NAMESPACE")
	}
	if len(req) > 0 {
		lgr.Fatal("bad configuration", zap.Strings("missing environment variables", req))
	}
	lgr.Debug("vault-path to the default engine", zap.String("VAULT_SECRET", defaultEngine))
	lgr.Debug("kubernetes namespace to apply secret", zap.String("K8S_NAMESPACE", k8sNamespace))

	// load optional config from environment
	override := os.Getenv("VAULT_SECRET_OVERRIDE")
	if len(override) < 1 {
		lgr.Info("No override engine provided.")
	} else {
		lgr.Debug("vault-path to an overriding engine", zap.String("VAULT_SECRET_OVERRIDE", override))
	}

	return &config{
		defaultEngine: defaultEngine,
		override:      override,
		k8sNamespace:  k8sNamespace,
	}
}
