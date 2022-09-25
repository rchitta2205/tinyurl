FROM golang:1.16.3-alpine3.13 AS GO_BUILD
RUN apk update && apk add build-base
COPY . /tinyurl
WORKDIR /tinyurl
RUN go build -o ./bin/tiny-url-service cmd/service/main.go
RUN go build -o ./bin/tiny-url-client cmd/client/main.go

FROM alpine:3.13.5
COPY --from=GO_BUILD /tinyurl ./
RUN ls
CMD ./bin/tiny-url-service
