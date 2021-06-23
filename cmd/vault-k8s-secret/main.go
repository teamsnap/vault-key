package main

import (
	"context"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/teamsnap/vault-key/pkg/k8s"
	"github.com/teamsnap/vault-key/pkg/vault"
)

var env = map[string]map[string]string{}

func main() {
	ctx := context.Background()

	vaultSecret := os.Getenv("VAULT_SECRET")
	log.Info("VAULT_SECRET=" + vaultSecret)
	if vaultSecret == "" {
		log.Fatal("You need to set VAULT_SECRET environment variable.")
	}

	k8sNamespace := os.Getenv("K8S_NAMESPACE")
	log.Info("K8S_NAMESPACE=" + k8sNamespace)
	if k8sNamespace == "" {
		log.Fatal("You need to set K8S_NAMESPACE environment variable.")
	}

	var envArr = []string{
		vaultSecret,
	}

	vault.GetSecrets(ctx, &env, envArr)

	k8sSecret := &k8s.Secret{
		Secrets:   env[vaultSecret],
		Namespace: k8sNamespace,
	}

	k8s.ApplySecret(k8sSecret)
}
