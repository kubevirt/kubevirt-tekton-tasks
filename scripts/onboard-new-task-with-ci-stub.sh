#!/usr/bin/env bash

set -e

SCRIPT_DIR="$(dirname "$(readlink -f "$0")")"
REPO_DIR="$(realpath "${SCRIPT_DIR}/..")"

read -p "What is the name of the new task: " TASK_NAME

if ! echo "${TASK_NAME}" |  grep -qE "^[a-z-]+$"; then
  echo "Invalid name! Should comply with ^[a-z-]+$ regex" 1>&2
  exit 1
fi

read -p "What is the name of the env variable for this task: " TASK_ENV_VAR

if ! echo "${TASK_ENV_VAR}" |  grep -qE "^[A-Z_]+_IMAGE$"; then
  echo "Invalid env variable name! Should comply with ^[A-Z_]+_IMAGE$ regex" 1>&2
  exit 1
fi

CONFIG_FILE="${REPO_DIR}/configs/${TASK_NAME}.yaml"

if [ ! -f "${CONFIG_FILE}" ]; then
echo "creating ${CONFIG_FILE}"
cat <<EOF > "${CONFIG_FILE}"
task_name: ${TASK_NAME}
task_category: ${TASK_NAME}
main_image: quay.io/kubevirt/tekton-tasks
EOF
fi

mkdir -p "${REPO_DIR}/modules/${TASK_NAME}"

echo "Update build/Containerfile file with new task name!"
