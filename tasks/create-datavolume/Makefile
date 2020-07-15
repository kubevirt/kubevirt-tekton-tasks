all: clean build

clean-dist:
	rm -rf dist

clean-manifests:
	rm -rf manifests

clean: clean-dist

build:
	ansible-playbook hack/generate-manifests.yaml

copy-bundle-to-manifests:
	mkdir -p manifests
	cp dist/create-datavolume.yaml \
		dist/namespace-role/create-datavolume-namespace-rbac.yaml \
		dist/cluster-role/create-datavolume-cluster-rbac.yaml \
		manifests

release: clean-dist clean-manifests build copy-bundle-to-manifests

.PHONY: \
	all \
	clean-dist \
	clean-manifests \
	clean \
	build \
	copy-bundle-to-manifests \
	release
