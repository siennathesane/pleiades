#!/bin/bash

#
# Copyright (c) 2022 Sienna Lloyd
#
# Licensed under the PolyForm Strict License 1.0.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License here:
#  https://github.com/mxplusb/pleiades/blob/mainline/LICENSE
#

set -e

echo "preparing kubeconfig"
mkdir -p "${HOME}"/.kube
echo "${KUBE_CONFIG}" | base64 -d > "${HOME}"/.kube/config
chmod go-r "${HOME}"/.kube/config

echo "adding minio repo"
helm repo add minio https://charts.min.io/
helm repo update

exists=$(helm ls -A -o json | jq -r '.[].name' | grep "${CHART_NAME}")

COMMAND="upgrade"
export COMMAND

if [ -z "${exists}" ]; then
  COMMAND="install"
fi

echo "applying minio helm chart"

set -x

helm "${COMMAND}" \
  "${CHART_NAME}" \
  --atomic \
  --set rootUser="${ROOT_USER}" \
  --set rootPassword="${ROOT_PASSWORD}" \
  --set mode=standalone \
  --set replicas=1 \
  --create-namespace \
  --namespace "${NAMESPACE}" \
  minio/minio

set +x

echo "creating ingress"

kubectl apply --namespace "${NAMESPACE}" -f ref-configs/minio/ingress.yaml
