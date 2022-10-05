package main

import (
	"context"
	"github.com/sirupsen/logrus"
	"net"
	"sync"
	"tinyurl/pkg"
	"tinyurl/pkg/config"
	"tinyurl/pkg/util"
)

func main() {
	var ctx context.Context

	ctx = context.Background()
	logEntry := logrus.NewEntry(logrus.New())

	// Initialize the configuration
	cfg := config.NewConfig()

	daprClient, err := util.Connect(ctx, net.JoinHostPort(cfg.DaprAddr, cfg.DaprPort))
	if err != nil || daprClient == nil {
		logEntry.Warnf("dapr not initialized, using mock db and cache: %+v", err.Error())
	}

	if daprClient != nil {
		defer daprClient.Close()
	}

	grpcService := pkg.NewGrpcService(ctx, cfg, logEntry, daprClient)
	restService := pkg.NewRestService(ctx, cfg, logEntry)

	wg := &sync.WaitGroup{}
	err = grpcService.Register()
	if err != nil {
		logEntry.Fatalf("Failed to create grpc server object and register apps: %+v", err.Error())
	}

	err = grpcService.Serve(wg)
	if err != nil {
		logEntry.Fatalf("Failed to start the gRPC server: %+v", err.Error())
	}

	err = restService.Register()
	if err != nil {
		logEntry.Fatalf("Failed to create rest server and register apps: %+v", err.Error())
	}

	err = restService.Serve(wg)
	if err != nil {
		logEntry.Fatalf("Failed to start the rest server: %+v", err.Error())
	}

	// Wait for all goroutines to be completed
	wg.Wait()
}
