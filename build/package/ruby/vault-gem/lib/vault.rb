require 'ffi'
require 'json'
require 'OS'

module VaultNative
  extend FFI::Library

  sharedObjectPath = OS.linux? ? "./native/vault.linux.so" : "./native/vault.darwin.so"

  ffi_lib File.expand_path(sharedObjectPath, File.dirname(__FILE__))
  attach_function :GetSecrets, [:string], :string
end

class Vault
  def self.getSecrets(secrets)
    secretsData = VaultNative.GetSecrets(JSON.generate(secrets))
    JSON.parse(secretsData)
  end
end
