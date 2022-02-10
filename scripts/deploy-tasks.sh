#!/usr/bin/env bash

set -e

SCRIPT_DIR="$(dirname "$(readlink -f "$0")")"
REPO_DIR="$(realpath "${SCRIPT_DIR}/..")"

source "${SCRIPT_DIR}/common.sh"


# run only for specified tasks in script arguments
# or default to all if no arguments specified

SCOPE="${SCOPE:-cluster}"
DEPLOY_NAMESPACE="${DEPLOY_NAMESPACE:-$(kubectl config view --minify --output 'jsonpath={..namespace}')}"

kubectl config set-context --current --namespace="${DEPLOY_NAMESPACE}" || true

visit "${REPO_DIR}/tasks"
  if [[ $# -eq 0 ]]; then
    TASK_NAMES=(*)
  else
    TASK_NAMES=("$@")
  fi
  for TASK_NAME in ${TASK_NAMES[*]}; do
    if [ ! -d  "${TASK_NAME}" ]; then
      continue
    fi
    CONFIG_FILE="${REPO_DIR}/configs/${TASK_NAME}.yaml"
    MAIN_IMAGE="$(sed -n  's/^main_image *: *//p' "${CONFIG_FILE}"):${VERSION}"
    CUSTOM_IMAGE="${TASK_NAME_TO_IMAGE[${TASK_NAME}]}"

    # cleanup first
    kubectl delete -f manifests 2> /dev/null || true
    kubectl delete clusterrolebinding "${TASK_NAME}-task" 2> /dev/null || true

    visit "${TASK_NAME}"

      if [[ -z ${CUSTOM_IMAGE} ]]; then
        kubectl apply -f "manifests/${TASK_NAME}.yaml"
      else
        sed "s!${MAIN_IMAGE}!${CUSTOM_IMAGE}!g" "manifests/${TASK_NAME}.yaml" | kubectl apply -f -
      fi

      # add cluster privileges if needed
      if [[ "${SCOPE}" == "cluster" ]] && grep -q ClusterRole "manifests/${TASK_NAME}.yaml"; then
        kubectl apply -f - << EOF
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: ${TASK_NAME}-task
roleRef:
  kind: ClusterRole
  name: ${TASK_NAME}-task
  apiGroup: rbac.authorization.k8s.io
subjects:
  - kind: ServiceAccount
    name:  ${TASK_NAME}-task
    namespace: ${DEPLOY_NAMESPACE}
EOF
      fi
    leave
  done
leave
