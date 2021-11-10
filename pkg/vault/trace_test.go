package vault

import (
	"context"
	"testing"

	"github.com/matryer/is"
)

type mockTracer struct {
	spans map[string]bool
}

func (m *mockTracer) trace(name string) func() {
	m.spans[name] = true

	return func() {}
}
func TestTracer(t *testing.T) {
	secretKey, secretValue = "existing-key", "existing-value"

	cluster := createTestVault(t)
	defer cluster.Cleanup()

	rootVaultClient := cluster.Cores[0].Client
	vc := &vaultClient{
		config: githubConfig(t),
		ctx:    context.Background(),
		client: rootVaultClient,
	}

	t.Run("trace create", test_createTrace(vc))
	t.Run("trace delete", test_deleteTrace(vc))
	t.Run("trace update", test_updateTrace(vc))
}

func test_createTrace(vc *vaultClient) func(*testing.T) {
	return func(t *testing.T) {
		engine, k, v := "kv/data/create/foo", "new-key", "new-value"
		is := is.New(t)
		vc.tracer = &mockTracer{spans: map[string]bool{}}

		_, err := vc.Create(engine, k, v)
		is.NoErr(err)

		val, ok := vc.tracer.(*mockTracer)
		is.Equal(ok, true)
		is.Equal(val.spans, map[string]bool{"vault/Create": true, "vault/write": true})
	}
}

func test_deleteTrace(vc *vaultClient) func(*testing.T) {
	return func(t *testing.T) {
		engine, k := "kv/data/foo", secretKey
		is := is.New(t)
		vc.tracer = &mockTracer{spans: map[string]bool{}}

		_, err := vc.Delete(engine, k)
		is.NoErr(err)

		val, ok := vc.tracer.(*mockTracer)
		is.Equal(ok, true)
		is.Equal(val.spans, map[string]bool{"vault/Delete": true})
	}
}

func test_updateTrace(vc *vaultClient) func(*testing.T) {
	return func(t *testing.T) {
		engine, k, v := "kv/data/foo", secretKey, secretValue
		is := is.New(t)
		vc.tracer = &mockTracer{spans: map[string]bool{}}

		_, err := vc.Update(engine, k, v)
		is.NoErr(err)

		val, ok := vc.tracer.(*mockTracer)
		is.Equal(ok, true)
		is.Equal(val.spans, map[string]bool{"vault/Update": true, "vault/write": true})
	}
}
