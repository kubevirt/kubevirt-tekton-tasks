DIST_DIR ?= dist

export GOFLAGS=-mod=vendor
export GO111MODULE=on

FOLDERS_WITHOUT_VENDOR = $(shell ls -d */ | grep -v "^vendor/")

lint:
	if [ -n "`gofmt -d $(FOLDERS_WITHOUT_VENDOR)`" ]; then gofmt -d $(FOLDERS_WITHOUT_VENDOR); exit 1; fi

test:
	mkdir -p $(DIST_DIR)
	go test -coverprofile $(DIST_DIR)/cover.out `go list ./... | grep -v utilstest`

cover: test
	go tool cover -html $(DIST_DIR)/cover.out


vendor:
	go mod tidy
	go mod vendor

.PHONY: \
	lint \
	test \
	vendor
