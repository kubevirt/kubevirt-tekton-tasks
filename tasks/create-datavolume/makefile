all: clean generate-task

CONFIG_FILE ?= ../../configs/create-datavolume.yaml
HAS_EXAMPLES = false

include ../../scripts/makefile-snippets/makefile-common.mk

include ../../scripts/makefile-snippets/makefile-clean.mk # clean-dist clean
include ../../scripts/makefile-snippets/makefile-generate-task.mk # clean-generated-task-dist clean-generated-task-release generate-task copy-generated-task-to-release
include ../../scripts/makefile-snippets/makefile-deploy.mk # undeploy deploy deploy-namespace

release: \
	clean \
	generate-task \
	clean-generated-task-release \
	copy-generated-task-to-release \
	clean

deploy-dev: deploy

deploy-dev-namespace: deploy-namespace

.PHONY: \
	all \
	release \
	deploy-dev \
	deploy-namespace-dev
