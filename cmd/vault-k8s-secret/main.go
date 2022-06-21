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

	lgr, err := newLogger("vault-k8s-sync")
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

	secretsToApply := []string{cfg.defaultSecretPath}

	var vaultSecrets []string
	lgr.Debug("getting secrets for engine", zap.String("engine", cfg.engine))

	// If a secret to apply does not exist in vault, do not attempt to apply it
	var verifiedSecretPaths []string
	for _, s := range secretsToApply {
		found := false

		for _, vs := range vaultSecrets {
			if getSecret(s) == vs {
				lgr.Debug("secret exists", zap.String("secret", s))
				verifiedSecretPaths = append(verifiedSecretPaths, s)
				found = true
				break
			}
		}

		if !found {
			lgr.Debug("secret does not exist", zap.String("secret", s))
		}
	}

	lgr.Info("getting vault secrets from verified secret paths", zap.Strings("verified-secret-paths", verifiedSecretPaths))
	secrets := map[string]map[string]string{}
	if err := vault.GetSecrets(ctx, &secrets, verifiedSecretPaths); err != nil {
		return fmt.Errorf("cannot get secrets from vault: %w", err)
	}

	lgr.Info("applying merged secrets to namespace", zap.String("namespace", cfg.k8sNamespace), zap.Int("number of secrets", len(secrets[cfg.defaultSecretPath])))
	if err := k8s.ApplySecret(
		&k8s.Secret{
			Secrets:   secrets[cfg.defaultSecretPath],
			Namespace: cfg.k8sNamespace,
		}); err != nil {
		return fmt.Errorf("unable to apply secrets to namespace %s: %w", cfg.k8sNamespace, err)
	}

	return nil
}
