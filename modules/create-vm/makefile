all: clean build

include ./configs/vars.mk

include ../../scripts/makefile-snippets/makefile-clean.mk # clean
include ../../scripts/makefile-snippets/makefile-golang.mk # build lint test
include ../../scripts/makefile-snippets/makefile-docker.mk # docker-build docker-push

release: lint test docker-build

release-with-push: release docker-push

.PHONY: \
	all \
	release \
	release-with-push
