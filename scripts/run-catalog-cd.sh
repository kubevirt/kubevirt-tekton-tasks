#!/usr/bin/env bash

set -x

if [ -z ${RELEASE_VERSION}]; then
    echo "tag not defined"
    exit 1
fi

catalog_cd_version=$(curl -s https://api.github.com/repos/openshift-pipelines/catalog-cd/releases | \
        jq '.[] | select(.prerelease==false) | .tag_name' | sort -V | tail -n1 | tr -d '"')

mkdir catalog-cd

curl -L https://github.com/openshift-pipelines/catalog-cd/releases/download/${catalog_cd_version}/catalog-cd_${catalog_cd_version:1}_linux_x86_64.tar.gz | \
        tar -C catalog-cd -xvzf -

task_names=""

for task_name in tasks/*/; do
    task_names="${task_names} ${task_name}"
done

mkdir release

./catalog-cd/catalog-cd release --output release --version=${RELEASE_VERSION} ${task_names}

# clean up
rm -r catalog-cd
rm -r release