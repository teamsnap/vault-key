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
	secret := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: "vault-secret",
		},
		Data: map[string][]byte{"key": []byte("value")},
	}

	cases := []struct {
		name   string
		secret *v1.Secret
		client *Client
		isErr  bool
	}{
		{
			name:   "failure",
			secret: secret,
			client: &Client{Clientset: testclient.NewSimpleClientset(secret)},
			isErr:  true,
		},
		{
			name:   "sucess",
			secret: secret,
			client: &Client{Clientset: testclient.NewSimpleClientset()},
			isErr:  false,
		},
	}

	for _, c := range cases {
		t.Run("create secret "+c.name, testCreateSecret(c.client, c.secret, c.isErr))
		t.Run("update secret "+c.name, testUpdateSecret(c.client, c.secret, c.isErr))
		t.Run("get secret "+c.name, testGetSecret(c.client, c.secret, c.isErr))
	}
}

func testCreateSecret(c *Client, secret *v1.Secret, isErr bool) func(*testing.T) {
	return func(t *testing.T) {
		is := is.New(t)

		s, err := c.createSecret(context.Background(), secret)
		if isErr {
			is.True(err != nil)
			is.Equal(s, nil)
		} else {
			is.NoErr(err)
		}
	}
}

func testUpdateSecret(c *Client, secret *v1.Secret, isErr bool) func(*testing.T) {
	return func(t *testing.T) {
		is := is.New(t)
		if isErr {
			secret.ObjectMeta.Name = "shrug"
		}

		err := c.updateSecret(context.Background(), secret)
		if isErr {
			is.True(err != nil)
		} else {
			is.NoErr(err)
		}
	}
}

func testGetSecret(c *Client, secret *v1.Secret, isErr bool) func(*testing.T) {
	return func(t *testing.T) {
		is := is.New(t)

		s, err := c.getSecret(context.Background(), secret)

		// test result is inverted since we succeed when the secret exists
		if !isErr {
			is.True(err != nil)
			is.Equal(s, nil)
		} else {
			is.NoErr(err)
		}
	}
}
