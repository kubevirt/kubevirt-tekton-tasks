all: lint

include ../../scripts/makefile-snippets/makefile-clean.mk # clean
include ../../scripts/makefile-snippets/makefile-golang.mk # lint lint-fix vendor

DEPLOY_NAMESPACE ?= $(shell kubectl config current-context | cut -d/ -f1)
TEST_NAMESPACE ?= $(DEPLOY_NAMESPACE)
SCOPE ?= "cluster"
DEBUG ?= "false"
STORAGE_CLASS ?= ""

e2e-tests:
	# mkdir -p ${DIST_DIR}
	# touch ./${DIST_DIR}/test.coverprofile
	# --cover --coverprofile ./${DIST_DIR}/test.coverprofile
	ginkgo  -r -p --randomizeAllSpecs --randomizeSuites --failOnPending --trace --race -- \
		--deploy-namespace=${DEPLOY_NAMESPACE} \
		--test-namespace=${TEST_NAMESPACE} \
		--kubeconfig-path=${KUBECONFIG} \
		--scope=${SCOPE} \
		--storage-class=${STORAGE_CLASS}
		--debug=${DEBUG}

.PHONY: \
	e2e-tests \
	all
