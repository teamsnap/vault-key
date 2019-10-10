package main

import (
	"context"
	"github.com/teamsnap/vault/pkg/k8s"
	"github.com/teamsnap/vault/pkg/vault"
	log "github.com/sirupsen/logrus"
	"os"
)

var env = map[string]map[string]string{}

func main() {
	ctx := context.Background()

	vaultSecret := os.Getenv("VAULT_SECRET")
	log.Info("VAULT_SECRET=" + vaultSecret)
	if vaultSecret == "" {
		log.Fatal("You need to set VAULT_SECRET environment variable.")
	}

	var envArr = []string{
		vaultSecret,
	}

	vault.GetSecrets(ctx, &env, envArr)

	k8sSecret := &k8s.Secret{
		Secrets: env[vaultSecret],
	}

	k8s.ApplySecret(k8sSecret)
}
