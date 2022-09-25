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
	@docker-compose build

up: build
	@docker-compose up

down:
	@docker-compose down
