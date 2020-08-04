all: clean build-manifests build

QUAY_USER ?= $(USER)
IMAGE_REGISTRY ?= quay.io/$(QUAY_USER)
IMAGE_TAG ?= latest
CREATE_VM_IMAGE ?= create-vm

DIST_DIR = dist
DIST_MANIFESTS_DIR = $(DIST_DIR)/manifests
MANIFESTS_DIR = manifests

clean-dist:
	rm -rf  $(DIST_DIR)

clean-dist-manifests:
	rm -rf $(DIST_MANIFESTS_DIR)

clean-manifests:
	rm -rf $(MANIFESTS_DIR)

clean: clean-dist

build:
	mkdir -p $(DIST_DIR)
	go build -o $(DIST_DIR)/$(CREATE_VM_IMAGE) cmd/create-vm/main.go

docker-build:
	docker build -f build/create-vm/Dockerfile -t $(IMAGE_REGISTRY)/$(CREATE_VM_IMAGE):$(IMAGE_TAG) .

docker-push:
	docker push $(IMAGE_REGISTRY)/$(CREATE_VM_IMAGE):$(IMAGE_TAG)

build-manifests:
	ansible-playbook scripts/generate-manifests.yaml

copy-bundle-to-manifests:
	mkdir -p manifests
	cp $(DIST_MANIFESTS_DIR)/create-vm-from-template.yaml \
		$(DIST_MANIFESTS_DIR)/namespace-role/create-vm-namespace-rbac.yaml \
		$(DIST_MANIFESTS_DIR)/cluster-role/create-vm-cluster-rbac.yaml \
		$(MANIFESTS_DIR)

release-manifests: clean-dist-manifests clean-manifests build-manifests copy-bundle-to-manifests

release: release-manifests docker-build

release-with-push: release docker-push

.PHONY: \
	all \
	clean-dist \
	clean-dist-manifests \
	clean-manifests \
	clean \
	build \
	docker-build \
	docker-push \
	build-manifests \
	copy-bundle-to-manifests \
	release-manifests \
	release \
	release-with-push
