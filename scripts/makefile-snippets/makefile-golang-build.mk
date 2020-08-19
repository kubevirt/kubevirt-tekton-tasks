DIST_DIR ?= dist
CMD_DIR ?= cmd

export GOFLAGS=-mod=vendor
export GO111MODULE=on

build:
	mkdir -p $(DIST_DIR)
	go build -o $(DIST_DIR)/$(TASK_NAME) $(CMD_DIR)/$(TASK_NAME)/main.go

.PHONY: \
	build
