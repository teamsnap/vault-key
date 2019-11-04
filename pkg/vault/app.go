package vault

import (
	"context"

	"github.com/hashicorp/vault/api"
)

func newApp() *App {
	a := &App{}
	return a
}

// App is the thing
type App struct {
	traceEnabled   bool
	tracePrefix    string
	project        string
	serviceAccount string
	vaultRole      string
	environment    string
	client         *api.Client
	ctx            context.Context
}
