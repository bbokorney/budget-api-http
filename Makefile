SHELL := /bin/bash

export AUTH_TOKEN=abc123
export LOG_LEVEL=debug
export BUILD_CONTAINER=go-build-container

.PHONY: build run reload

build:
	go build -mod=vendor -o build/_outputs/budget-server ./cmd/server/main.go

run: build
	build/_outputs/budget-server

reload:
	hack/reload.sh 'make run' pkg/ cmd/

start-build-container:
	docker run --name $(BUILD_CONTAINER) -v $(PWD):/app -w /app -e CGO_ENABLED=1 --rm -d --platform linux/amd64 golang:1.15 tail -f /dev/null

stop-build-container:
	docker rm -f $(BUILD_CONTAINER)

build-linux:
	docker exec -it $(BUILD_CONTAINER) go  build -mod=vendor -o build/_outputs/budget-server-linux-amd64 cmd/server/main.go
