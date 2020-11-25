#!/usr/bin/env bash

set -e

SCRIPT_DIR="$(dirname "$(readlink -f "$0")")"
REPO_DIR="$(realpath "${SCRIPT_DIR}/..")"

source "${SCRIPT_DIR}/common.sh"

export SCOPE="${SCOPE:-cluster}"
export DEPLOY_NAMESPACE="${DEPLOY_NAMESPACE:-$(oc project --short)}"

oc patch configs.imageregistry.operator.openshift.io/cluster --patch '{"spec":{"defaultRoute":true}}' --type=merge
IMAGE_REGISTRY="$(oc get route default-route -n openshift-image-registry --template='{{ .spec.host }}')"

# wait for the route
sleep 5

podman login -u kubeadmin -p "$(oc whoami -t)" --tls-verify=false "$IMAGE_REGISTRY"

visit "${REPO_DIR}"
  visit modules
    for TASK_NAME in $(ls | grep -vE "^(${EXCLUDED_NON_IMAGE_MODULES})$"); do
      visit "${TASK_NAME}"
        IMAGE_NAME_AND_TAG="tekton-task-${TASK_NAME}:latest"
        export IMAGE="${IMAGE_REGISTRY}/${DEPLOY_NAMESPACE}/${IMAGE_NAME_AND_TAG}"
        podman build -f "build/${TASK_NAME}/Dockerfile" -t "${IMAGE}" .
        podman push "${IMAGE}" --tls-verify=false

        # set inside-cluster registry
        export IMAGE="image-registry.openshift-image-registry.svc:5000/${DEPLOY_NAMESPACE}/${IMAGE_NAME_AND_TAG}"
        export ${TASK_NAME_TO_ENV_NAME[${TASK_NAME}]}="${IMAGE}"
      leave
    done
  leave
leave

"${REPO_DIR}/scripts/deploy-tasks.sh"
