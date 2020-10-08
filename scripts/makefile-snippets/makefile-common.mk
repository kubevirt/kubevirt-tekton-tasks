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

SUBTASK_NAMES ?= $(shell sed -n -e  '/^subtask_names *: */,/^ *^[-]/p' $(CONFIG_FILE) | sed -n  's/^ *-//p')

ifeq ($(strip $(SUBTASK_NAMES)),)
$(error SUBTASK_NAMES is empty, at least one subtask has to be defined)
endif

CONTAINER_ENGINE ?=  $(shell \
	if podman ps >/dev/null; then \
	  echo podman ; \
    elif docker ps >/dev/null; then \
      echo docker ; \
	fi)

ifeq ($(strip $(CONTAINER_ENGINE)),)
$(error no working container runtime found. Neither docker nor podman seems to work.)
endif

IMAGE_REGISTRY ?= quay.io
IMAGE_REGISTRY_USER ?= $(USER)
IMAGE_TAG ?= latest
IMAGE_NAME ?= $(TASK_NAME)
IMAGE ?= $(IMAGE_REGISTRY)/$(IMAGE_REGISTRY_USER)/$(IMAGE_NAME):$(IMAGE_TAG)
