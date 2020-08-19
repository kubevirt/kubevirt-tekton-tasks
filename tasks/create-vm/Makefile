all: clean build-manifests build

DIST_DIR ?= dist
DIST_MANIFESTS_DIR ?= $(DIST_DIR)/manifests
MANIFESTS_DIR ?= manifests
TASK_SCRIPTS_DIR ?= scripts

TASK_NAME = create-vm
SUBTASK_NAMES = create-vm-from-template.yaml # add other-subtask.yaml (list)

QUAY_USER ?= $(USER)
IMAGE_REGISTRY ?= quay.io/$(QUAY_USER)
IMAGE_TAG ?= latest
IMAGE_NAME ?= $(TASK_NAME)

include ../../scripts/makefile-snippets/makefile-clean.mk # clean-dist clean-dist-manifests clean-manifests clean
include ../../scripts/makefile-snippets/makefile-manifests.mk # build-manifests copy-bundle-to-manifests
include ../../scripts/makefile-snippets/makefile-deploy.mk # undeploy deploy deploy-namespace
include ../../scripts/makefile-snippets/makefile-golang.mk # build lint test
include ../../scripts/makefile-snippets/makefile-docker.mk # docker-build docker-push

release-manifests: clean-dist-manifests clean-manifests build-manifests copy-bundle-to-manifests

release: lint test release-manifests docker-build

release-with-push: release docker-push

.PHONY: \
	all \
	release-manifests \
	release \
	release-with-push
