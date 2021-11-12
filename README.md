# Vault integration with GCP and GitHub for Golang, NodeJS, and Ruby

Vault-Key makes it easy to use Vault with the [Google Cloud Auth Method](https://www.vaultproject.io/docs/auth/gcp.html). It uses a [GCP service account](https://cloud.google.com/iam/docs/service-accounts) and JSON web tokens to log in to Vault securely and without a password. Then it retrieves the secrets you need and makes them available in your code, hassle free.

This project also makes it easy to interact with Vault using [Github Authentication](https://www.vaultproject.io/docs/auth/github).  Interaction is not only limited to retrieval, pkg/vault allows you to manage secrets as well.

## Usage

### Golang

```go
package main

import (
    "context"
    "fmt"
    "github.com/teamsnap/vault-key/pkg/vault"
)

var env = map[string]map[string]string{}

var envArr = []string{
    "secret-engine/data/secret-name",
    "secret-engine-2/data/another-secret-name",
}

func main() {
    ctx := context.Background()

    vault.GetSecrets(ctx, &env, envArr)

    fmt.Println("Secret values:", env)
    fmt.Println("secret-key value = " + env["secret-engine/data/secret-name"]["secret-key"])
    fmt.Println("secret-key-2 value = " + env["secret-engine-2/data/another-secret-name"]["secret-key-2"])
}
```

### NodeJS

```js
const vault = require('@teamsnap/vault-key')

const secrets = [
  'secret-engine/data/secret-name',
  'secret-engine-2/data/another-secret-name'
]

const secretData = vault.getSecrets(secrets)

console.log('Secret values:', JSON.stringify(secretData, null, 4))
console.log(`secret-key value = ${secretData['secret-engine/data/secret-name']['secret-key']}`)
console.log(`secret-key-2 value = ${secretData['secret-engine-2/data/another-secret-name']['secret-key-2']}`)
```

### Ruby

```ruby
require 'vault-key'

secrets = [
  "secret-engine/data/secret-name",
  "secret-engine-2/data/another-secret-name"
]

secretData = Vault.getSecrets(secrets)

puts secretData

puts secretsData["secret-engine/data/secret-name"]["secret-key"]
puts secretsData["secret-engine-2/data/another-secret-name"]["secret-key-2"]
```

### Environment Variable Configuration

|      Environment Variable        |     Default      | Required (GCP) | Required (other environments) |    Example   | Description |
| -------------------------------- | ---------------- | -------------- | ----------------------------- | ---------------------------------------------------- | ---------------------------------------------------------------------------------- |
| `ENVIRONMENT`                    | `"development"`  | No             | No                            | `production`                                         | If set to anything but `production`, prints `trace` level logs                     |
| `FUNCTION_IDENTITY`              | `""`             | No             | Yes                           | `my-project-123@appspot.gserviceaccount.com`         | Email address associated with service account (Required for Google Authentication) |
| `GITHUB_OAUTH_TOKEN`             | `""`             | No             | No                            | `1234abcd`                                           | GitHub Personal Access Token (When set, disables Google Authentication)            |
| `GCLOUD_PROJECT`                 | `""`             | No             | No                            | `my-project-123`                                     | Project ID the service account belongs to                                          |
| `GOOGLE_APPLICATION_CREDENTIALS` | `""`             | No             | No                            | `service-account/my-project-123.serviceaccount.json` | Path to service account credentials file                                           |
| `TRACE_ENABLED`                  | `"false"`        | No             | No                            | `true`                                               | Whether or to enable `opencensus` tracing                                          |
| `TRACE_PREFIX`                   | `"vault"`        | No             | No                            | `my-company`                                         | Prefix added to name of tracing spans                                              |
| `VAULT_ADDR`                     | `""`             | Yes            | Yes                           | `https://vault.my-company.com`                       | Vault address including protocol                                                   |
| `VAULT_ROLE`                     | `""`             | Yes            | No                            | `vault-role-cloud-functions`                         | Name of role created in Vault for GCP auth.  (Required for Google Auth)            |


## GitHub Auth Method

This project also allows you to use GitHub Personal Access tokens for Vault. You'll need to configure a [personal access token](https://docs.github.com/en/free-pro-team@latest/github/authenticating-to-github/creating-a-personal-access-token) for a [user configured with Vault access](https://www.vaultproject.io/api-docs/auth/github). Note that this authentication method is only enabled when the `GITHUB_OAUTH_TOKEN` environment variable is set.  When not set, this project defaults to Google authentication method specified below.

## Google Cloud Auth Method

Because this project uses the [Google Cloud auth method](https://www.vaultproject.io/api/auth/gcp/index.html) for Vault, you'll need to configure a role for the service account you're using. By default, for Google Cloud Functions that will be `<project-id>@appspot.gserviceaccount.com`. You can use the [Terraform example](./examples/terraform/gcp-auth.tf) to get you started.


## Kubernetes

Integrating Vault with Kubernetes is easy to do with this project.

There are examples of two different strategies.

1. [Using an init container](./examples/kubernetes/vault-init) and a shared volume to write a secret to a `.env` file that your app can read in when it's container starts
2. Running a [job](./examples/kubernetes/vault-k8s-secret/job.yaml) or [cronjob](./examples/kubernetes/vault-k8s-secret/cronjob.yaml) to sync Vault secrets with Kubernetes secrets that your deployments can read in like they would any other k8s secrets.

# References

- [Vault Google Cloud Auth](https://www.vaultproject.io/api/auth/gcp/index.html)
- [node-8.16 v8 docs](https://v8docs.nodesource.com/node-8.16/)
- [nodeaddons.com](https://nodeaddons.com/)
