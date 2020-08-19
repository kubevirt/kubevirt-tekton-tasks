QUAY_USER ?= $(USER)
IMAGE_REGISTRY ?= quay.io/$(QUAY_USER)
IMAGE_TAG ?= latest
IMAGE_NAME ?= $(TASK_NAME)
IMAGE ?= $(IMAGE_REGISTRY)/$(IMAGE_NAME):$(IMAGE_TAG)

docker-build:
	docker build -f build/$(IMAGE_NAME)/Dockerfile -t $(IMAGE) .

docker-push:
	docker push $(IMAGE)

.PHONY: \
	docker-build \
	docker-push
