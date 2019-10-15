# Non-GCP environment usage

This example demonstrates how to use this repo in ruby from a local machine, other cloud provider, or anywhere that's not a Google Cloud Platform environment.

It assumes there's a secret-engine in Vault named `test`, a version 2 key-value secret named `test`, and a key within that secret named `hello`.

It also assumes you've created a Vault role for GCP auth, and that role defines access to the `test` secret.

Make sure to export the following environment variables (refer to the [main README.md](../../../README.md)) for help.

- `GCLOUD_PROJECT`
- `FUNCTION_IDENTITY`
- `VAULT_ADDR`
- `VAULT_ROLE`
- `GOOGLE_APPLICATION_CREDENTIALS`

This requires you to build the gem before being able to use it.

1. move to the [build/package/ruby](../../../build/package/ruby) directory
2. `make build_all`

After all necessary environment variables are exported, and the gem is built run, move back into `examples/ruby/non-gcp-environment` and run:

```sh
gem install ../../../build/package/ruby/vault-gem/vault-0.0.0.gem

ruby main.rb
```
