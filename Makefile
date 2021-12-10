GOPATH=$(shell go env GOPATH)
IMAGE_NAME=MIHP
CURRENT_PATH=$(shell pwd)
GO111MODULE=on

.PHONY: all test clean build docker

build-windows:
	GOOS=windows GOARCH=amd64 go build -a -o $(IMAGE_NAME).exe cmd/*.go

build-linux:
	GO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o $(IMAGE_NAME).app cmd/*.go

test: build-linux
	export GO111MODULE on; \
	go test ./... -cover -vet -all -v -short -covermode=count -coverprofile=coverage.out

run-setup: build-linux
	./$(IMAGE_NAME).app -setup

distribute: build-linux
	scp -C MIHP.app root@MIHP1:/root
	scp -C MIHP.app root@MIHP2:/root
	scp -C MIHP.app root@MIHP3:/root
	scp -C MIHP.app root@MIHP4:/root
