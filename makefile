#current version of tekton tasks
export VERSION ?=v0.0.9

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

lint:
	./scripts/lint.sh

lint-fix:
	./scripts/lint-fix.sh

test:
	./scripts/test.sh

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

build-release-images:
	./scripts/build-release-images.sh

push-release-images:
	./scripts/push-release-images.sh

release: generate-yaml-tasks build-release-images push-release-images

vendor:
	./scripts/vendor.sh

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
	release \
	e2e-tests \
	onboard-new-task-with-ci-stub \
	vendor
