#!/usr/bin/env bash

set -e

export NAMESPACE="${NAMESPACE:-}"
IMAGE_REGISTRY="$(oc get route default-route -n openshift-image-registry --template='{{ .spec.host }}')"

oc adm prune images --registry-url=$IMAGE_REGISTRY  --all=false  --keep-younger-than=0s --keep-tag-revisions=0 --confirm

pushd modules
  for TEST_NS in $NAMESPACE $(oc get namespaces -o name | grep -Eo "e2e-tests-[0-9]{5}"); do
    oc project "$TEST_NS"
    for MODULE_DIR in $(ls | grep -v "^shared$"); do
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

    if [[ "$NAMESPACE" != "$TEST_NS" ]]; then
      oc delete namespace "$TEST_NS"
    fi
  done
popd

oc project default
