
GOPATH:=$(shell go env GOPATH)


.PHONY: build
build:

	go build -o config-grpc-srv main.go

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: docker
docker:
	docker build . -t config-grpc-srv:latest
