#!/usr/bin/env bash

GOFLAGS='' go install github.com/jstemmer/go-junit-report@latest

GOFLAGS='' go install github.com/onsi/ginkgo/v2/ginkgo@v2.17.1
GOFLAGS='' go install github.com/onsi/gomega/...@v1.32.0
