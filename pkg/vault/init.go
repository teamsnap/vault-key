package vault

import (
	"context"
	log "github.com/sirupsen/logrus"
	"os"
)

var traceEnabled bool
var tracePrefix string
var project string
var serviceAccount string
var vaultRole string
var environment string

func initialize(ctx context.Context) {
	environment = os.Getenv("ENVIRONMENT")

	if environment == "production" {
		log.SetFormatter(&log.JSONFormatter{})
		log.SetLevel(log.WarnLevel)
	} else {
		log.SetFormatter(&log.TextFormatter{})
		log.SetLevel(log.TraceLevel)
	}

	getConfigFromEnv(ctx)
}
