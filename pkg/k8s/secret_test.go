package k8s_test

import (
	"context"
	"testing"

	"github.com/matryer/is"
	"github.com/teamsnap/vault-key/pkg/k8s"
	apiv1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	testclient "k8s.io/client-go/kubernetes/fake"
)

func TestClient(t *testing.T) {
	secret := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: "secret",
		},
		Data: map[string][]byte{"key": []byte("value")},
	}

	cases := []struct {
		name   string
		secret *v1.Secret
		client *k8s.Client
		isErr  bool
	}{
		{
			name:   "failure",
			secret: secret,
			client: &k8s.Client{Clientset: testclient.NewSimpleClientset(secret)},
			isErr:  true,
		},
		{
			name:   "sucess",
			secret: secret,
			client: &k8s.Client{Clientset: testclient.NewSimpleClientset()},
			isErr:  false,
		},
	}

	for _, c := range cases {
		t.Run("apply secret "+c.name, testApply(c.client, c.secret, c.isErr))
	}
}

func testApply(c *k8s.Client, secret *apiv1.Secret, isErr bool) func(*testing.T) {
	return func(t *testing.T) {
		is := is.New(t)

		err := c.ApplySecret(context.Background(), secret)
		if isErr {
			is.True(err != nil)
		} else {
			is.NoErr(err)
		}
	}
}
