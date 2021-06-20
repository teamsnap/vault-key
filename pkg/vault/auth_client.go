package vault

import (
	"os"

	log "github.com/sirupsen/logrus"
)

// AuthClient is a type that satifies the necesary authorization layer for a vault client.
type AuthClient interface {
	GetVaultToken(vc *vaultClient) (string, error)
}

func NewAuthClient(c *config) AuthClient {
	switch {
	case len(c.githubToken) > 0:
		return NewGithubAuthClient()
	case len(c.project) > 0:
		return NewGcpAuthClient()
	default:
		log.Error("GetVaultToken: configuration error, one of [githubAuth, googleAuth] must be set to true")
		os.Exit(1)
	}

	return nil
}
