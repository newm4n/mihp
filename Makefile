GOPATH=$(shell go env GOPATH)
IMAGE_NAME ?= $(shell basename `pwd`)
CURRENT_PATH=$(shell pwd)
GO111MODULE=on

.PHONY: all test clean build docker

build:
#	export GO111MODULE=on; \
#	GO_ENABLED=0 go build -a -o $(IMAGE_NAME).app cmd/main/Main.go
#   Use bellow if you're running on linux.
	GO_ENABLED=0 go build -a -o $(IMAGE_NAME).app cmd/Main.go cmd/SetupConfig.go

test: build
	export GO111MODULE on; \
	go test ./... -cover -vet -all -v -short -covermode=count -coverprofile=coverage.out

