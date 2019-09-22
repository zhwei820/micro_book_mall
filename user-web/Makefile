
GOPATH:=$(shell go env GOPATH)

.PHONY: build
build:

	go build -o user-web main.go plugin.go

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: docker
docker:
	docker build . -t user-web:latest

start:
	export ConfigAddress=192.168.1.4:9600; cd ./user-web; ps -ef|grep user-web|grep -v grep|cut -c 9-15|xargs kill -9 | true; nohup ./user-web &>nohup.out; nohup ./user-web &>nohup.out; nohup ./user-web &>nohup.out;
