all: lint test

lint:
	if [ -n "`gofmt -d .`" ]; then gofmt -d .; exit 1; fi

test:
	go test `go list ./... | grep -v utilstest`

.PHONY: \
	all \
	lint \
	test
