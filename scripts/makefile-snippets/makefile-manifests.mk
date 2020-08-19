TASK_SCRIPTS_DIR ?= scripts

DIST_DIR ?= dist
DIST_MANIFESTS_DIR ?= $(DIST_DIR)/manifests
MANIFESTS_DIR ?= manifests

clean-dist-manifests:
	rm -rf $(DIST_MANIFESTS_DIR)

clean-manifests:
	rm -rf $(MANIFESTS_DIR)

build-manifests:
	ansible-playbook $(TASK_SCRIPTS_DIR)/generate-manifests.yaml

copy-bundle-to-manifests:
	mkdir -p manifests
	cp $(DIST_MANIFESTS_DIR)/namespace-role/$(TASK_NAME)-namespace-rbac.yaml \
		$(DIST_MANIFESTS_DIR)/cluster-role/$(TASK_NAME)-cluster-rbac.yaml \
		$(MANIFESTS_DIR)
	set -e; $(foreach SUBTASK_NAME, $(SUBTASK_NAMES), cp $(DIST_MANIFESTS_DIR)/$(SUBTASK_NAME).yaml $(MANIFESTS_DIR);)

.PHONY: \
	clean-dist-manifests \
	clean-manifests \
	build-manifests \
	copy-bundle-to-manifests

