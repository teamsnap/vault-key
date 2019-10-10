# Non-GCP environment usage

This example demonstrates how to use this repo in golang from a local machine, other cloud provider, or anywhere that's not a Google Cloud Platform environment.

It assumes there's a secret-engine in Vault named `test`, a version 2 key-value secret named `test`, and a key within that secret named `hello`.

It also assumes you've created a Vault role for GCP auth, and that role defines access to the `test` secret.

Make sure to export the following environment variables (refer to the [main README.md](../../../README.md)) for help.

- `GCLOUD_PROJECT`
- `FUNCTION_IDENTITY`
- `VAULT_ADDR`
- `VAULT_ROLE`
- `GOOGLE_APPLICATION_CREDENTIALS`

After they're all exported, just run:

```sh
npm install

node .
```
