all: clean generate-task

CONFIG_FILE ?= ../../configs/create-vm.yaml
HAS_EXAMPLES = false

ifdef CREATE_VM_IMAGE
IMAGE = $(CREATE_VM_IMAGE)
endif

include ../../scripts/makefile-snippets/makefile-common.mk

include ../../scripts/makefile-snippets/makefile-clean.mk # clean
include ../../scripts/makefile-snippets/makefile-generate-task.mk # clean-generated-task-dist clean-generated-task-release generate-task copy-generated-task-to-release test-generated-tasks-consistency
include ../../scripts/makefile-snippets/makefile-deploy.mk # undeploy deploy deploy-namespace
include ../../scripts/makefile-snippets/makefile-deploy-dev.mk # deploy-dev deploy-dev-namespace

release: \
	clean \
	generate-task \
	clean-generated-task-release \
	copy-generated-task-to-release \
	clean

.PHONY: \
	all \
	release
