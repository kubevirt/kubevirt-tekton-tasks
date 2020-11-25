#!/usr/bin/env bash

set -e

PRUNE_IMAGES="${PRUNE_IMAGES:-true}"

IMAGE_REGISTRY="$(oc get route default-route -n openshift-image-registry --template='{{ .spec.host }}')"

if [[ "$PRUNE_IMAGES" == "true" ]]; then
  oc adm prune images --registry-url=$IMAGE_REGISTRY  --all=false  --keep-younger-than=0s --keep-tag-revisions=0 --confirm
fi

pushd modules
  for TEST_NS in $(oc get namespaces -o name | grep -Eo "e2e-tests-[0-9]{5}"); do
    oc project "$TEST_NS"
    for MODULE_DIR in $(ls | grep -vE "^(${EXCLUDED_NON_IMAGE_MODULES})$"); do
      podman rmi "$IMAGE_REGISTRY/$TEST_NS/`basename $MODULE_DIR`:latest" || true
    done

    for RESOURCE in pipelineresources.tekton.dev \
      pipelineruns.tekton.dev \
      pipelines.tekton.dev \
      taskruns.tekton.dev \
      virtualmachines.kubevirt.io \
      templates.template.openshift.io; do
        oc delete $RESOURCE --all
    done

    oc delete namespace "$TEST_NS"
  done
popd

oc project default
