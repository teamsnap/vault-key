package k8s

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"k8s.io/client-go/util/retry"

	// enable gcp auth for k8s client
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

// Secret holds a map of the secrets that need to be created in Kubernetes
type Secret struct {
	Secrets   map[string]string
	Namespace string
}

// ApplySecret takes a Vault secret and k8s namespace and creates the k8s secret based on the data
func ApplySecret(vaultSecret *Secret) {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}

	flag.Parse()
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	secretData := map[string][]byte{}
	for key, val := range vaultSecret.Secrets {
		secretData[key] = []byte(val)
	}

	secret := &apiv1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: "vault-secret",
		},
		Data: secretData,
	}

	client := client{
		clientset: clientset,
	}

	ctx := context.Background()
	cs, err := client.createSecret(ctx, secret)
	if err != nil {
		fmt.Println(err)
		err = retry.RetryOnConflict(retry.DefaultRetry, func() error {
			// Retrieve the latest version of Secret before attempting update
			// RetryOnConflict uses exponential backoff to avoid exhausting the apiserver
			gs, err := client.getSecret(ctx, cs)
			if err != nil {
				panic(err)
			}

			err = client.updateSecret(ctx, gs)
			if err != nil {
				panic(err)
			}

			return err
		})

		if err != nil {
			fmt.Printf("apply secret failed: %s", err)
		}
	}

	fmt.Println("applied secret...")
}
