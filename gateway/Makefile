
GOPATH:=$(shell go env GOPATH)


.PHONY: build
build:

	go build -o gateway main.go plugin.go

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: docker
docker:
	docker build . -t gateway:latest
