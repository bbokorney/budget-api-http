SHELL := /bin/bash

export AUTH_TOKEN=abc123
export LOG_LEVEL=debug

.PHONY: build run reload

build:
	go build -mod=vendor -o build/_outputs/budget-server ./cmd/server/main.go

run: build
	build/_outputs/budget-server

reload:
	hack/reload.sh 'make run' pkg/ cmd/
