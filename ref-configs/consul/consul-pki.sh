#!/bin/bash

vault secrets enable -path=consul-pki pki

vault secrets tune -max-lease-ttl=87600h consul-pki

vault write -field=certificate consul-pki/root/generate/internal \
        common_name="" \
        ttl=87600h > consul_ca.crt

vault write consul-pki/config/urls \
        issuing_certificates="http://127.0.0.1:8200/v1/pki/ca" \
        crl_distribution_points="http://127.0.0.1:8200/v1/pki/crl"

vault secrets enable -path=consul-pki-int pki

vault secrets tune -max-lease-ttl=43800h consul-pki-int

vault write -format=json consul-pki-int/intermediate/generate/internal \
        common_name="dc1.consul Intermediate Authority G1" \
        | jq -r '.data.csr' > consul_pki_intermediate.csr

vault write -format=json consul-pki/root/sign-intermediate csr=@consul_pki_intermediate.csr \
        format=pem_bundle ttl="43800h" \
        | jq -r '.data.certificate' > consul_intermediate_cert.pem

vault write consul-pki-int/intermediate/set-signed certificate=@consul_intermediate_cert.pem

vault write consul-pki-int/roles/consul-dc1 \
  allowed_domains="dc1.consul" \
  allow_subdomains=true \
  generate_lease=true \
  max_ttl="720h"

vault secrets enable -path=connect-pki pki
vault secrets tune -max-lease-ttl=87600h connect-pki
vault secrets enable -path=connect-pki-int pki
vault secrets tune -max-lease-ttl=43800h connect-pki-int

