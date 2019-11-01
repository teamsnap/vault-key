package vault

import (
	"context"

	"github.com/hashicorp/vault/api"
)

func newApp() *app {
	a := &app{}
	return a
}

// Vaults is the thing
type app struct {
	traceEnabled   bool
	tracePrefix    string
	project        string
	serviceAccount string
	vaultRole      string
	environment    string
	client         *api.Client
	ctx            context.Context
}
