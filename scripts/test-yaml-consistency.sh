#!/usr/bin/env bash

set -ex

make generate-yaml-tasks

git_porcelain="$(git status --untracked-files=no --porcelain)"
if [[ -n "${git_porcelain}" ]] ; then
  echo "There are uncommited changes"
  echo "${git_porcelain}"
  exit 1
fi
