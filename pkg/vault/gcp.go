package vault

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/vault/api"
	log "github.com/sirupsen/logrus"
	"go.opencensus.io/trace"
	"google.golang.org/api/iam/v1"
	"time"
)

// generateSignedJWT returns a signed JWT response using IAM
func generateSignedJWT(ctx context.Context, iamClient *iam.Service, project, serviceAccount, vaultRole string) *iam.SignJwtResponse {
	if traceEnabled {
		_, span := trace.StartSpan(ctx, tracePrefix+"/generateSignedJWT")
		defer span.End()
	}

	resourceName := fmt.Sprintf("projects/%s/serviceAccounts/%s", project, serviceAccount)
	jwtPayload := map[string]interface{}{
		"aud": "vault/" + vaultRole,
		"sub": serviceAccount,
		"exp": time.Now().Add(time.Minute * 10).Unix(),
	}

	payloadBytes, err := json.Marshal(jwtPayload)
	if err != nil {
		log.Fatal(err)
	}
	signJwtReq := &iam.SignJwtRequest{
		Payload: string(payloadBytes),
	}

	resp, err := iamClient.Projects.ServiceAccounts.SignJwt(resourceName, signJwtReq).Do()
	if err != nil {
		log.Fatal(err)
	}

	log.Debug("Successfully generated signed JWT")

	return resp
}

// vaultLogin takes signed JWT and sends login request to vault
func vaultLogin(ctx context.Context, resp *iam.SignJwtResponse, vaultRole string) *api.Secret {
	if traceEnabled {
		_, span := trace.StartSpan(ctx, tracePrefix+"/vaultLogin")
		defer span.End()
	}

	vaultClient, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		log.Fatal(err)
	}

	vaultResp, err := vaultClient.Logical().Write(
		"auth/gcp/login",
		map[string]interface{}{
			"role": vaultRole,
			"jwt":  resp.SignedJwt,
		})

	if err != nil {
		log.Fatal(err)
	}

	log.Debug("Successfully logged into Vault with auth/gcp/login")

	return vaultResp
}

// getVaultToken uses a service account to get a vault auth token
func getVaultToken(ctx context.Context, project, serviceAccount, vaultRole string) string {
	if traceEnabled {
		_, span := trace.StartSpan(ctx, tracePrefix+"/getVaultToken")
		defer span.End()
	}

	var (
		err       error
		iamClient *iam.Service
	)

	iamClient, err = iam.NewService(ctx)

	if err != nil {
		log.Fatal(err)
	}

	log.Debug("Successfully created IAM client")

	resp := generateSignedJWT(ctx, iamClient, project, serviceAccount, vaultRole)
	vaultResp := vaultLogin(ctx, resp, vaultRole)

	log.Debug("Successfully got vault token:", vaultResp.Auth.ClientToken)

	return vaultResp.Auth.ClientToken
}
