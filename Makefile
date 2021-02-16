SHELL := /bin/bash

export AUTH_TOKEN=abc123
export LOG_LEVEL=debug

.PHONY: build run reload

build:
	go build -mod=vendor -o build/_outputs/budget-server ./cmd/server/main.go

build-linux:
	docker run -v $(PWD):/app -w /app -e CGO_ENABLED=1 --rm -it --platform linux/amd64 golang:1.15 go build -mod=vendor -o build/_outputs/budget-server-linux-amd64 cmd/server/main.go

run: build
	build/_outputs/budget-server

reload:
	hack/reload.sh 'make run' pkg/ cmd/
