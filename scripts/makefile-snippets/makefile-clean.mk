DIST_DIR ?= dist
DIST_MANIFESTS_DIR ?= $(DIST_DIR)/manifests
MANIFESTS_DIR ?= manifests

clean-dist:
	rm -rf  $(DIST_DIR)

clean-dist-manifests:
	rm -rf $(DIST_MANIFESTS_DIR)

clean-manifests:
	rm -rf $(MANIFESTS_DIR)

clean: clean-dist

.PHONY: \
	clean-dist \
	clean-dist-manifests \
	clean-manifests \
	clean
