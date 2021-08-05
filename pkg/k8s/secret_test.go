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
		client *client
		isErr  bool
	}{
		{
			name:   "failure",
			secret: secret,
			client: &client{clientset: testclient.NewSimpleClientset(secret)},
			isErr:  true,
		},
		{
			name:   "sucess",
			secret: secret,
			client: &client{clientset: testclient.NewSimpleClientset()},
			isErr:  false,
		},
	}

	for _, c := range cases {
		t.Run("do "+c.name, testDo(c.client, c.name))
		t.Run("create secret "+c.name, testCreateSecret(c.client, c.secret, c.isErr))
		t.Run("update secret "+c.name, testUpdateSecret(c.client, c.secret))
		t.Run("get secret "+c.name, testGetSecret(c.client, c.secret, c.name))
	}
}

func testDo(c *client, name string) func(*testing.T) {
	return func(t *testing.T) {
		is := is.New(t)

		secret := &v1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name: "vault-do-secret",
			},
			Data: map[string][]byte{"key": []byte("value")},
		}

		if name == "failure" {
			c.createSecret(context.Background(), secret)
		}

		err := c.do(context.Background(), secret)
		is.NoErr(err)
	}
}

func testCreateSecret(c *client, secret *v1.Secret, isErr bool) func(*testing.T) {
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

func testUpdateSecret(c *client, secret *v1.Secret) func(*testing.T) {
	return func(t *testing.T) {
		is := is.New(t)
		err := c.updateSecret(context.Background(), secret)
		is.NoErr(err)
	}
}

func testGetSecret(c *client, secret *v1.Secret, name string) func(*testing.T) {
	return func(t *testing.T) {
		is := is.New(t)
		s, err := c.getSecret(context.Background(), secret)

		switch name {
		case "success":
			is.NoErr(err)
		case "fail":
			is.True(err != nil)
			is.Equal(s, nil)
		}
	}
}
