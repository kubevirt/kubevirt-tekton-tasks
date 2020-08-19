all: clean build-manifests

CONFIG_FILE ?= ../../configs/create-datavolume.yaml
include ../../scripts/makefile-snippets/makefile-load-config.mk

include ../../scripts/makefile-snippets/makefile-clean.mk # clean-dist clean
include ../../scripts/makefile-snippets/makefile-manifests.mk # clean-dist-manifests clean-manifests build-manifests copy-bundle-to-manifests
include ../../scripts/makefile-snippets/makefile-deploy.mk # undeploy deploy deploy-namespace

release: clean-dist-manifests clean-manifests build-manifests copy-bundle-to-manifests

.PHONY: \
	all \
	release
