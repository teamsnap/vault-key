const vault = require('./build/Release/vault-auto')

exports.getSecrets = function(secrets) {
  const secretData = JSON.parse(vault.secrets(JSON.stringify(secrets)))

  return secretData
}
