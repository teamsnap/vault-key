package k8s

import (
	"context"
	"fmt"
	"log"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/util/retry"
)

type Client struct {
	Clientset kubernetes.Interface
}

func (c Client) ApplySecret(ctx context.Context, secret *apiv1.Secret) error {
	_, err := c.createSecret(ctx, secret)
	if err != nil {
		log.Print(err)
		err := retry.RetryOnConflict(retry.DefaultRetry, func() error {
			// Retrieve the latest version of Secret before attempting update
			// RetryOnConflict uses exponential backoff to avoid exhausting the apiserver
			gs, err := c.getSecret(ctx, secret)
			if err != nil {
				return err
			}

			err = c.updateSecret(ctx, gs)
			if err != nil {
				return err
			}

			return err
		})

		if err != nil {
			return err
		}
	}

	log.Print("applied secret...")
	return nil
}

func (c Client) createSecret(ctx context.Context, secret *apiv1.Secret) (*apiv1.Secret, error) {
	log.Print("creating secret...")
	_, err := c.Clientset.CoreV1().Secrets(secret.Namespace).Create(ctx, secret, metav1.CreateOptions{})
	if err != nil {
		return nil, fmt.Errorf("secret probably already exists %s \n", err)
	}

	return secret, nil
}

func (c Client) getSecret(ctx context.Context, secret *apiv1.Secret) (*apiv1.Secret, error) {
	log.Print("getting secret...")
	gs, err := c.Clientset.CoreV1().Secrets(secret.Namespace).Get(ctx, "vault-secret", metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get latest version of secret: %s", err)
	}

	gs.Data = secret.Data
	return gs, nil
}

func (c Client) updateSecret(ctx context.Context, secret *apiv1.Secret) error {
	log.Print("updating secret...")
	_, err := c.Clientset.CoreV1().Secrets(secret.Namespace).Update(ctx, secret, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("failed to udpate latest version of secret: %s", err)
	}

	return err
}
