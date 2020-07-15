#!/usr/bin/env bash

set -e

TMP_DV_YAML_FILENAME="/tmp/datavolume.yaml"
TMP_DV_RESULTS_FILENAME="/tmp/dv_results"

echo "$(inputs.params.manifest)" > "$TMP_DV_YAML_FILENAME"

if ! grep -q "kind: DataVolume$" "$TMP_DV_YAML_FILENAME"; then
    1>&2 echo "manifest does not contain DataVolume kind!"
    exit 1
fi

oc create -f "$TMP_DV_YAML_FILENAME" -o  jsonpath='{.metadata.name} {.metadata.namespace}' > "$TMP_DV_RESULTS_FILENAME"

sed -i 's/ /\n/g' "$TMP_DV_RESULTS_FILENAME"
readarray -t DV_OUTPUT_ARRAY < "$TMP_DV_RESULTS_FILENAME"

DV_NAME="${DV_OUTPUT_ARRAY[0]}"
DV_NAMESPACE="${DV_OUTPUT_ARRAY[1]}"

echo "$DV_NAME" > /tekton/results/name
echo "$DV_NAMESPACE" > /tekton/results/namespace

echo "Created $DV_NAME Datavolume in $DV_NAMESPACE namespace."

if [ "$(inputs.params.waitForSuccess)" == true ]; then
    echo "Waiting for Ready condition."
    # TODO: detect failed imports and don't wait until wait timeouts
    oc wait "datavolumes.cdi.kubevirt.io/$DV_NAME" --namespace="$DV_NAMESPACE" --for="condition=Ready" --timeout=720h
fi
