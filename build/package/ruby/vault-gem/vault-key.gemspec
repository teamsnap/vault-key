Gem::Specification.new do |s|
  s.name        = 'vault-key'
  s.version     = '0.0.1'
  s.date        = '2020-12-17'
  s.summary     = "Vault with GCP auth"
  s.description = "A gem to integrate Vault with GCP auth"
  s.authors     = ["TeamSnap SREs"]
  s.email       = ''
  s.files       = [
    "lib/vault.rb",
    "lib/native/vault.darwin.h",
    "lib/native/vault.darwin.so",
    "lib/native/vault.linux.h",
    "lib/native/vault.linux.so"
  ]
  s.homepage    = ''
  s.license       = 'MIT'
  s.add_dependency('ffi', '1.11.1')
  s.add_dependency('os', '1.0.1')
end
