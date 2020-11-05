DIST_DIR ?= dist

export GOFLAGS=-mod=vendor
export GO111MODULE=on

test:
	@mkdir -p $(DIST_DIR)
	@go test -coverprofile $(DIST_DIR)/cover.out `go list ./... | grep -v utilstest`

cover: test
	@go tool cover -html $(DIST_DIR)/cover.out

.PHONY: \
	test \
	cover
