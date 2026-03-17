package k8s

import (
	"context"
	"testing"

	"github.com/matryer/is"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	testclient "k8s.io/client-go/kubernetes/fake"
)

func newTestSecret() *v1.Secret {
	return &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: "vault-secret",
		},
		Data: map[string][]byte{"key": []byte("value")},
	}
}

func TestClient(t *testing.T) {
	t.Run("create secret failure", func(t *testing.T) {
		is := is.New(t)
		secret := newTestSecret()
		c := &Client{Clientset: testclient.NewSimpleClientset(secret)}

		s, err := c.createSecret(context.Background(), secret)
		is.True(err != nil)
		is.Equal(s, nil)
	})

	t.Run("create secret sucess", func(t *testing.T) {
		is := is.New(t)
		secret := newTestSecret()
		c := &Client{Clientset: testclient.NewSimpleClientset()}

		s, err := c.createSecret(context.Background(), secret)
		is.NoErr(err)
		is.True(s != nil)
	})

	t.Run("update secret failure", func(t *testing.T) {
		is := is.New(t)
		secret := newTestSecret()
		c := &Client{Clientset: testclient.NewSimpleClientset()}

		err := c.updateSecret(context.Background(), secret)
		is.True(err != nil)
	})

	t.Run("update secret sucess", func(t *testing.T) {
		is := is.New(t)
		secret := newTestSecret()
		c := &Client{Clientset: testclient.NewSimpleClientset(secret)}

		err := c.updateSecret(context.Background(), secret)
		is.NoErr(err)
	})

	t.Run("get secret failure", func(t *testing.T) {
		is := is.New(t)
		secret := newTestSecret()
		c := &Client{Clientset: testclient.NewSimpleClientset()}

		s, err := c.getSecret(context.Background(), secret)
		is.True(err != nil)
		is.Equal(s, nil)
	})

	t.Run("get secret sucess", func(t *testing.T) {
		is := is.New(t)
		secret := newTestSecret()
		c := &Client{Clientset: testclient.NewSimpleClientset(secret)}

		s, err := c.getSecret(context.Background(), secret)
		is.NoErr(err)
		is.True(s != nil)
	})
}
