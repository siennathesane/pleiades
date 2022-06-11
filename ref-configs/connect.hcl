# Consul Managed PKI Mounts
path "/sys/mounts" {
  capabilities = [ "read" ]
}

path "/sys/mounts/connect-pki" {
  capabilities = [ "create", "read", "update", "delete", "list" ]
}

path "/sys/mounts/connect-pki-int" {
  capabilities = [ "create", "read", "update", "delete", "list" ]
}

path "/connect-pki/*" {
  capabilities = [ "create", "read", "update", "delete", "list" ]
}

path "/connect-pki-int/*" {
  capabilities = [ "create", "read", "update", "delete", "list" ]
}

# Work with consul secrets engine
path "consul*" {
  capabilities = [ "create", "read", "update", "delete", "list", "sudo" ]
}

path "identity/oidc/*" {
  capabilities = ["read"]
}
