package vault

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/vault/api"
	"go.opencensus.io/trace"
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
	if vc.config.traceEnabled {
		var span *trace.Span
		vc.ctx, span = trace.StartSpan(vc.ctx, fmt.Sprintf("%s/gcp/GetVaultToken", vc.config.tracePrefix))
		defer span.End()
	}

	var err error
	a.iamService, err = iam.NewService(vc.ctx)
	if err != nil {
		return "", fmt.Errorf("getting new iam service: google: could not find default credentials")
	}

	err = a.generateSignedJWT(vc.ctx, vc.config)
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
func (a *gcpAuthClient) generateSignedJWT(ctx context.Context, c *config) error {
	if c.traceEnabled {
		var span *trace.Span
		ctx, span = trace.StartSpan(ctx, fmt.Sprintf("%s/generateSignedJWT", c.tracePrefix))
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
		return fmt.Errorf("marshaling payload: %v", err)
	}

	signJwtReq := &iam.SignJwtRequest{
		Payload: string(payloadBytes),
	}

	a.resp, err = a.iamService.Projects.ServiceAccounts.SignJwt(resourceName, signJwtReq).Do()
	if err != nil {
		return fmt.Errorf("sigining jwt: %v", err)
	}

	return nil
}

// gcpSaAuth takes signed JWT and sends login request to vault
func (a *gcpAuthClient) gcpSaAuth(vc *vaultClient) (*api.Secret, error) {
	if vc.config.traceEnabled {
		var span *trace.Span
		vc.ctx, span = trace.StartSpan(vc.ctx, fmt.Sprintf("%s/gcp/vaultLogin", vc.config.tracePrefix))
		defer span.End()
	}

	vaultResp, err := vc.client.Logical().Write(
		"auth/"+vc.config.gcpAuthPath+"/login",
		map[string]interface{}{
			"role": vc.config.vaultRole,
			"jwt":  a.resp.SignedJwt,
		})

	if err != nil {
		return nil, fmt.Errorf("logging into vault:%v", err)
	}

	return vaultResp, nil
}
