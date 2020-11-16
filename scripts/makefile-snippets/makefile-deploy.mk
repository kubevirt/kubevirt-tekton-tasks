TARGET_NAMESPACE ?= $(shell oc project --short)
MANIFESTS_DIR ?= manifests

undeploy:
	oc delete -f $(MANIFESTS_DIR) 2> /dev/null || echo "undeployed only available resources"

deploy: undeploy
	sed "s/TARGET_NAMESPACE/$(TARGET_NAMESPACE)/" $(MANIFESTS_DIR)/$(TASK_NAME)-cluster-rbac.yaml | oc apply -f -
	set -e; $(foreach SUBTASK_NAME, $(SUBTASK_NAMES), oc apply -f $(MANIFESTS_DIR)/$(SUBTASK_NAME).yaml;)

deploy-namespace: undeploy
	oc apply -f manifests/$(TASK_NAME)-namespace-rbac.yaml
	set -e; $(foreach SUBTASK_NAME, $(SUBTASK_NAMES), oc apply -f $(MANIFESTS_DIR)/$(SUBTASK_NAME).yaml;)


.PHONY: \
	undeploy \
	deploy \
	deploy-namespace
