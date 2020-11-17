TASK_SCRIPTS_DIR ?= scripts

DIST_DIR ?= dist
DIST_MANIFESTS_DIR ?= $(DIST_DIR)/manifests
DIST_EXAMPLES_DIR ?= $(DIST_DIR)/examples
MANIFESTS_DIR ?= manifests
EXAMPLES_DIR ?= examples

HAS_EXAMPLES ?= true

clean-generated-task-dist:
	if [ "$(HAS_EXAMPLES)" == true ]; then rm -rf $(DIST_EXAMPLES_DIR); fi
	rm -rf $(DIST_MANIFESTS_DIR)

clean-generated-task-release:
	if [ "$(HAS_EXAMPLES)" == true ]; then rm -rf $(EXAMPLES_DIR); fi
	rm -rf $(MANIFESTS_DIR)

generate-task:
	ansible-playbook $(TASK_SCRIPTS_DIR)/generate-task.yaml

copy-generated-task-to-release:
	mkdir -p $(MANIFESTS_DIR)
	cp $(DIST_MANIFESTS_DIR)/namespace-role/$(TASK_NAME)-namespace-rbac.yaml \
		$(DIST_MANIFESTS_DIR)/cluster-role/$(TASK_NAME)-cluster-rbac.yaml \
		$(MANIFESTS_DIR)
	set -e; $(foreach SUBTASK_NAME, $(SUBTASK_NAMES), cp $(DIST_MANIFESTS_DIR)/$(SUBTASK_NAME).yaml $(MANIFESTS_DIR);)
	if [ "$(HAS_EXAMPLES)" == true ]; then cp -r $(DIST_EXAMPLES_DIR) $(EXAMPLES_DIR); fi

test-yaml-consistency: MANIFESTS_DIR_OLD := $(MANIFESTS_DIR)
test-yaml-consistency: EXAMPLES_DIR_OLD := $(EXAMPLES_DIR)
test-yaml-consistency: MANIFESTS_DIR = $(DIST_DIR)/consistency/manifests
test-yaml-consistency: EXAMPLES_DIR = $(DIST_DIR)/consistency/examples
test-yaml-consistency: clean generate-task copy-generated-task-to-release
	diff -r $(MANIFESTS_DIR_OLD) $(MANIFESTS_DIR)
	if [ "$(HAS_EXAMPLES)" == true ]; then diff -r $(EXAMPLES_DIR_OLD) $(EXAMPLES_DIR); fi

.PHONY: \
	clean-generated-task-dist \
	clean-generated-task-release \
	generate-task \
	copy-generated-task-to-release \
	test-yaml-consistency
