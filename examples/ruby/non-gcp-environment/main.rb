require 'vault'

secrets = [
  "test/data/test"
]

secretData = Vault.getSecrets(secrets)

puts secretData
