all: clean build-manifests

CONFIG_FILE ?= ../../configs/execute-in-vm.yaml

include ../../scripts/makefile-snippets/makefile-common.mk


include ../../scripts/makefile-snippets/makefile-clean.mk # clean
include ../../scripts/makefile-snippets/makefile-generate-task.mk # clean-generated-task-dist clean-generated-task-release generate-task copy-generated-task-to-release test-generated-tasks-consistency
include ../../scripts/makefile-snippets/makefile-deploy.mk # undeploy deploy deploy-namespace
include ../../scripts/makefile-snippets/makefile-deploy-dev.mk # deploy-dev deploy-namespace-dev

release: clean-dist-manifests clean-manifests build-manifests copy-bundle-to-manifests

.PHONY: \
	all \
	release
