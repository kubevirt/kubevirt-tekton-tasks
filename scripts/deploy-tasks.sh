#!/usr/bin/env bash

set -e

SCRIPT_DIR="$(dirname "$(readlink -f "$0")")"
REPO_DIR="$(realpath "${SCRIPT_DIR}/..")"

source "${SCRIPT_DIR}/common.sh"


# run only for specified tasks in script arguments
# or default to all if no arguments specified

SCOPE="${SCOPE:-cluster}"
DEPLOY_NAMESPACE="${DEPLOY_NAMESPACE:-$(oc project --short)}"

oc project "${DEPLOY_NAMESPACE}" || true

visit "${REPO_DIR}/tasks"
  for TASK_NAME in *; do
    CONFIG_FILE="${REPO_DIR}/configs/${TASK_NAME}.yaml"
    MAIN_IMAGE="$(sed -n  's/^main_image *: *//p' "${CONFIG_FILE}")"
    CUSTOM_IMAGE="${TASK_NAME_TO_IMAGE[${TASK_NAME}]}"

    visit "${TASK_NAME}"
      oc delete -f manifests 2> /dev/null || true
      if [[ $SCOPE == "cluster" ]]; then
        sed "s/TARGET_NAMESPACE/${DEPLOY_NAMESPACE}/" "manifests/${TASK_NAME}-cluster-rbac.yaml" | oc apply -f -
      else
        oc apply -f "manifests/${TASK_NAME}-namespace-rbac.yaml"
      fi

      for SUBTASK_NAME in $(ls manifests | grep -v rbac); do
        if [[ -z ${CUSTOM_IMAGE} ]]; then
          oc apply -f "manifests/${SUBTASK_NAME}"
        else
          sed "s!${MAIN_IMAGE}!${CUSTOM_IMAGE}!g" "manifests/${SUBTASK_NAME}" | oc apply -f -
        fi
      done
    leave
  done
leave
