all: clean build

CONFIG_FILE ?= ../../configs/execute-in-vm.yaml
include ../../scripts/makefile-snippets/makefile-common.mk

include ../../scripts/makefile-snippets/makefile-clean.mk # clean
include ../../scripts/makefile-snippets/makefile-golang.mk # lint test vendor
include ../../scripts/makefile-snippets/makefile-golang-build.mk # build

release-dev:
	$(CONTAINER_ENGINE) build -f build/$(IMAGE_NAME)/Dockerfile -t $(IMAGE) .

release-dev-with-push: release-dev
	$(CONTAINER_ENGINE) push $(IMAGE) $(ARGS)

release: lint test
	$(CONTAINER_ENGINE) build -f build/$(IMAGE_NAME)/Dockerfile -t $(MAIN_IMAGE) .

release-with-push: release
	$(CONTAINER_ENGINE) push $(MAIN_IMAGE) $(ARGS)


.PHONY: \
	all \
	release \
	release-with-push \
	release-prod \
	release-prod-with-push
