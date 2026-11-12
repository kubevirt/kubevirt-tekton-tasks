#!/bin/bash

set -euo pipefail

SCRIPT_DIR="$(dirname "$(readlink -f "$0")")"
REPO_DIR="$(realpath "${SCRIPT_DIR}/..")"

source "./scripts/common.sh"

resolve_deploy_namespace

USE_RESOLVER_IN_MANIFESTS=false make generate-pipelines

visit "${REPO_DIR}/release/pipelines"
    for PIPELINE_NAME in "windows-efi-installer" "windows-customize"; do
        if ! oc apply -n "${DEPLOY_NAMESPACE}" -f "${PIPELINE_NAME}/${PIPELINE_NAME}.yaml"; then
            echo "Failed to apply pipeline '${PIPELINE_NAME}' to namespace '${DEPLOY_NAMESPACE}'." >&2
            exit 1
        fi
    done
leave
