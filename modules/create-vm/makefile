all: clean build

CONFIG_FILE ?= ../../configs/create-vm.yaml
include ../../scripts/makefile-snippets/makefile-load-config.mk

include ../../scripts/makefile-snippets/makefile-clean.mk # clean
include ../../scripts/makefile-snippets/makefile-golang.mk # lint test vendor
include ../../scripts/makefile-snippets/makefile-golang-build.mk # build
include ../../scripts/makefile-snippets/makefile-docker.mk # docker-build docker-push

release-dev: lint test docker-build

release-dev-with-push: release docker-push

release: IMAGE = $(MAIN_IMAGE)
release: release-dev

release-with-push: IMAGE = $(MAIN_IMAGE)
release-with-push: release-dev-with-push

.PHONY: \
	all \
	release-dev \
	release-dev-with-push \
	release \
	release-with-push
