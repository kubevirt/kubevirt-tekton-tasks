#!/usr/bin/env bash

function visit() {
  pushd "${1}" > /dev/null
}

function leave() {
  popd > /dev/null
}
