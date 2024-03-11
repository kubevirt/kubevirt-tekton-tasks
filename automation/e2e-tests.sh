#!/usr/bin/env bash

set -ex

export DEPLOY_NAMESPACE="${DEPLOY_NAMESPACE:-e2e-tests-$(shuf -i10000-99999 -n1)}"
source "./automation/e2e-source.sh"

./automation/e2e-deploy-resources.sh

kubectl get namespaces -o name | grep -Eq "^namespace/$DEPLOY_NAMESPACE$" || kubectl create namespace "$DEPLOY_NAMESPACE"
kubectl config set-context --current --namespace="${DEPLOY_NAMESPACE}"

if [[ "$DEV_MODE" == "true" ]]; then
  make cluster-sync
else
  make deploy
fi

make cluster-test
