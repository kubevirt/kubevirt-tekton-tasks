#!/usr/bin/env bash

GOFLAGS='' go install github.com/jstemmer/go-junit-report@latest

GOFLAGS='' go install github.com/onsi/ginkgo/v2/ginkgo@v2.27.2
GOFLAGS='' go install github.com/onsi/gomega/...@v1.38.2
