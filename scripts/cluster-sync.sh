#!/usr/bin/env bash

set -e

SCRIPT_DIR="$(dirname "$(readlink -f "$0")")"
REPO_DIR="$(realpath "${SCRIPT_DIR}/..")"
source "./scripts/common.sh"

wait_for_image_registry_default_route() {
  while true; do
    if ! host="$(oc get route default-route -n openshift-image-registry --ignore-not-found \
      -o go-template='{{range .status.ingress}}{{range .conditions}}{{if and (eq .type "Admitted") (eq .status "True")}}{{$.spec.host}}{{end}}{{end}}{{end}}')"; then
      echo "failed to query route default-route (permanent error)" >&2
      exit 1
    fi
    [[ -n "${host}" ]] && break
    echo "waiting for image-registry default-route to be admitted..."
    sleep 5
  done
}

resolve_deploy_namespace

IMAGE_REGISTRY=""
IN_CLUSTER_IMAGE_REGISTRY=""

if [[ "${IS_OKD}" == "true" ]]; then
  oc patch configs.imageregistry.operator.openshift.io/cluster --patch '{"spec":{"defaultRoute":true}}' --type=merge
  export -f wait_for_image_registry_default_route
  timeout 10m bash -c 'wait_for_image_registry_default_route'
  IMAGE_REGISTRY="$(oc get route default-route -n openshift-image-registry --template='{{ .spec.host }}')"
  IN_CLUSTER_IMAGE_REGISTRY="image-registry.openshift-image-registry.svc:5000"

  podman login -u kubeadmin -p "$(oc whoami -t)" --tls-verify=false "$IMAGE_REGISTRY"
elif [[ "${IS_MINIKUBE}" == "true" ]]; then
  if ! minikube addons list | grep -q "registry .*enabled"; then
     echo "minikube should have registry addon enabled" >&2
     exit 2
  fi
  IMAGE_REGISTRY="$(minikube ip):5000"
  IN_CLUSTER_IMAGE_REGISTRY="$(kubectl get service registry -n kube-system --output 'jsonpath={.spec.clusterIP}')"
else
  echo "only minikube or OKD is supported" >&2
  exit 3
fi

IMAGE_NAME_AND_TAG="tekton-tasks:latest"
export IMAGE="${IMAGE_REGISTRY}/${DEPLOY_NAMESPACE}/${IMAGE_NAME_AND_TAG}"
podman build -f "build/Containerfile" -t "${IMAGE}" .
podman push "${IMAGE}" --tls-verify=false

# set inside-cluster registry
export IMAGE="${IN_CLUSTER_IMAGE_REGISTRY}/${DEPLOY_NAMESPACE}/${IMAGE_NAME_AND_TAG}"
export TEKTON_TASKS_IMAGE="${IMAGE}"

IMAGE_NAME_AND_TAG="tekton-tasks-disk-virt:latest"
export IMAGE="${IMAGE_REGISTRY}/${DEPLOY_NAMESPACE}/${IMAGE_NAME_AND_TAG}"
podman build -f "build/Containerfile.DiskVirt" -t "${IMAGE}" .
podman push "${IMAGE}" --tls-verify=false

# set inside-cluster registry
export IMAGE="${IN_CLUSTER_IMAGE_REGISTRY}/${DEPLOY_NAMESPACE}/${IMAGE_NAME_AND_TAG}"
export TEKTON_TASKS_DISK_VIRT_IMAGE="${IMAGE}"

"${REPO_DIR}/scripts/deploy-tasks.sh" "$@"
"${REPO_DIR}/scripts/deploy-pipelines.sh"
