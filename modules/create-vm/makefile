all: clean build

CONFIG_FILE ?= ../../configs/create-vm.yaml
include ../../scripts/makefile-snippets/makefile-load-config.mk

include ../../scripts/makefile-snippets/makefile-clean.mk # clean
include ../../scripts/makefile-snippets/makefile-golang.mk # build lint test
include ../../scripts/makefile-snippets/makefile-docker.mk # docker-build docker-push

release: lint test docker-build

release-with-push: release docker-push

release-prod: IMAGE = $(MAIN_IMAGE)
release-prod: release

release-prod-with-push: IMAGE := $(MAIN_IMAGE)
release-prod-with-push: release-with-push

.PHONY: \
	all \
	release \
	release-with-push \
	release-prod \
	release-prod-with-push
