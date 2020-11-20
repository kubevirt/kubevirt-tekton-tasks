MODULES_DIR = ./modules
UNIT_TESTS_DIR = $(shell ls -d $(MODULES_DIR)/* | grep -v "/tests$$")

all: clean

clean:
	./scripts/clean.sh

generate-yaml-tasks:
	./scripts/generate-yaml-tasks.sh

test-yaml-consistency:
	./scripts/test-yaml-consistency.sh

deploy:
	./scripts/deploy-tasks.sh

undeploy:
	./scripts/undeploy-tasks.sh

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

cluster-clean-and-skip-images:
	PRUNE_IMAGES=false ./scripts/cluster-clean.sh

e2e-tests:
	./automation/e2e-tests.sh

onboard-new-task-with-ci-stub:
	./scripts/onboard-new-task-with-ci-stub.sh


.PHONY: \
	all \
	clean \
	generate-yaml-tasks \
	test-yaml-consistency \
	deploy \
	undeploy \
	lint \
	lint-fix \
	test \
	test-with-reports \
	cluster-sync \
	cluster-test \
	cluster-clean \
	cluster-clean-and-skip-images \
	e2e-tests \
	onboard-new-task-with-ci-stub
