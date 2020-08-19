DIST_DIR ?= dist
CMD_DIR ?= cmd

build:
	mkdir -p $(DIST_DIR)
	go build -o $(DIST_DIR)/$(TASK_NAME) $(CMD_DIR)/$(TASK_NAME)/main.go

lint:
	if [ -n "`gofmt -d .`" ]; then gofmt -d .; exit 1; fi

test:
	go test `go list ./... | grep -v utilstest`

.PHONY: \
	build \
	lint \
	test
