#!/usr/bin/env bash

set -e

SCRIPT_DIR="$(dirname "$(readlink -f "$0")")"

source "${SCRIPT_DIR}/common.sh"

PRUNE_IMAGES="${PRUNE_IMAGES:-true}"
FORCE="${FORCE:-true}"
MINIKUBE_CONTAINER_RUNTIME="${MINIKUBE_CONTAINER_RUNTIME:-docker}"
IMAGE_REGISTRY=""

if [[ "${PRUNE_IMAGES}" == "true" ]]; then
  if [[ "$IS_OKD" == "true" ]]; then
    IMAGE_REGISTRY="$(oc get route default-route -n openshift-image-registry --template='{{ .spec.host }}')"
    if [[ "$FORCE" == "true" ]]; then
      oc adm prune images --registry-url=$IMAGE_REGISTRY  --all=false  --keep-younger-than=0s --keep-tag-revisions=0 --confirm
    fi
  elif [[ "${IS_MINIKUBE}" == "true" ]]; then
    IMAGE_REGISTRY="$(minikube ip):5000"
    if [[ "$FORCE" == "true" ]]; then
      IN_CLUSTER_IMAGE_REGISTRY="$(kubectl get service registry -n kube-system --output 'jsonpath={.spec.clusterIP}')"
      if [ -n "${IN_CLUSTER_IMAGE_REGISTRY}" ]; then
        minikube ssh -- "${MINIKUBE_CONTAINER_RUNTIME}"' rmi $('"${MINIKUBE_CONTAINER_RUNTIME}"' images --format "{{.Repository}}" | grep '"${IN_CLUSTER_IMAGE_REGISTRY})" > /dev/null 2>&1 || true
      fi
      if minikube addons list | grep -q "registry .*enabled"; then
        minikube addons disable registry
        minikube addons enable registry
      fi
    fi
  fi
fi

pushd modules
  for TEST_NS in $(kubectl get namespaces -o name | grep -Eo "e2e-tests-[0-9]{5}"); do
    kubectl config set-context --current --namespace="$TEST_NS"
    if [ -n "${IMAGE_REGISTRY}" ]; then
      for MODULE_DIR in $(ls | grep -vE "^(${EXCLUDED_NON_IMAGE_MODULES})$"); do
        podman rmi -f $(podman images --format "{{.Repository}}" | grep "$(basename "${MODULE_DIR}")") || true
      done
    fi

    for RESOURCE in pipelineresources.tekton.dev \
      pipelineruns.tekton.dev \
      pipelines.tekton.dev \
      taskruns.tekton.dev \
      virtualmachines.kubevirt.io; do
        kubectl delete $RESOURCE --all
    done

    if [[ "$IS_OKD" == "true" ]]; then
      kubectl delete templates.template.openshift.io --all
    fi

    kubectl delete namespace "$TEST_NS"
  done
popd

kubectl config set-context --current --namespace=default > /dev/null 2>&1
