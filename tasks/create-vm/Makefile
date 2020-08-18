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
CREATE_VM_IMAGE ?= create-vm

include ../../scripts/makefile-snippets/makefile-clean.mk # clean-dist clean-dist-manifests clean-manifests clean
include ../../scripts/makefile-snippets/makefile-manifests.mk # build-manifests copy-bundle-to-manifests
include ../../scripts/makefile-snippets/makefile-deploy.mk # undeploy deploy deploy-namespace

build:
	mkdir -p $(DIST_DIR)
	go build -o $(DIST_DIR)/$(CREATE_VM_IMAGE) cmd/create-vm/main.go

lint:
	if [ -n "`gofmt -d .`" ]; then gofmt -d .; exit 1; fi

test:
	go test `go list ./... | grep -v utilstest`

docker-build:
	docker build -f build/create-vm/Dockerfile -t $(IMAGE_REGISTRY)/$(CREATE_VM_IMAGE):$(IMAGE_TAG) .

docker-push:
	docker push $(IMAGE_REGISTRY)/$(CREATE_VM_IMAGE):$(IMAGE_TAG)

release-manifests: clean-dist-manifests clean-manifests build-manifests copy-bundle-to-manifests

release: lint test release-manifests docker-build

release-with-push: release docker-push

.PHONY: \
	all \
	build \
	lint \
	test \
	docker-build \
	docker-push \
	release-manifests \
	release \
	release-with-push
