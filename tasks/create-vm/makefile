all: clean build-manifests

CONFIG_FILE ?= ../../configs/create-vm.yaml
include ../../scripts/makefile-snippets/makefile-load-config.mk

include ../../scripts/makefile-snippets/makefile-clean.mk # clean-dist clean
include ../../scripts/makefile-snippets/makefile-manifests.mk # clean-dist-manifests clean-manifests build-manifests copy-bundle-to-manifests
include ../../scripts/makefile-snippets/makefile-deploy.mk # undeploy deploy deploy-namespace
include ../../scripts/makefile-snippets/makefile-deploy-dev.mk # deploy-dev deploy-namespace-dev

release: clean-dist-manifests clean-manifests build-manifests copy-bundle-to-manifests

.PHONY: \
	all \
	release
