TASKS_DIR = ./tasks
MODULES_DIR = ./modules

all: $(TASKS_DIR)/* $(MODULES_DIR)/*
	set -e; $(foreach TASK_DIR, $^, $(MAKE) -C $(TASK_DIR);)

clean: $(TASKS_DIR)/* $(MODULES_DIR)/*
	set -e; $(foreach TASK_DIR, $^, $(MAKE) -C $(TASK_DIR) clean;)

release-manifests: $(TASKS_DIR)/*
	set -e; $(foreach TASK_DIR, $^, $(MAKE) -C $(TASK_DIR) release;)

undeploy: $(TASKS_DIR)/*
	set -e; $(foreach TASK_DIR, $^, $(MAKE) -C $(TASK_DIR) undeploy;)

deploy: $(TASKS_DIR)/*
	set -e; $(foreach TASK_DIR, $^, $(MAKE) -C $(TASK_DIR) deploy;)

deploy-namespace: $(TASKS_DIR)/*
	set -e; $(foreach TASK_DIR, $^, $(MAKE) -C $(TASK_DIR) deploy-namespace;)

deploy-dev: $(TASKS_DIR)/*
	set -e; $(foreach TASK_DIR, $^, $(MAKE) -C $(TASK_DIR) deploy-dev;)

deploy-dev-namespace: $(TASKS_DIR)/*
	set -e; $(foreach TASK_DIR, $^, $(MAKE) -C $(TASK_DIR) deploy-dev-namespace;)

test-generated-tasks-consistency: $(TASKS_DIR)/*
	set -e; $(foreach TASK_DIR, $^, $(MAKE) -C $(TASK_DIR) test-generated-tasks-consistency;)

lint: $(MODULES_DIR)/*
	set -e; $(foreach TASK_DIR, $^, $(MAKE) -C $(TASK_DIR) lint;)

lint-fix: $(MODULES_DIR)/*
	set -e; $(foreach TASK_DIR, $^, $(MAKE) -C $(TASK_DIR) lint-fix;)

test: $(MODULES_DIR)/*
	set -e; $(foreach TASK_DIR, $^, $(MAKE) -C $(TASK_DIR) test;)

cluster-sync:
	./scripts/cluster-sync.sh

e2e-tests:
	./automation/e2e-tests.sh

.PHONY: \
	all \
	clean \
	release-manifests \
	undeploy \
	deploy \
	deploy-namespace \
	deploy-dev \
	deploy-dev-namespace \
	test-generated-tasks-consistency \
	lint \
	lint-fix \
	test \
	cluster-sync \
	e2e-tests
