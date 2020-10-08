#!/usr/bin/env bash

set -ex

export SCOPE="${SCOPE:-cluster}"
export NAMESPACE="${NAMESPACE:-$(oc config current-context | cut -d/ -f1)}"
export TARGET_NAMESPACE="$NAMESPACE"
export IMAGE_REGISTRY_USER="${IMAGE_REGISTRY_USER:-$NAMESPACE}"

oc patch configs.imageregistry.operator.openshift.io/cluster --patch '{"spec":{"defaultRoute":true}}' --type=merge

IMAGE_REGISTRY="$(oc get route default-route -n openshift-image-registry --template='{{ .spec.host }}')"
export IMAGE_REGISTRY="$IMAGE_REGISTRY"

set +x
podman login -u kubeadmin -p "$(oc whoami -t)" --tls-verify=false "$IMAGE_REGISTRY"
set -x

pushd modules
  for MODULE_DIR in $(echo ./* | grep -v "^shared$"); do
    pushd "$MODULE_DIR"
      make release-dev-with-push ARGS="--tls-verify=false"
    popd
  done
popd

export IMAGE_REGISTRY="image-registry.openshift-image-registry.svc:5000"
if [[ "$SCOPE" == "cluster" ]]; then
  make deploy-dev
else
  make deploy-dev-namespace
fi
