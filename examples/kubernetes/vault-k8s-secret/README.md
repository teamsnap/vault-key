# Kubernetes secret sync

This example will show you how to use this project to sync a Vault secret with a Kubernetes secret. This works by pulling a secret from Vault and generating a generic Kubernetes secret that contains all the data that was in the Vault secret.

There are two examples:

1. cronjob
2. job

The cronjob example will sync a Vault secret with Kubernetes every hour.

The job example will only run once, which might be useful as part of a CI/CD pipeline before deploying a new version of an app.

You'll need to replace `base64encodedjsonfile=` with the contents of your service account JSON file. You can get this value by running `cat ./path/to/service-account.json | base64`.

At that point you're ready to `kubectl apply -f job.yaml` or `kubectl apply -f cronjob.yaml` and then you're good to go.
