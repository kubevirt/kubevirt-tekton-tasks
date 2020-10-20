export GOFLAGS=-mod=vendor
export GO111MODULE=on

FOLDERS_WITHOUT_VENDOR = $(shell ls -d */ | grep -v "^vendor/")

lint:
	@if [ -n "`gofmt -d $(FOLDERS_WITHOUT_VENDOR)`" ]; then gofmt -d $(FOLDERS_WITHOUT_VENDOR); exit 1; fi

lint-fix:
	@gofmt -w $(FOLDERS_WITHOUT_VENDOR)

vendor:
	go mod tidy
	go mod vendor

.PHONY: \
	lint \
	lint-fix \
	vendor
