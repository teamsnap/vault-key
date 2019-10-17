const vault = require('@teamsnap/vault-key')

const secrets = [
  'test/data/test'
]

const secretData = vault.getSecrets(secrets)

console.log('Secret map:', JSON.stringify(secretData, null, 4))
