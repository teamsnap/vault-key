package vault

import (
	"fmt"

	"github.com/hashicorp/vault/api"
)

type githubAuthClient struct {
}

// NewAuthClient returns a new instance of an auth client
func NewGithubAuthClient() AuthClient {
	return &githubAuthClient{}
}
func (a *githubAuthClient) GetVaultToken(vc *vaultClient) (string, error) {
	vc.tracer.trace(fmt.Sprintf("%s/github/GetVaultToken", vc.config.tracePrefix))

	vaultResp, err := a.githubVaultAuth(vc)
	if err != nil {
		return "", err
	}

	return vaultResp.Auth.ClientToken, nil
}

// githubVaultAuth takes GitHub access token and sends login request to vault
func (a *githubAuthClient) githubVaultAuth(vc *vaultClient) (*api.Secret, error) {
	vc.tracer.trace(fmt.Sprintf("%s/github/vaultLogin", vc.config.tracePrefix))

	vaultResp, err := vc.client.Logical().Write(
		"auth/github/login",
		map[string]interface{}{
			"token": vc.config.githubToken,
		})

	if err != nil {
		return nil, fmt.Errorf("logging into vault with github:%w", err)
	}

	return vaultResp, nil
}
