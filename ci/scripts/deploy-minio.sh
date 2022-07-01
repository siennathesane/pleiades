#!/bin/bash

#
# Copyright (c) 2022 Sienna Lloyd
#
# Licensed under the PolyForm Strict License 1.0.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License here:
#  https://github.com/mxplusb/pleiades/blob/mainline/LICENSE
#

set -eux

echo "preparing kubeconfig"
mkdir -p "${HOME}"/.kube
echo "${KUBE_CONFIG}" | base64 -d > "${HOME}"/.kube/config
chmod go-r "${HOME}"/.kube/config

echo "adding minio repo"
helm repo add minio https://charts.min.io/
helm repo update

exists="$(helm ls -A -o json | jq -r '.[].name' | grep ${CHART_NAME} || true)"

COMMAND="upgrade"
export COMMAND

if [ -z "${exists}" ]; then
  COMMAND="install"
fi

echo "applying minio helm chart"

helm "${COMMAND}" \
  "${CHART_NAME}" \
  --atomic \
  --set rootUser="${ROOT_USER}" \
  --set rootPassword="${ROOT_PASSWORD}" \
  --set mode=standalone \
  --set replicas=1 \
  --set resources.requests.memory="512Mi" \
  --set buckets[0].name="${BINARIES_BUCKET}",buckets[0].policy=none,buckets[0].purge=false \
  --create-namespace \
  --namespace "${NAMESPACE}" \
  minio/minio

kubectl apply --namespace "${NAMESPACE}" -f ref-configs/minio/ingress.yaml
