TASK_SCRIPTS_DIR ?= scripts

DIST_DIR ?= dist
DIST_MANIFESTS_DIR ?= $(DIST_DIR)/manifests
MANIFESTS_DIR ?= manifests

build-manifests:
	ansible-playbook $(TASK_SCRIPTS_DIR)/generate-manifests.yaml

copy-bundle-to-manifests:
	mkdir -p manifests
	cp $(DIST_MANIFESTS_DIR)/namespace-role/$(TASK_NAME)-namespace-rbac.yaml \
		$(DIST_MANIFESTS_DIR)/cluster-role/$(TASK_NAME)-cluster-rbac.yaml \
		$(MANIFESTS_DIR)
	set -e; $(foreach SUBTASK_NAME, $(SUBTASK_NAMES), cp $(DIST_MANIFESTS_DIR)/$(SUBTASK_NAME) $(MANIFESTS_DIR);)

.PHONY: \
	build-manifests \
	copy-bundle-to-manifests
