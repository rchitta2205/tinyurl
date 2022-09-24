package main

import (
	"context"
	"github.com/sirupsen/logrus"
	"tinyurl/pkg"
	"tinyurl/pkg/config"
)

func main() {
	var ctx context.Context

	ctx = context.Background()
	logEntry := logrus.NewEntry(logrus.New())

	// Initialize the configuration
	cfg := config.NewConfig()
	grpcService := pkg.NewService(ctx, cfg, logEntry)

	err := grpcService.Register()
	if err != nil {
		logEntry.Fatalf("Failed to create grpc server object and register apps: %+v", err.Error())
	}

	err = grpcService.Serve()
	if err != nil {
		logEntry.Fatalf("Failed to start the server: %+v", err.Error())
	}
}
