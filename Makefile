TASKS_DIR = ./tasks

all: build

build: $(TASKS_DIR)/*
	$(MAKE) -C $^

clean: $(TASKS_DIR)/*
	$(MAKE) -C $^ clean

release: $(TASKS_DIR)/*
	rm -rf manifests
	mkdir -p manifests
	$(MAKE) -C $^ release
	cp $^/manifests/* manifests

.PHONY: \
	all \
	build \
	clean \
	release
