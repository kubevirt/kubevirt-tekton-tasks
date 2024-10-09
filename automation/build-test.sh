#!/usr/bin/env bash

GOFLAGS='' go install github.com/jstemmer/go-junit-report@latest

GOFLAGS='' go install github.com/onsi/ginkgo/v2/ginkgo@v2.20.1
GOFLAGS='' go install github.com/onsi/gomega/...@v1.34.2
