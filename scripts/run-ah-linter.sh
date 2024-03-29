#!/usr/bin/env bash

set -e

mkdir -p ah
AH_VERSION=$(curl -s https://api.github.com/repos/artifacthub/hub/releases | \
            jq -r '[.[]|select(.prerelease==false) | .tag_name] | sort | last')

curl -L "https://github.com/artifacthub/hub/releases/download/${AH_VERSION}/ah_${AH_VERSION:1}_linux_amd64.tar.gz" | tar -C ah/ -xzf -

./ah/ah lint -k tekton-task -p release/tasks/
#AH linter does not support git base versioning
#https://github.com/artifacthub/hub/issues/3691
#./ah lint -k tekton-pipeline -p pipelines/

rm -r ah
