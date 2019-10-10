const vault = require('@teamsnap/vault')

const secrets = [
	'test/data/test'
]

const secretData = vault.getSecrets(secrets)

exports.vaultNodeJS = (pubSubEvent, context) => {
  return new Promise((resolve, reject) => {
    console.log('Secret values:', JSON.stringify(secretData, null, 4))
		resolve()
  })
}
