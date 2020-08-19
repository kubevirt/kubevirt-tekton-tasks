all: clean build-manifests

DIST_DIR ?= dist
DIST_MANIFESTS_DIR ?= $(DIST_DIR)/manifests
MANIFESTS_DIR ?= manifests
TASK_SCRIPTS_DIR ?= scripts

TASK_NAME = create-datavolume
SUBTASK_NAMES = create-datavolume.yaml # add other-subtask.yaml (list)

include ../../scripts/makefile-snippets/makefile-clean.mk # clean-dist clean-dist-manifests clean-manifests clean
include ../../scripts/makefile-snippets/makefile-manifests.mk # build-manifests copy-bundle-to-manifests
include ../../scripts/makefile-snippets/makefile-deploy.mk # undeploy deploy deploy-namespace

release-manifests: clean-dist-manifests clean-manifests build-manifests copy-bundle-to-manifests

release: release-manifests

.PHONY: \
	all \
	release-manifests \
	release
