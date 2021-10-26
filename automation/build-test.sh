#!/usr/bin/env bash

export GO111MODULE=on
go get github.com/jstemmer/go-junit-report

go get github.com/onsi/ginkgo/ginkgo@v1.15.2
go get github.com/onsi/gomega/...@v1.15.0
go get github.com/onsi/ginkgo/extensions/table@v1.15.2
