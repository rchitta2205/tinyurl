.PHONY: protobuf cert service client vendor build up down

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
	@docker build -t rchitta2205/tinyurl .

tag:
	@docker tag rchitta2205/tinyurl rchitta2205/tinyurl

push:
	@docker push rchitta2205/tinyurl

up: build tag push
	@helm install tinyurl ./deploy

down:
	@helm uninstall tinyurl
