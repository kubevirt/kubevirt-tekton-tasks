#!/usr/bin/env bash

set -e

SCRIPT_DIR="$(dirname "$(readlink -f "$0")")"
REPO_DIR="$(realpath "${SCRIPT_DIR}/..")"

source "${SCRIPT_DIR}/common.sh"


# run only for specified tasks in script arguments
# or default to all if no arguments specified

SCOPE="${SCOPE:-cluster}"
DEPLOY_NAMESPACE="${DEPLOY_NAMESPACE:-$(oc project --short)}"

CREATE_VM_IMAGE="${CREATE_VM_IMAGE:-}"

declare -A CUSTOM_IMAGES
CUSTOM_IMAGES["create-vm"]="${CREATE_VM_IMAGE}"

oc project "${DEPLOY_NAMESPACE}" || true

visit "${REPO_DIR}/tasks"
  for TASK_NAME in *; do
    CONFIG_FILE="${REPO_DIR}/configs/${TASK_NAME}.yaml"
    MAIN_IMAGE="$(sed -n  's/^main_image *: *//p' "${CONFIG_FILE}")"
    SUBTASK_NAMES=( "$(sed -n -e  '/^subtask_names *: */,/^ *^[-]/p'  "${CONFIG_FILE}" | sed -n  's/^ *-//p')" )
    CUSTOM_IMAGE="${CUSTOM_IMAGES[${TASK_NAME}]}"

    visit "${TASK_NAME}"
      oc delete -f manifests 2> /dev/null || true
      if [[ $SCOPE == "cluster" ]]; then
        sed "s/TARGET_NAMESPACE/${DEPLOY_NAMESPACE}/" "manifests/${TASK_NAME}-cluster-rbac.yaml" | oc apply -f -
      else
        oc apply -f "manifests/${TASK_NAME}-namespace-rbac.yaml"
      fi

      for SUBTASK_NAME in ${SUBTASK_NAMES[*]}; do
        if [[ -z ${CUSTOM_IMAGE} ]]; then
          oc apply -f "manifests/${SUBTASK_NAME}.yaml"
        else
          sed "s!${MAIN_IMAGE}!${CUSTOM_IMAGE}!g" "manifests/${SUBTASK_NAME}.yaml" | oc apply -f -
        fi
      done
    leave
  done
leave
