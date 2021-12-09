package main

import (
	"context"
	"fmt"
	"os"
	"runtime"

	"github.com/teamsnap/vault-key/pkg/k8s"
	"github.com/teamsnap/vault-key/pkg/vault"
	"go.uber.org/automaxprocs/maxprocs"
	"go.uber.org/zap"
)

func main() {

	lgr, err := newLogger("vault-staging-k8s-sync")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer lgr.Sync()

	if err := run(context.Background(), lgr); err != nil {
		lgr.Error("run", zap.Error(err))
		lgr.Sync()
		os.Exit(1)
	}
}

func run(ctx context.Context, lgr *zap.Logger) error {

	// =========================================================================
	// GOMAXPROCS

	// Set the correct number of threads for the service
	// based on what is available either by the machine or quotas.
	if _, err := maxprocs.Set(); err != nil {
		return fmt.Errorf("maxprocs: %w", err)
	}
	lgr.Info("startup", zap.Int("GOMAXPROCS", runtime.GOMAXPROCS(0)))

	cfg := newConfig(lgr)

	engines := []string{cfg.defaultEngine}
	if cfg.override != "" {
		engines = append(engines, cfg.override)
	}

	var verifiedEngines []string
	for _, p := range engines {
		lgr.Debug("verifying engine exists for path", zap.String("path", p))
		if _, err := vault.ListEngines(ctx, p); err != nil {
			verifiedEngines = append(verifiedEngines, p)
		} else {
			lgr.Error("cannot verify engine exists at path", zap.Error(err))
		}
	}

	lgr.Info("getting vault secrets from verified engines", zap.Strings("verified-engines", verifiedEngines))
	secrets := map[string]map[string]string{}
	if err := vault.GetSecrets(ctx, &secrets, verifiedEngines); err != nil {
		return fmt.Errorf("cannot get secrets from vault: %w", err)
	}

	mergedSecrets := map[string]string{}
	for _, e := range verifiedEngines {
		if e == cfg.override {
			lgr.Info("creating overrides")
			mergedSecrets = secrets[cfg.override]
		}
	}

	lgr.Info("setting defaults")
	for k, v := range secrets[cfg.defaultEngine] {
		if _, ok := mergedSecrets[k]; !ok {
			mergedSecrets[k] = v
		} else {
			lgr.Debug("overriding default value for key", zap.String("key", k))
		}
	}

	lgr.Info("applying merged secrets to namespace", zap.String("namespace", cfg.k8sNamespace))
	if err := k8s.ApplySecret(&k8s.Secret{
		Secrets:   mergedSecrets,
		Namespace: cfg.k8sNamespace,
	}); err != nil {
		return err
	}

	return nil
}
