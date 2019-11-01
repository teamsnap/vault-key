package vault

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/vault/api"
	log "github.com/sirupsen/logrus"
	"go.opencensus.io/trace"
	"google.golang.org/api/iam/v1"
)

// generateSignedJWT returns a signed JWT response using IAM
func (a *app) generateSignedJWT(iamClient *iam.Service) (*iam.SignJwtResponse, error) {
	if a.traceEnabled {
		var span *trace.Span
		a.ctx, span = trace.StartSpan(a.ctx, fmt.Sprintf("%s/generateSignedJWT", a.tracePrefix))
		defer span.End()
	}

	resourceName := fmt.Sprintf("projects/%s/serviceAccounts/%s", a.project, a.serviceAccount)

	jwtPayload := map[string]interface{}{
		"aud": "vault/" + a.vaultRole,
		"sub": a.serviceAccount,
		"exp": time.Now().Add(time.Minute * 10).Unix(),
	}

	payloadBytes, err := json.Marshal(jwtPayload)
	if err != nil {
		return nil, fmt.Errorf("Failed to marshal json Error: %v", err)
	}
	signJwtReq := &iam.SignJwtRequest{
		Payload: string(payloadBytes),
	}

	resp, err := iamClient.Projects.ServiceAccounts.SignJwt(resourceName, signJwtReq).Do()
	if err != nil {
		return nil, fmt.Errorf("Failed to sign jwt: %v", err)
	}
	log.Debug("Successfully generated signed JWT")

	return resp, nil
}

// vaultLogin takes signed JWT and sends login request to vault
func (a *app) vaultLogin(resp *iam.SignJwtResponse) (*api.Secret, error) {
	if a.traceEnabled {
		var span *trace.Span
		a.ctx, span = trace.StartSpan(a.ctx, fmt.Sprintf("%s/vaultLogin", a.tracePrefix))
		defer span.End()
	}

	vaultClient, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		return nil, fmt.Errorf("Failed to create new vault api client: %v", err)
	}

	vaultResp, err := vaultClient.Logical().Write(
		"auth/gcp/login",
		map[string]interface{}{
			"role": a.vaultRole,
			"jwt":  resp.SignedJwt,
		})
	if err != nil {
		return nil, fmt.Errorf("Failed to login to vault with auth/gcp/login: %v", err)
	}
	log.Debug("Successfully logged into Vault with auth/gcp/login")

	return vaultResp, nil
}

// getVaultToken uses a service account to get a vault auth token
func (a *app) getVaultToken() (string, error) {
	if a.traceEnabled {
		var span *trace.Span
		a.ctx, span = trace.StartSpan(a.ctx, fmt.Sprintf("%s/getVaultToken", a.tracePrefix))
		defer span.End()
	}

	iamClient, err := iam.NewService(a.ctx)
	if err != nil {
		return "", err
	}
	log.Debug("Successfully created IAM client")

	resp, err := a.generateSignedJWT(iamClient)
	if err != nil {
		return "", err
	}

	vaultResp, err := a.vaultLogin(resp)
	if err != nil {
		return "", err
	}
	log.Debug("Successfully got vault token:", vaultResp.Auth.ClientToken)

	return vaultResp.Auth.ClientToken, nil
}
