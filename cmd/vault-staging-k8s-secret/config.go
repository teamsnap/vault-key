package main

import (
	"os"

	"go.uber.org/zap"
)

type config struct {
	defaultEngine string
	override      string
	k8sNamespace  string
	lgr           *zap.Logger
}

func newConfig() *config {
	lgr := newLogger()

	defaultEngine := os.Getenv("VAULT_SECRET")
	lgr.Debug("vault-path to the default engine", zap.String("VAULT_SECRET", defaultEngine))
	if len(defaultEngine) < 1 {
		lgr.Fatal("You need to set VAULT_SECRET environment variable.")
	}

	override := os.Getenv("VAULT_SECRET_OVERRIDE")
	lgr.Debug("vault-path to an overriding engine", zap.String("VAULT_SECRET_OVERRIDE", override))
	if len(override) < 1 {
		lgr.Info("No override engine provided.")
	}

	k8sNamespace := os.Getenv("K8S_NAMESPACE")
	lgr.Debug("kubernetes namespace to apply secret", zap.String("K8S_NAMESPACE", k8sNamespace))
	if len(k8sNamespace) < 1 {
		lgr.Fatal("You need to set K8S_NAMESPACE environment variable.")
	}

	return &config{
		defaultEngine: defaultEngine,
		override:      override,
		k8sNamespace:  k8sNamespace,
		lgr:           lgr,
	}
}
