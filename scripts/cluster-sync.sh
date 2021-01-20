#!/usr/bin/env bash

set -e

SCRIPT_DIR="$(dirname "$(readlink -f "$0")")"
REPO_DIR="$(realpath "${SCRIPT_DIR}/..")"

source "${SCRIPT_DIR}/common.sh"

export SCOPE="${SCOPE:-cluster}"
export DEPLOY_NAMESPACE="${DEPLOY_NAMESPACE:-$(kubectl config view --minify --output 'jsonpath={..namespace}')}"
IMAGE_REGISTRY=""
IN_CLUSTER_IMAGE_REGISTRY=""

if [[ "${IS_OPENSHIFT}" == "true" ]]; then
  oc patch configs.imageregistry.operator.openshift.io/cluster --patch '{"spec":{"defaultRoute":true}}' --type=merge
  IMAGE_REGISTRY="$(oc get route default-route -n openshift-image-registry --template='{{ .spec.host }}')"
  IN_CLUSTER_IMAGE_REGISTRY="image-registry.openshift-image-registry.svc:5000"
  # wait for the route
  sleep 5

  podman login -u kubeadmin -p "$(oc whoami -t)" --tls-verify=false "$IMAGE_REGISTRY"
elif [[ "${IS_MINIKUBE}" == "true" ]]; then
  if ! minikube addons list | grep -q "registry .*enabled"; then
     echo "minikube should have registry addon enabled" >&2
     exit 2
  fi
  IMAGE_REGISTRY="$(minikube ip):5000"
  IN_CLUSTER_IMAGE_REGISTRY="$(kubectl get service registry -n kube-system --output 'jsonpath={.spec.clusterIP}')"
else
  echo "only minikube or openshift is supported" >&2
  exit 3
fi


visit "${REPO_DIR}"
  visit modules
    if [[ $# -eq 0 ]]; then
      TASK_NAMES=(*)
    else
      TASK_NAMES=("$@")
    fi
    for TASK_NAME in ${TASK_NAMES[*]}; do
      if echo "${TASK_NAME}" | grep -vqE "^(${EXCLUDED_NON_IMAGE_MODULES})$"; then
        if [ ! -d  "${TASK_NAME}" ]; then
          continue
        fi
        visit "${TASK_NAME}"
          IMAGE_NAME_AND_TAG="tekton-task-${TASK_NAME}:latest"
          export IMAGE="${IMAGE_REGISTRY}/${DEPLOY_NAMESPACE}/${IMAGE_NAME_AND_TAG}"
          podman build -f "build/${TASK_NAME}/Dockerfile" -t "${IMAGE}" .
          podman push "${IMAGE}" --tls-verify=false

          # set inside-cluster registry
          export IMAGE="${IN_CLUSTER_IMAGE_REGISTRY}/${DEPLOY_NAMESPACE}/${IMAGE_NAME_AND_TAG}"
          export ${IMAGE_MODULE_NAME_TO_ENV_NAME[${TASK_NAME}]}="${IMAGE}"
        leave
      fi
    done
  leave
leave

"${REPO_DIR}/scripts/deploy-tasks.sh" "$@"
