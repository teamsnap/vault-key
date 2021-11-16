package vault_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/matryer/is"
	"github.com/teamsnap/vault-key/pkg/vault"
)

func TestGetSecrets(t *testing.T) {
	is := is.New(t)
	res := &map[string]map[string]string{}
	loginServer := vaultLoginServer()
	defer loginServer.Close()

	os.Setenv("VAULT_ADDR", loginServer.URL)

	err := vault.GetSecrets(context.Background(), res, []string{"my-key"})
	is.NoErr(err)
}

func TestGetSecretVersions(t *testing.T) {
	is := is.New(t)
	res := &map[string]int64{}
	loginServer := vaultLoginServer()
	defer loginServer.Close()

	os.Setenv("VAULT_ADDR", loginServer.URL)

	err := vault.GetSecretVersions(context.Background(), res, []string{"my-key"})
	is.NoErr(err)
}

func TestCreateSecret(t *testing.T) {
	is := is.New(t)
	loginServer := vaultLoginServer()
	defer loginServer.Close()

	os.Setenv("VAULT_ADDR", loginServer.URL)

	err := vault.CreateSecret(context.Background(), "staging/classic/dotenv", "key", "value")
	is.NoErr(err)
}

func TestUpdateSecret(t *testing.T) {
	is := is.New(t)
	loginServer := vaultLoginServer()
	defer loginServer.Close()

	os.Setenv("VAULT_ADDR", loginServer.URL)

	err := vault.UpdateSecret(context.Background(), "data", "my-key", "value")
	is.NoErr(err)
}

func TestDeleteSecret(t *testing.T) {
	is := is.New(t)
	loginServer := vaultLoginServer()
	defer loginServer.Close()

	os.Setenv("VAULT_ADDR", loginServer.URL)

	err := vault.DeleteSecret(context.Background(), "data", "my-key")
	is.NoErr(err)
}

func TestListEngines(t *testing.T) {
	is := is.New(t)
	loginServer := vaultLoginServer()
	defer loginServer.Close()

	os.Setenv("VAULT_ADDR", loginServer.URL)

	_, err := vault.ListEngines(context.Background(), "data")
	is.NoErr(err)
}

func vaultLoginServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/login") {
			json.NewEncoder(w).Encode(api.Secret{
				Auth: &api.SecretAuth{ClientToken: "vault-test-token"},
			})
			return
		}

		switch r.Method {
		case http.MethodGet:
			json.NewEncoder(w).Encode(
				&api.Secret{
					Data: map[string]interface{}{
						"data":            map[string]interface{}{"my-key": "bar"},
						"current_version": 1,
					},
				},
			)
		case http.MethodPost, http.MethodPut:
			var incoming map[string]interface{}
			json.NewDecoder(r.Body).Decode(&incoming)
		}
	}))
}
