TASKS_DIR = ./tasks
MODULES_DIR = ./modules

all: $(TASKS_DIR)/* $(MODULES_DIR)/*
	set -e; $(foreach TASK_DIR, $^, $(MAKE) -C $(TASK_DIR);)

clean: $(TASKS_DIR)/* $(MODULES_DIR)/*
	set -e; $(foreach TASK_DIR, $^, $(MAKE) -C $(TASK_DIR) clean;)

release-manifests: $(TASKS_DIR)/*
	set -e; $(foreach TASK_DIR, $^, $(MAKE) -C $(TASK_DIR) release;)

undeploy: $(TASKS_DIR)/*
	set -e; $(foreach TASK_DIR, $^, $(MAKE) -C $(TASK_DIR) undeploy;)

deploy: $(TASKS_DIR)/*
	set -e; $(foreach TASK_DIR, $^, $(MAKE) -C $(TASK_DIR) deploy;)

deploy-namespace: $(TASKS_DIR)/*
	set -e; $(foreach TASK_DIR, $^, $(MAKE) -C $(TASK_DIR) deploy-namespace;)


.PHONY: \
	all \
	clean \
	release-manifests \
	undeploy \
	deploy \
	deploy-namespace
