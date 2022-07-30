#!/bin/bash

#
# Copyright (c) 2022 Sienna Lloyd
#
# Licensed under the PolyForm Strict License 1.0.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License here:
#  https://github.com/mxplusb/pleiades/blob/mainline/LICENSE
#

export ROOT_DOMAIN="a13s.io"
export PREFIX="a13s-io"

sudo mkdir -p /var/lib/rancher/k3s/server/tls/

k8s_certs=(
  request-header-ca
  client-ca
  server-ca
  client-kube-apiserver
  client-admin client-auth-proxy
  client-controller
  client-k3s-cloud-controller
  client-k3s-controller
  client-kube-apiserver
  client-kube-proxy
  client-scheduler
  serving-kube-apiserver
  serving-kubelet
)

for target in "${k8s_certs}"; do
  ENGINE_PREFIX="${PREFIX}-k3s-${target}"
  vault secrets disable "${ENGINE_PREFIX}"
  vault secrets enable -path="${ENGINE_PREFIX}" pki
  vault secrets tune -max-lease-ttl=43800h "${ENGINE_PREFIX}"
  # generate the intermediate, ensuring to export the private key
  vault write -format=json "${ENGINE_PREFIX}"/intermediate/generate/exported \
    common_name="${ROOT_DOMAIN} k3s ${target} Intermediate Authority G1" >"${target}"-csr.json
  # build the csr
  jq -r '.data.csr' "${target}"-csr.json >"${target}".csr
  # sign the intermediate
  vault write -format=json pki/root/sign-intermediate csr=@"${target}".csr \
    format=pem_bundle ttl="43800h" >"${target}".json
  # write the target files
  jq -r '.data.certificate' "${target}".json >"${target}".crt
  jq -r '.data.private_key' "${target}"-csr.json >"${target}".key
  # verify the modulus
  openssl x509 -noout -modulus -in "${target}".crt | openssl md5
  openssl rsa -noout -modulus -in "${target}".key | openssl md5
  # set the roles
  vault write "${ENGINE_PREFIX}"/roles/"${ROOT_DOMAIN}" \
    allowed_domains="k3s.${ROOT_DOMAIN}" \
    allow_subdomains=true \
    max_ttl="26280h"
  # sign things
  vault write "${ENGINE_PREFIX}"/intermediate/set-signed certificate=@<(cat "${target}".crt)
done

sudo mv ./*.{crt,key} /var/lib/rancher/k3s/server/tls/
rm *.{json,csr}

etcd_certs=(
  client
  peer-ca
  peer-server-client
  server-ca
  server-client
)

for target in "${etcd_certs[@]}"; do
  ENGINE_PREFIX="${PREFIX}-k3s-etcd-${target}"
  vault secrets disable "${ENGINE_PREFIX}"
  vault secrets enable -path="${ENGINE_PREFIX}" pki
  vault secrets tune -max-lease-ttl=43800h "${ENGINE_PREFIX}"
  # generate the intermediate, ensuring to export the private key
  vault write -format=json "${ENGINE_PREFIX}"/intermediate/generate/exported \
    common_name="${ROOT_DOMAIN} k3s ${target} Intermediate Authority G1" >"${target}"-csr.json
  # build the csr
  jq -r '.data.csr' "${target}"-csr.json >"${target}".csr
  # sign the intermediate
  vault write -format=json pki/root/sign-intermediate csr=@"${target}".csr \
    format=pem_bundle ttl="43800h" >"${target}".json
  # write the target files
  jq -r '.data.certificate' "${target}".json >"${target}".crt
  jq -r '.data.private_key' "${target}"-csr.json >"${target}".key
  # verify the modulus
  openssl x509 -noout -modulus -in "${target}".crt | openssl md5
  openssl rsa -noout -modulus -in "${target}".key | openssl md5
  # set the roles
  vault write "${ENGINE_PREFIX}"/roles/"${ROOT_DOMAIN}" \
    allowed_domains="etcd-k3s.${ROOT_DOMAIN}" \
    allow_subdomains=true \
    max_ttl="26280h"
  # sign things
  vault write "${ENGINE_PREFIX}"/intermediate/set-signed certificate=@<(cat "${target}".crt)
done

sudo mv ./*.{crt,key} /var/lib/rancher/k3s/server/tls/etcd/
rm *.{json,csr}
