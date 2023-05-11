.PHONY: build
build:
	swag init -g ./app/apiserver/server.go -o ./docs
	go build -v ./cmd/apiserver

.DEFAULT_GOAL := build