
GOPATH:=$(shell go env GOPATH)

.PHONY: buildall
buildall:
	cd ./auth; make build; make docker;
	cd ./config-grpc-srv; make build; make docker;
	cd ./inventory-srv; make build; make docker;
	cd ./orders-srv; make build; make docker;
	cd ./orders-web; make build; make docker;
	cd ./payment-srv; make build; make docker;
	cd ./payment-web; make build; make docker;
	cd ./user-srv; make build; make docker;
	cd ./user-web; make build; make docker;

buildallraw:
	cd ./auth; make build; make build;
	cd ./config-grpc-srv; make build; make build;
	cd ./inventory-srv; make build; make build;
	cd ./orders-srv; make build; make build;
	cd ./orders-web; make build; make build;
	cd ./payment-srv; make build; make build;
	cd ./payment-web; make build; make build;
	cd ./user-srv; make build; make build;
	cd ./user-web; make build; make build;

.PHONY: clean
clean:
	cd ./config-grpc-srv; rm ./config-grpc-srv | true;rm ./nohup.out | true;
	cd ./auth; rm ./auth-srv | true;rm ./nohup.out | true;
	cd ./inventory-srv; rm ./inventory-srv | true;rm ./nohup.out | true;
	cd ./orders-srv; rm ./orders-srv | true;rm ./nohup.out | true;
	cd ./orders-web; rm ./orders-web | true;rm ./nohup.out | true;
	cd ./payment-srv; rm ./payment-srv | true;rm ./nohup.out | true;
	cd ./payment-web; rm ./payment-web | true;rm ./nohup.out | true;
	cd ./user-srv; rm ./user-srv | true;rm ./nohup.out | true;
	cd ./user-web; rm ./user-web | true; rm ./nohup.out | true;

restartall:
	export ConfigAddress=192.168.1.4:9600; cd ./config-grpc-srv; ps -ef|grep config-grpc-srv|grep -v grep|cut -c 9-15|xargs kill -9 | true; nohup ./config-grpc-srv &
	export ConfigAddress=192.168.1.4:9600; cd ./auth; ps -ef|grep auth-srv|grep -v grep|cut -c 9-15|xargs kill -9 | true; nohup ./auth-srv &>nohup.out; nohup ./auth-srv &>nohup.out; nohup ./auth-srv &>nohup.out;
	export ConfigAddress=192.168.1.4:9600; cd ./inventory-srv; ps -ef|grep inventory-srv|grep -v grep|cut -c 9-15|xargs kill -9 | true; nohup ./inventory-srv &>nohup.out; nohup ./inventory-srv &>nohup.out; nohup ./inventory-srv &>nohup.out;
	export ConfigAddress=192.168.1.4:9600; cd ./orders-srv; ps -ef|grep orders-srv|grep -v grep|cut -c 9-15|xargs kill -9 | true; nohup ./orders-srv &>nohup.out; nohup ./orders-srv &>nohup.out; nohup ./orders-srv &>nohup.out;
	export ConfigAddress=192.168.1.4:9600; cd ./orders-web; ps -ef|grep orders-web|grep -v grep|cut -c 9-15|xargs kill -9 | true; nohup ./orders-web &>nohup.out; nohup ./orders-web &>nohup.out; nohup ./orders-web &>nohup.out;
	export ConfigAddress=192.168.1.4:9600; cd ./payment-srv; ps -ef|grep payment-srv|grep -v grep|cut -c 9-15|xargs kill -9 | true; nohup ./payment-srv &>nohup.out; nohup ./payment-srv &>nohup.out; nohup ./payment-srv &>nohup.out;
	export ConfigAddress=192.168.1.4:9600; cd ./payment-web; ps -ef|grep payment-web|grep -v grep|cut -c 9-15|xargs kill -9 | true; nohup ./payment-web &>nohup.out; nohup ./payment-web &>nohup.out; nohup ./payment-web &>nohup.out;
	export ConfigAddress=192.168.1.4:9600; cd ./user-srv; ps -ef|grep user-srv|grep -v grep|cut -c 9-15|xargs kill -9 | true; nohup ./user-srv &>nohup.out; nohup ./user-srv &>nohup.out; nohup ./user-srv &>nohup.out;
	export ConfigAddress=192.168.1.4:9600; cd ./user-web; ps -ef|grep user-web|grep -v grep|cut -c 9-15|xargs kill -9 | true; nohup ./user-web &>nohup.out; nohup ./user-web &>nohup.out; nohup ./user-web &>nohup.out;
