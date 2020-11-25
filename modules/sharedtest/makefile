all: lint test

include ../../scripts/makefile-snippets/makefile-clean.mk # clean
include ../../scripts/makefile-snippets/makefile-golang.mk # lint lint-fix vendor
include ../../scripts/makefile-snippets/makefile-golang-test.mk # test cover junit-report

.PHONY: \
	all
