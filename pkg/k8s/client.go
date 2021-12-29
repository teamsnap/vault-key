package k8s

import (
	"context"
	"fmt"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/util/retry"
)

type Client struct {
	Clientset kubernetes.Interface
}

func (c Client) ApplySecret(ctx context.Context, secret *apiv1.Secret) error {
	if _, err := c.createSecret(ctx, secret); err != nil {
		err := retry.RetryOnConflict(retry.DefaultRetry, func() error {
			// Fetch the resource here; you need to refetch it on every try, since
			// if you got a conflict on the last update attempt then you need to get
			// the current version before making your own changes.

			// RetryOnConflict uses exponential backoff to avoid exhausting the apiserver
			gs, err := c.getSecret(ctx, secret)
			if err != nil {
				return fmt.Errorf("retrieving the latest version %w", err)
			}

			gs.Data = secret.Data
			if err := c.updateSecret(ctx, gs); err != nil {
				return fmt.Errorf("update secret %w", err)
			}

			return nil
		})

		if err != nil {
			// May be conflict if max retries were hit, or may be something unrelated
			// like permissions or a network error

			return fmt.Errorf("", err)
		}
	}

	return nil
}

func (c Client) createSecret(ctx context.Context, secret *apiv1.Secret) (*apiv1.Secret, error) {
	if _, err := c.Clientset.CoreV1().Secrets(secret.Namespace).Create(ctx, secret, metav1.CreateOptions{}); err != nil {
		return nil, fmt.Errorf("secret probably already exists %s", err)
	}

	return secret, nil
}

func (c Client) getSecret(ctx context.Context, secret *apiv1.Secret) (*apiv1.Secret, error) {
	gs, err := c.Clientset.CoreV1().Secrets(secret.Namespace).Get(ctx, "vault-secret", metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get latest version of secret: %s", err)
	}

	gs.Data = secret.Data
	return gs, nil
}

func (c Client) updateSecret(ctx context.Context, secret *apiv1.Secret) error {
	if _, err := c.Clientset.CoreV1().Secrets(secret.Namespace).Update(ctx, secret, metav1.UpdateOptions{}); err != nil {
		return fmt.Errorf("failed to udpate latest version of secret: %s", err)
	}

	return nil
}
