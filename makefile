TASKS_DIR = ./tasks
MODULES_DIR = ./modules
UNIT_TESTS_DIR = $(shell ls -d $(MODULES_DIR)/* | grep -v "/tests$$")

all: $(TASKS_DIR)/* $(MODULES_DIR)/*
	set -e; $(foreach TASK_DIR, $^, $(MAKE) -C $(TASK_DIR);)

clean: $(TASKS_DIR)/* $(MODULES_DIR)/*
	rm -rf dist
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

test-yaml-consistency:
	./scripts/test-yaml-consistency.sh

lint: $(MODULES_DIR)/*
	set -e; $(foreach MODULE_DIR, $^, $(MAKE) -C $(MODULE_DIR) lint;)

lint-fix: $(MODULES_DIR)/*
	set -e; $(foreach MODULE_DIR, $^, $(MAKE) -C $(MODULE_DIR) lint-fix;)

test: $(UNIT_TESTS_DIR)
	set -e; $(foreach UNIT_TEST_DIR, $^, $(MAKE) -C $(UNIT_TEST_DIR) test;)

test-with-reports:
	./scripts/test-with-reports.sh

cluster-sync:
	./scripts/cluster-sync.sh

cluster-test:
	./scripts/cluster-test.sh

cluster-clean:
	./scripts/cluster-clean.sh

cluster-clean-without-images:
	PRUNE_IMAGES=false ./scripts/cluster-clean.sh

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
	test-yaml-consistency \
	lint \
	lint-fix \
	test \
	test-with-reports \
	cluster-sync \
	cluster-test \
	cluster-clean \
	cluster-clean-without-images \
	e2e-tests
