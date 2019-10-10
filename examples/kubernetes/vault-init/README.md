# Kubernetes init container

This example will show you how to use this project in an init container. This works by sharing an `emptyDir` volume called `vault-data` between the `vault-init` container and the `cat` container. The init container will pull a secret from Vault and write it to a `.env` file in the shared volume. When the `vault-init` init container finishes, the `cat` container will print the contents of `.env`.

You'll need to replace `base64encodedjsonfile=` with the contents of your service account JSON file. You can get this value by running `cat ./path/to/service-account.json | base64`.

At that point you're ready to `kubectl apply -f deployment.yaml` and then you're good to go.
