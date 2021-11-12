package vault

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/vault/api"
	"google.golang.org/api/iam/v1"
)

type gcpAuthClient struct {
	iamService *iam.Service
	resp       *iam.SignJwtResponse
}

// NewAuthClient returns a new instance of an auth client
func NewGcpAuthClient() AuthClient {
	return &gcpAuthClient{}
}

func (a *gcpAuthClient) GetVaultToken(vc *vaultClient) (string, error) {
	vc.tracer.trace(fmt.Sprintf("%s/gcp/GetVaultToken", vc.config.tracePrefix))

	var err error
	a.iamService, err = iam.NewService(vc.ctx)
	if err != nil {
		return "", fmt.Errorf("getting new iam service: %w", err)
	}

	err = a.generateSignedJWT(vc)
	if err != nil {
		return "", fmt.Errorf("generating signed jwt, sigining jwt: Post")
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

	resourceName := fmt.Sprintf("projects/%s/serviceAccounts/%s", vc.config.project, vc.config.serviceAccount)
	jwtPayload := map[string]interface{}{
		"aud": "vault/" + vc.config.vaultRole,
		"sub": vc.config.serviceAccount,
		"exp": time.Now().Add(time.Minute * 10).Unix(),
	}

	payloadBytes, err := json.Marshal(jwtPayload)
	if err != nil {
		return fmt.Errorf("marshaling payload: %w", err)
	}

	signJwtReq := &iam.SignJwtRequest{
		Payload: string(payloadBytes),
	}

	a.resp, err = a.iamService.Projects.ServiceAccounts.SignJwt(resourceName, signJwtReq).Do()
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
