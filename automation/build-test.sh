#!/usr/bin/env bash

GOFLAGS='' go install github.com/jstemmer/go-junit-report@latest

GOFLAGS='' go install github.com/onsi/ginkgo/v2/ginkgo@v2.21.0
GOFLAGS='' go install github.com/onsi/gomega/...@v1.35.1
