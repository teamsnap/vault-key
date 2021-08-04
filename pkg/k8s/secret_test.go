package k8s

import (
	"context"
	"testing"

	"github.com/matryer/is"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	testclient "k8s.io/client-go/kubernetes/fake"
)

func TestClient(t *testing.T) {
	s := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: "vault-secret",
		},
		Data: map[string][]byte{"key": []byte("value")},
	}

	cases := []struct {
		secret *v1.Secret
		client client
		isErr  bool
	}{
		{
			secret: s,
			client: client{clientset: testclient.NewSimpleClientset(s)},
			isErr:  true,
		},
		{
			secret: s,
			client: client{clientset: testclient.NewSimpleClientset()},
			isErr:  false,
		},
	}

	for _, c := range cases {
		t.Run("create secret", testCreateSecret(c.client, c.secret, c.isErr))
		t.Run("update secret", testUpdateSecret(c.client, c.secret))
		t.Run("get secret", testGetSecret(c.client, c.secret))
	}
}

func testCreateSecret(c client, secret *v1.Secret, isErr bool) func(*testing.T) {
	return func(t *testing.T) {
		is := is.New(t)

		_, err := c.createSecret(context.Background(), secret)
		if isErr {
			is.True(err != nil)
		} else {
			is.NoErr(err)
		}
	}
}

func testUpdateSecret(c client, secret *v1.Secret) func(*testing.T) {
	return func(t *testing.T) {
		is := is.New(t)
		err := c.updateSecret(context.Background(), secret)
		is.NoErr(err)
	}
}

func testGetSecret(c client, secret *v1.Secret) func(*testing.T) {
	return func(t *testing.T) {
		is := is.New(t)
		s, err := c.getSecret(context.Background(), secret)
		is.NoErr(err)
		is.Equal(s.Data, secret.Data)
	}
}
