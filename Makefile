.PHONY: build run

build:
	go build -mod=vendor -o build/_outputs/budget-server ./cmd/server/main.go

run: build
	build/_outputs/budget-server

reload:
	hack/reload.sh 'make run' pkg/ cmd/
