package vault

import (
	"encoding/json"
	"fmt"
	"time"

	credentials "cloud.google.com/go/iam/credentials/apiv1"
	"cloud.google.com/go/iam/credentials/apiv1/credentialspb"
	"github.com/hashicorp/vault/api"
)

type gcpAuthClient struct {
	credentialsClient *credentials.IamCredentialsClient
	resp              *credentialspb.SignJwtResponse
}

// NewAuthClient returns a new instance of an auth client
func NewGcpAuthClient() AuthClient {
	return &gcpAuthClient{}
}

func (a *gcpAuthClient) GetVaultToken(vc *vaultClient) (string, error) {
	vc.tracer.trace(fmt.Sprintf("%s/gcp/GetVaultToken", vc.config.tracePrefix))

	var err error
	a.credentialsClient, err = credentials.NewIamCredentialsClient(vc.ctx)
	if err != nil {
		return "", fmt.Errorf("getting new iam credentials client: %w", err)
	}

	err = a.generateSignedJWT(vc)
	if err != nil {
		return "", fmt.Errorf("generate signed jwt:  %w", err)
	}

	vaultResp, err := a.gcpSaAuth(vc)
	if err != nil {
		return "", err
	}

	return vaultResp.Auth.ClientToken, nil
}

// generateSignedJWT returns a signed JWT response using IAM
func (a *gcpAuthClient) generateSignedJWT(vc *vaultClient) error {
	vc.tracer.trace(fmt.Sprintf("%s/gcp/generateSignedJWT", vc.config.tracePrefix))

	// `projects/-/serviceAccounts/{ACCOUNT_EMAIL_OR_UNIQUEID}`. The `-` wildcard
	// character is required; replacing it with a project ID is invalid.
	// https://pkg.go.dev/cloud.google.com/go/iam@v1.1.8/credentials/apiv1/credentialspb#SignJwtRequest
	jwtPayload := map[string]any{
		"aud": "vault/" + vc.config.vaultRole,
		"sub": vc.config.serviceAccount,
		"exp": time.Now().Add(time.Minute * 10).Unix(),
	}

	payloadBytes, err := json.Marshal(jwtPayload)
	if err != nil {
		return fmt.Errorf("marshaling payload: %w", err)
	}

	signJwtReq := &credentialspb.SignJwtRequest{
		Name:      fmt.Sprintf("projects/-/serviceAccounts/%s", vc.config.serviceAccount),
		Delegates: []string{fmt.Sprintf("projects/-/serviceAccounts/%s", vc.config.serviceAccount)},
		Payload:   string(payloadBytes),
	}

	a.resp, err = a.credentialsClient.SignJwt(vc.ctx, signJwtReq)
	if err != nil {
		return fmt.Errorf("sigining jwt: %w", err)
	}

	return nil
}

// gcpSaAuth takes signed JWT and sends login request to vault
func (a *gcpAuthClient) gcpSaAuth(vc *vaultClient) (*api.Secret, error) {
	vc.tracer.trace(fmt.Sprintf("%s/gcp/vaultLogin", vc.config.tracePrefix))

	vaultResp, err := vc.client.Logical().Write(
		"auth/"+vc.config.gcpAuthPath+"/login",
		map[string]interface{}{
			"role": vc.config.vaultRole,
			"jwt":  a.resp.SignedJwt,
		})

	if err != nil {
		return nil, fmt.Errorf("logging into vault:%w", err)
	}

	return vaultResp, nil
}
