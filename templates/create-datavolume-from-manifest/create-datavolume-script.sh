#!/usr/bin/env bash

set -e

TMP_OBJ_YAML_FILENAME="/tmp/datavolume.yaml"
TMP_OBJ_RESULTS_FILENAME="/tmp/dv_results"

echo "$(inputs.params.manifest)" > "$TMP_OBJ_YAML_FILENAME"

if [[ ! "$(cat $TMP_OBJ_YAML_FILENAME)" == *"kind: DataVolume"* ]] && [[ ! "$(cat $TMP_OBJ_YAML_FILENAME)" == *"kind: DataSource"* ]]; then
    1>&2 echo "manifest does not contain DataVolume or DataSource kind!"
    exit 1
fi

oc create -f "$TMP_OBJ_YAML_FILENAME" -o  jsonpath='{.metadata.name} {.metadata.namespace}' > "$TMP_OBJ_RESULTS_FILENAME"

sed -i 's/ /\n/g' "$TMP_OBJ_RESULTS_FILENAME"
readarray -t OBJ_OUTPUT_ARRAY < "$TMP_OBJ_RESULTS_FILENAME"

OBJ_NAME="${OBJ_OUTPUT_ARRAY[0]}"
OBJ_NAMESPACE="${OBJ_OUTPUT_ARRAY[1]}"

echo -n "$OBJ_NAME" > /tekton/results/name
echo -n "$OBJ_NAMESPACE" > /tekton/results/namespace

echo "Created $OBJ_NAME Datavolume in $OBJ_NAMESPACE namespace."

if [[ "$(inputs.params.waitForSuccess)" == true ]] && [[ "$TMP_OBJ_YAML_FILENAME" == *"kind: DataVolume"* ]]; then
    echo "Waiting for Ready condition."
    # TODO: detect failed imports and don't wait until wait timeouts
    oc wait "datavolumes.cdi.kubevirt.io/$OBJ_NAME" --namespace="$OBJ_NAMESPACE" --for="condition=Ready" --timeout=720h
fi

if [[ "$(inputs.params.waitForSuccess)" == true ]] && [[ "$TMP_OBJ_YAML_FILENAME" == *"kind: DataSource"* ]]; then
    echo "Waiting for Ready condition."
    # TODO: detect failed imports and don't wait until wait timeouts
    oc wait "datasources.cdi.kubevirt.io/$OBJ_NAME" --namespace="$OBJ_NAMESPACE" --for="condition=Ready" --timeout=720h
fi