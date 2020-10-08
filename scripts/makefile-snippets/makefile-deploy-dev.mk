TARGET_NAMESPACE ?= $(shell kubectl config current-context | cut -d/ -f1)
MANIFESTS_DIR ?= manifests

deploy-dev: undeploy
	sed "s/TARGET_NAMESPACE/$(TARGET_NAMESPACE)/" $(MANIFESTS_DIR)/$(TASK_NAME)-cluster-rbac.yaml | kubectl apply -f -
	set -e; $(foreach SUBTASK_NAME, $(SUBTASK_NAMES), sed 's!$(MAIN_IMAGE)!$(IMAGE)!g' $(MANIFESTS_DIR)/$(SUBTASK_NAME).yaml | kubectl apply -f -;)

deploy-dev-namespace: undeploy
	kubectl apply -f manifests/$(TASK_NAME)-namespace-rbac.yaml
	set -e; $(foreach SUBTASK_NAME, $(SUBTASK_NAMES), sed 's!$(MAIN_IMAGE)!$(IMAGE)!g' $(MANIFESTS_DIR)/$(SUBTASK_NAME).yaml | kubectl apply -f -;)


.PHONY: \
	deploy-dev \
	deploy-dev-namespace
