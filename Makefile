SHELL := /bin/bash
.PHONY: build run reload

build:
	go build -mod=vendor -o build/_outputs/budget-server ./cmd/server/main.go

run: build
	build/_outputs/budget-server

reload:
	LOG_LEVEL=debug
	hack/reload.sh 'make run' pkg/ cmd/
