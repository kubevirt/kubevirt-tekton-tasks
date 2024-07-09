#!/usr/bin/env bash

cd modules || exit 1

go test ./...

exit $?
