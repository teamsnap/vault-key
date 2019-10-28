require 'vault-key'

secrets = [
  "test/data/test"
]

secretData = VaultKey.getSecrets(secrets)

puts secretData
