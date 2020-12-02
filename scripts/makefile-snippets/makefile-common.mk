ifndef CONFIG_FILE
$(error CONFIG_FILE is not set)
endif

TASK_NAME ?= $(shell sed -n  's/^task_name *: *//p' $(CONFIG_FILE))

ifeq ($(strip $(TASK_NAME)),)
$(error TASK_NAME is empty)
endif

MAIN_IMAGE ?= $(shell sed -n  's/^main_image *: *//p' $(CONFIG_FILE))

ifeq ($(strip $(MAIN_IMAGE)),)
$(error MAIN_IMAGE is empty)
endif

CONTAINER_ENGINE ?=  $(shell \
	if podman ps >/dev/null; then \
	  echo podman ; \
    elif docker ps >/dev/null; then \
      echo docker ; \
    else \
      echo 'no-container-engine-found:'; \
	fi)

IMAGE_REGISTRY ?= quay.io
IMAGE_REGISTRY_USER ?= $(USER)
IMAGE_TAG ?= latest
IMAGE_NAME ?= $(TASK_NAME)
IMAGE ?= $(IMAGE_REGISTRY)/$(IMAGE_REGISTRY_USER)/$(IMAGE_NAME):$(IMAGE_TAG)
