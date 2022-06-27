#!/usr/bin/env bash

OFFENSIVE_WORDS="black[ -]?list|white[ -]?list|master|slave"
ALLOW_LIST=".+[:/=]master[a-zA-Z]*/?"
MODULES_VENDOR=(':!modules/tests/vendor' ':!modules/sharedtest/vendor' ':!modules/shared/vendor' ':!modules/modify-vm-template/vendor' ':!modules/generate-ssh-keys/vendor' ':!modules/execute-in-vm/vendor' ':!modules/disk-virt-sysprep/vendor' ':!modules/disk-virt-customize/vendor' ':!modules/create-vm/vendor' ':!modules/copy-template/vendor' ':!modules/wait-for-vmi-status/vendor' ':!modules/create-data-object/vendor')
MODULES_GOMOD=(':!modules/tests/go.mod' ':!modules/sharedtest/go.mod' ':!modules/shared/go.mod' ':!modules/modify-vm-template/go.mod' ':!modules/generate-ssh-keys/go.mod' ':!modules/execute-in-vm/go.mod' ':!modules/disk-virt-sysprep/go.mod' ':!modules/disk-virt-customize/go.mod' ':!modules/create-vm/go.mod' ':!modules/copy-template/go.mod' ':!modules/wait-for-vmi-status/go.mod' ':!modules/create-data-object/go.mod')

if git grep -inE "${OFFENSIVE_WORDS}" -- "${MODULES_VENDOR[@]}" "${MODULES_GOMOD[@]}" ":!${BASH_SOURCE[0]}" ':!.github/workflows/' | grep -viE "${ALLOW_LIST}"; then
  echo "Validation failed. Found offensive language"
  exit 1
fi
