require 'vault-key'

secrets = [
  "test/data/test"
]

secretData = Vault.getSecrets(secrets)

puts secretData
