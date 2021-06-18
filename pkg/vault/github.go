package vault

import (
	"fmt"

	"github.com/hashicorp/vault/api"
	log "github.com/sirupsen/logrus"
	"go.opencensus.io/trace"
)

type githubAuthClient struct {
}

// NewAuthClient returns a new instance of an auth client
func NewGithubAuthClient() AuthClient {
	a := &githubAuthClient{}
	return a
}
func (a *githubAuthClient) GetVaultToken(vc *vaultClient) (string, error) {
	vaultResp, err := a.githubVaultAuth(vc)
	if err != nil {
		return "", err
	}

	log.Debug("Successfully logged into Vault with auth/github/login")

	return vaultResp.Auth.ClientToken, nil
}

// githubVaultAuth takes GitHub access token and sends login request to vault
func (a *githubAuthClient) githubVaultAuth(vc *vaultClient) (*api.Secret, error) {
	if vc.config.traceEnabled {
		var span *trace.Span
		vc.ctx, span = trace.StartSpan(vc.ctx, fmt.Sprintf("%s/github/vaultLogin", vc.config.tracePrefix))
		defer span.End()
	}

	vaultResp, err := vc.client.Logical().Write(
		"auth/github/login",
		map[string]interface{}{
			"token": vc.config.githubToken,
		})

	if err != nil {
		return nil, fmt.Errorf("logging into vault with github:%v", err)
	}

	return vaultResp, nil
}
