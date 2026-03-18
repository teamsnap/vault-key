# Kubernetes secret sync

This example will show you how to use this project to sync a Vault secret with a Kubernetes secret. This works by pulling a secret from Vault and generating a generic Kubernetes secret that contains all the data that was in the Vault secret.

There are two examples:

1. cronjob
2. job

The cronjob example will sync a Vault secret with Kubernetes every hour.

The job example will only run once, which might be useful as part of a CI/CD pipeline before deploying a new version of an app.

## Environment Variables

| Variable | Required | Default | Description |
|---|---|---|---|
| `VAULT_SECRET` | Yes | | Path to the Vault secret (e.g., `staging/applications/data/myapp/dotenv`) |
| `K8S_NAMESPACE` | Yes | | Kubernetes namespace to create the secret in |
| `K8S_SECRET_NAME` | No | `vault-secret` | Name of the Kubernetes secret to create. Use a custom name to support versioned secrets or avoid overwriting existing secrets. |
| `VAULT_ADDR` | Yes | | Vault server URL |
| `VAULT_ROLE` | Yes | | Vault auth role |
| `GCLOUD_PROJECT` | Yes | | GCP project ID |
| `GCLOUD_REGION` | Yes | | GCP region (e.g., `us-east1`) |
| `CLUSTER_NAME` | Yes | | GKE cluster name |
| `FUNCTION_IDENTITY` | Yes | | GCP service account email |
| `GOOGLE_APPLICATION_CREDENTIALS` | Yes | | Path to the GCP SA key file inside the container |
| `GCP_AUTH_PATH` | No | | Vault GCP auth mount path |
| `VERBOSITY` | No | `info` | Log level: `debug`, `info`, `warn`, `error` |

## Setup

You'll need to replace `base64encodedjsonfile=` with the contents of your service account JSON file. You can get this value by running `cat ./path/to/service-account.json | base64`.

At that point you're ready to `kubectl apply -f job.yaml` or `kubectl apply -f cronjob.yaml` and then you're good to go.
