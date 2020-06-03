package vault

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/vault/api"
	log "github.com/sirupsen/logrus"
	"go.opencensus.io/trace"
	"google.golang.org/api/iam/v1"
)

type gcpAuthClient struct {
	iamService *iam.Service
	resp       *iam.SignJwtResponse
}

// NewAuthClient returns a new instance of an auth client
func NewAuthClient() AuthClient {
	a := &gcpAuthClient{}

	return a
}

// GetVaultToken uses a service account to get a vault auth token
func (a *gcpAuthClient) GetVaultToken(vc *vaultClient) (string, error) {
	if vc.config.traceEnabled {
		var span *trace.Span
		vc.ctx, span = trace.StartSpan(vc.ctx, fmt.Sprintf("%s/getVaultToken", vc.config.tracePrefix))
		defer span.End()
	}

	var err error
	a.iamService, err = iam.NewService(vc.ctx)
	if err != nil {
		log.Debugf("initialze client: getting new iam service: %v", err)
		return "", fmt.Errorf("getting new iam service: google: could not find default credentials")
	}

	log.Debug("Successfully created IAM client")

	err = a.generateSignedJWT(vc.ctx, vc.config)
	if err != nil {
		log.Debugf("generating signed jwt: %v", err)
		return "", fmt.Errorf("generating signed jwt, sigining jwt: Post")
	}

	log.Debug("Successfully generated signed JWT")

	vaultResp, err := a.openVault(vc)
	if err != nil {
		return "", err
	}

	log.Debug("Successfully logged into Vault with auth/gcp/login")

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

// openVault takes signed JWT and sends login request to vault
func (a *gcpAuthClient) openVault(vc *vaultClient) (*api.Secret, error) {
	if vc.config.traceEnabled {
		var span *trace.Span
		vc.ctx, span = trace.StartSpan(vc.ctx, fmt.Sprintf("%s/vaultLogin", vc.config.tracePrefix))
		defer span.End()
	}

	vaultResp, err := vc.client.Logical().Write(
		"auth/gcp/login",
		map[string]interface{}{
			"role": vc.config.vaultRole,
			"jwt":  a.resp.SignedJwt,
		})

	if err != nil {
		return nil, fmt.Errorf("logging into vault:%v", err)
	}

	return vaultResp, nil
}
