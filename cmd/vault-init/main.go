package main

import (
	"context"
	"github.com/teamsnap/vault-key/pkg/vault"
	log "github.com/sirupsen/logrus"
	"os"
)

var env = map[string]map[string]string{}

// writeStringToFile takes a file path and a content string and writes the
// content to a file with the given path
func writeStringToFile(filePath, str string) {
	f, err := os.Create(filePath)
	if err != nil {
		log.Fatal("Error writing to file", filePath, err)
	}

	defer f.Close()

	f.WriteString(str)

	// flush
	f.Sync()
}

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

	writeStringToFile("/usr/share/vault/data/secrets", GenerateDotEnv(DotEnvVariables{Secrets: env[vaultSecret]}))
}
