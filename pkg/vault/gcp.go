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

// getVaultToken uses a service account to get a vault auth token
func getVaultToken(c *vaultClient) (string, error) {
	if c.traceEnabled {
		var span *trace.Span
		c.ctx, span = trace.StartSpan(c.ctx, fmt.Sprintf("%s/getVaultToken", c.tracePrefix))
		defer span.End()
	}

	iamClient, err := iam.NewService(c.ctx)
	if err != nil {
		return "", fmt.Errorf("Error getting vault token: %v", err)
	}
	log.Debug("Successfully created IAM client")

	resp, err := generateSignedJWT(c, iamClient)
	if err != nil {
		return "", err
	}
	log.Debug("Successfully generated signed JWT")

	vaultResp, err := vaultLogin(c, resp)
	if err != nil {
		return "", err
	}
	log.Debug("Successfully logged into Vault with auth/gcp/login")

	return vaultResp.Auth.ClientToken, nil
}

// generateSignedJWT returns a signed JWT response using IAM
func generateSignedJWT(c *vaultClient, iamClient *iam.Service) (*iam.SignJwtResponse, error) {
	if c.traceEnabled {
		var span *trace.Span
		c.ctx, span = trace.StartSpan(c.ctx, fmt.Sprintf("%s/generateSignedJWT", c.tracePrefix))
		defer span.End()
	}

	resourceName := fmt.Sprintf("projects/%s/serviceAccounts/%s", c.project, c.serviceAccount)
	jwtPayload := map[string]interface{}{
		"aud": "vault/" + c.vaultRole,
		"sub": c.serviceAccount,
		"exp": time.Now().Add(time.Minute * 10).Unix(),
	}

	payloadBytes, err := json.Marshal(jwtPayload)
	if err != nil {
		return nil, fmt.Errorf("Error marshaling payload: %v", err)
	}

	signJwtReq := &iam.SignJwtRequest{
		Payload: string(payloadBytes),
	}

	resp, err := iamClient.Projects.ServiceAccounts.SignJwt(resourceName, signJwtReq).Do()
	if err != nil {
		return nil, fmt.Errorf("Error sigining jwt: %v", err)
	}

	return resp, nil
}

// vaultLogin takes signed JWT and sends login request to vault
func vaultLogin(c *vaultClient, resp *iam.SignJwtResponse) (vaultResp *api.Secret, err error) {
	if c.traceEnabled {
		var span *trace.Span
		c.ctx, span = trace.StartSpan(c.ctx, fmt.Sprintf("%s/vaultLogin", c.tracePrefix))
		defer span.End()
	}

	vaultResp, err = c.client.Logical().Write(
		"auth/gcp/login",
		map[string]interface{}{
			"role": c.vaultRole,
			"jwt":  resp.SignedJwt,
		})

	if err != nil {
		return nil, fmt.Errorf("Error getting logging into vault:%v", err)
	}

	return vaultResp, nil
}
