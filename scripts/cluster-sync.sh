#!/usr/bin/env bash

set -ex

export SCOPE="${SCOPE:-cluster}"
export DEPLOY_NAMESPACE="${DEPLOY_NAMESPACE:-$(oc project --short)}"
export IMAGE_REGISTRY_USER="${DEPLOY_NAMESPACE}"

oc patch configs.imageregistry.operator.openshift.io/cluster --patch '{"spec":{"defaultRoute":true}}' --type=merge

IMAGE_REGISTRY="$(oc get route default-route -n openshift-image-registry --template='{{ .spec.host }}')"
export IMAGE_REGISTRY="$IMAGE_REGISTRY"

# wait for the route
sleep 5

set +x
podman login -u kubeadmin -p "$(oc whoami -t)" --tls-verify=false "$IMAGE_REGISTRY"
set -x

pushd modules
  for MODULE_DIR in $(ls | grep -vE "^(shared|tests)$"); do
    pushd "$MODULE_DIR"
      make release-dev-with-push ARGS="--tls-verify=false"
    popd
  done
popd

export IMAGE_REGISTRY="image-registry.openshift-image-registry.svc:5000"
export TARGET_NAMESPACE="$DEPLOY_NAMESPACE"
oc project $DEPLOY_NAMESPACE

if [[ "$SCOPE" == "cluster" ]]; then
  make deploy-dev
else
  make deploy-dev-namespace
fi
