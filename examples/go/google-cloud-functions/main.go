package gcf

import (
	"context"
	"fmt"
	"log"

	"contrib.go.opencensus.io/exporter/stackdriver"
	"github.com/teamsnap/vault-key/pkg/vault"
	"go.opencensus.io/trace"
)

var env = map[string]map[string]string{}

var envMap = []string{
	"test/data/test",
}

// PubSubMessage represents the data passed by Google PubSub to trigger the
// function. We'll ignore this in the example.
type PubSubMessage struct {
	Data []byte `json:"data"`
}

func createTraceExporter() {
	// Create and register a OpenCensus Stackdriver Trace exporter.
	exporter, err := stackdriver.NewExporter(stackdriver.Options{})
	if err != nil {
		log.Fatal(err)
	}
	trace.RegisterExporter(exporter)

	// necessary or it won't take enough samples for us to notice
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})
}

// init runs during package initialization. this will only run during an
// an instance's cold start
func init() {
	ctx := context.Background()

	createTraceExporter()

	vault.GetSecrets(ctx, &env, envMap)
}

func vaultOnInit(ctx context.Context, m PubSubMessage) error {
	fmt.Println("Environment Values:", env)
	fmt.Println("hello = " + env["test/data/test"]["hello"])

	return nil
}
