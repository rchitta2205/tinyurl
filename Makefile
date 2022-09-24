.PHONY: protobuf cert service client vendor

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

