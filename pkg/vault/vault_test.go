package vault

import (
	"context"
	"os"
	"testing"

	"github.com/matryer/is"
)

func TestGetSecrets(t *testing.T) {
	is := is.New(t)
	res := &map[string]map[string]string{}
	loginServer := vaultLoginServer()
	defer loginServer.Close()

	os.Setenv("VAULT_ADDR", loginServer.URL)

	err := GetSecrets(context.Background(), res, []string{"my-key"})
	is.NoErr(err)
}

func TestGetSecretVersions(t *testing.T) {
	is := is.New(t)
	res := &map[string]int64{}
	loginServer := vaultLoginServer()
	defer loginServer.Close()

	os.Setenv("VAULT_ADDR", loginServer.URL)

	err := GetSecretVersions(context.Background(), res, []string{"my-key"})
	is.NoErr(err)
}
