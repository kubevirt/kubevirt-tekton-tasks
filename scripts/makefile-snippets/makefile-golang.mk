export GOFLAGS=-mod=vendor
export GO111MODULE=on

FOLDERS_WITHOUT_VENDOR = $(shell ls -d */ | grep -v "^vendor/")

lint:
	if [ -n "`gofmt -d $(FOLDERS_WITHOUT_VENDOR)`" ]; then gofmt -d $(FOLDERS_WITHOUT_VENDOR); exit 1; fi

test:
	go test `go list ./... | grep -v utilstest`

vendor:
	go mod tidy
	go mod vendor

.PHONY: \
	lint \
	test \
	vendor
