package k8s

import (
	"context"
	"fmt"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type client struct {
	clientset kubernetes.Interface
}

func (c client) createSecret(ctx context.Context, secret *apiv1.Secret) (*apiv1.Secret, error) {
	fmt.Println("creating secret...")
	cs, err := c.clientset.CoreV1().Secrets(secret.Namespace).Create(ctx, secret, metav1.CreateOptions{})
	if err != nil {
		return nil, fmt.Errorf("secret probably already exists %s \n", err)
	}

	return cs, nil
}

func (c client) getSecret(ctx context.Context, secret *apiv1.Secret) (*apiv1.Secret, error) {
	fmt.Println("getting secret...")
	gs, err := c.clientset.CoreV1().Secrets(secret.Namespace).Get(ctx, "vault-secret", metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get latest version of secret: %s", err)
	}

	gs.Data = secret.Data
	return gs, nil
}

func (c client) updateSecret(ctx context.Context, secret *apiv1.Secret) error {
	fmt.Println("updating secret...")
	_, err := c.clientset.CoreV1().Secrets(secret.Namespace).Update(ctx, secret, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("failed to udpate latest version of secret: %s", err)
	}

	return err
}
