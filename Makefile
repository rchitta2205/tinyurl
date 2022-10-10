.PHONY: protobuf cert service client vendor build test-build up down mini mini-down

mini:
	./minikube_start.sh

mini-down:
	minikube stop && minikube delete

vendor:
	go mod tidy && go mod vendor

protobuf:
	@buf generate

cert:
	cd cert; ./gen.sh; cd ..

service:
	go build -o ./bin/tiny-url-service cmd/service/main.go

client:
	go build -o ./bin/tiny-url-client cmd/client/main.go

build:
	@docker build -t rchitta2205/tinyurl -f TinyUrlDockerfile .

test-build:
	@docker build -t rchitta2205/test-tinyurl -f TestDockerfile .

tag:
	@docker tag rchitta2205/tinyurl rchitta2205/tinyurl

test-tag:
	@docker tag rchitta2205/test-tinyurl rchitta2205/test-tinyurl

push:
	@docker push rchitta2205/tinyurl

test-push:
	@docker push rchitta2205/test-tinyurl

install:
	@helm install tinyurl ./deploy

up: build tag push test-build test-tag test-push install

down:
	@helm uninstall tinyurl
