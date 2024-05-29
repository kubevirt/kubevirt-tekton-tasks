#!/bin/bash

SCRIPT_DIR="$(dirname "$(readlink -f "$0")")"
REPO_DIR="$(realpath "${SCRIPT_DIR}/..")"

source "./scripts/common.sh"

USE_RESOLVER_IN_MANIFESTS=false make generate-pipelines

visit "${REPO_DIR}/release/pipelines"
    for PIPELINE_NAME in "windows-efi-installer" "windows-customize"; do
        oc apply -f "${PIPELINE_NAME}/${PIPELINE_NAME}.yaml"
        # uncomment accepting eula in autounattend.xml
        sed -i "s/<AcceptEula>false<\/AcceptEula>/<AcceptEula>true<\/AcceptEula>/g" "${PIPELINE_NAME}/configmaps/${PIPELINE_NAME}-configmaps.yaml"
        oc apply -f "${PIPELINE_NAME}/configmaps"
    done
leave
