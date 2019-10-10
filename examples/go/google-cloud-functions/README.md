# Deploy to Cloud Functions

This example demonstrates how to use this repo in golang on Google Cloud Functions.

It assumes there's a secret-engine in Vault named `test`, a version 2 key-value secret named `test`, and a key within that secret named `hello`.

It also assumes you've created a Vault role for GCP auth, and that role defines access to the `test` secret.

```sh
export GCF_FUNCTION_NAME="VaultOnInit"
export VAULT_ROLE="vault-role-cloud-functions"
export VAULT_ADDR="https://vault.your-domain.com"

gcloud functions deploy "$GCF_FUNCTION_NAME" \
  --runtime go111 \
  --trigger-topic vault_go \
  --region us-central1 \
  --set-env-vars "VAULT_ROLE=$VAULT_ROLE,VAULT_ADDR=$VAULT_ADDR"
```
