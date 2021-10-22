package k8s

import (
	"context"
	"flag"
	"log"
	"path/filepath"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"

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

	client := &Client{
		Clientset: clientset,
	}

	ctx := context.Background()
	client.ApplySecret(ctx, secret)
	if err != nil {
		log.Printf("apply secret failed: %s", err)
	}
}
