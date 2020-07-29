all: clean build-manifests

QUAY_USER ?= $(USER)
IMAGE_REGISTRY ?= quay.io/$(QUAY_USER)
IMAGE_TAG ?= latest
CREATE_VM_IMAGE ?= create-vm

clean-dist:
	rm -rf dist

clean-manifests:
	rm -rf manifests

clean: clean-dist

docker-build:
	docker build -f build/create-vm/Dockerfile -t $(IMAGE_REGISTRY)/$(CREATE_VM_IMAGE):$(IMAGE_TAG) .

docker-push:
	docker push $(IMAGE_REGISTRY)/$(CREATE_VM_IMAGE):$(IMAGE_TAG)

build-manifests:
	ansible-playbook scripts/generate-manifests.yaml

copy-bundle-to-manifests:
	mkdir -p manifests
	cp dist/create-vm-from-template.yaml \
		dist/namespace-role/create-vm-namespace-rbac.yaml \
		dist/cluster-role/create-vm-cluster-rbac.yaml \
		manifests

release-manifests: clean-dist clean-manifests build-manifests copy-bundle-to-manifests

release: release-manifests docker-build

release-with-push: release docker-push

.PHONY: \
	all \
	clean-dist \
	clean-manifests \
	clean \
	docker-build \
	docker-push \
	build-manifests \
	copy-bundle-to-manifests \
	release-manifests \
	release \
	release-with-push
