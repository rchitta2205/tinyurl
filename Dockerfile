FROM golang:latest
COPY . /
WORKDIR /
RUN go build -o ./bin/tiny-url-service cmd/service/main.go
CMD ./bin/tiny-url-service