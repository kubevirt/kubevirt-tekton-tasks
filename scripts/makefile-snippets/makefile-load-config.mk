ifndef CONFIG_FILE
$(error CONFIG_FILE is not set)
endif

TASK_NAME ?= $(shell sed -n  's/^task_name *: *//p' $(CONFIG_FILE))

ifeq ($(strip $(TASK_NAME)),)
$(error TASK_NAME is empty)
endif

MAIN_IMAGE ?= $(shell sed -n  's/^main_image *: *//p' $(CONFIG_FILE))
MAIN_IMAGE_WITHOUT_TAG ?= $(shell echo $(MAIN_IMAGE) | sed 's/:.*$$//')

ifeq ($(strip $(MAIN_IMAGE)),)
$(error MAIN_IMAGE is empty)
endif

SUBTASK_NAMES ?= $(shell sed -n -e  '/^subtask_names *: */,/^ *^[-]/p' $(CONFIG_FILE) | sed -n  's/^ *-//p')

ifeq ($(strip $(SUBTASK_NAMES)),)
$(error SUBTASK_NAMES is empty, at least one subtask has to be defined)
endif

