#!/usr/bin/env bash
if [[ "$(go fmt ./...)" != "" ]]; then
    echo "Please go fmt your code"
    exit 1
fi

if [[ "$(go vet ./...)" != "" ]]; then
    echo "Please go vet your code"
    exit 1
fi

if [[ "$(golint ./...)" != "" ]]; then
    echo "Please golint your code"
    exit 1
fi