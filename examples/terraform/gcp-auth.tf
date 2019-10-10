provider "vault" { }

resource "vault_gcp_auth_backend" "gcp" {
   credentials  = "${file("../../secrets/my-project-123.service-account-key.json")}"
}

resource "vault_policy" "test-policy" {
  name = "test-policy"

  policy = <<EOF
path "test/*" {
  capabilities = ["create", "read", "update", "delete", "list"]
}
EOF
}

resource "vault_gcp_auth_backend_role" "gcp" {
    role = "vault-role-cloud-functions"
    type = "iam"
    bound_projects         = ["my-project-123"]
    backend                = "${vault_gcp_auth_backend.gcp.path}"
    bound_service_accounts = ["my-project-123@appspot.gserviceaccount.com"]
    token_policies         = ["test-policy"]
}
