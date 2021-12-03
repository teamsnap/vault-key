package main

import (
	"context"

	"github.com/teamsnap/vault-key/pkg/k8s"
	"github.com/teamsnap/vault-key/pkg/vault"
	"go.uber.org/zap"
)

func main() {
	run(context.Background(), newConfig())
}

func run(ctx context.Context, cfg *config) {
	engines := []string{cfg.defaultEngine, cfg.override}
	verifiedEngines := []string{}

	for _, p := range engines {
		cfg.lgr.Debug("verifying engine exists for path", zap.String("path", p))

		if _, err := vault.ListEngines(ctx, p); err != nil {
			verifiedEngines = append(verifiedEngines, p)
		} else {
			cfg.lgr.Error("cannot verify engine exists at path", zap.Error(err))
		}
	}

	cfg.lgr.Info("getting secrets from vault from the verified engines", zap.Strings("verified-engines", verifiedEngines))
	secrets := map[string]map[string]string{}
	if err := vault.GetSecrets(ctx, &secrets, verifiedEngines); err != nil {
		cfg.lgr.Fatal("cannot get secrets from vault", zap.Error(err))
	}

	cfg.lgr.Info("creating overrides")
	mergedSecrets := secrets[cfg.override]

	cfg.lgr.Info("setting defaults")
	for k, v := range secrets[cfg.defaultEngine] {
		if _, ok := mergedSecrets[k]; !ok {
			mergedSecrets[k] = v
		} else {
			cfg.lgr.Debug("overriding default value for key", zap.String("key", k))
		}
	}

	cfg.lgr.Info("applying merged secrets to namespace", zap.String("namespace", cfg.k8sNamespace))
	k8s.ApplySecret(&k8s.Secret{
		Secrets:   mergedSecrets,
		Namespace: cfg.k8sNamespace,
	})
}
