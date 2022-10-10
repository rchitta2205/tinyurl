package main

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net"
	"net/http"
	"os"
	"tinyurl/pkg/config"
	"tinyurl/pkg/util"
	"tinyurl/test/benchmark"
)

// Integration test for the tiny url application. We basically will have a test pod that will solely test the
// running tinyUrl service. To work with it, we need to run the test code in a different pod, connect with the
// tinyUrl service and make multiple api requests to see whether it is working correctly. We will be performing
// stress testing in this pod as well, and will have multiple goroutines do the api testing. For this purpose
// we can leverage dapr's service-to-service invocation strategy.
func main() {
	ctx := context.Background()
	cfg := config.NewConfig()
	logEntry := logrus.NewEntry(logrus.New())

	router := mux.NewRouter()
	daprClient, err := util.Connect(ctx, net.JoinHostPort(cfg.DaprAddr, cfg.DaprPort))
	if err != nil {
		logEntry.Warnf("dapr not initialized, due to error: %+v", err.Error())
		os.Exit(1)
	} else if daprClient == nil {
		logEntry.Fatal("dapr not initialized, due to internal unknown error")
		os.Exit(1)
	}
	defer daprClient.Close()

	b := benchmark.NewBenchMark(ctx, cfg, logEntry, daprClient)

	router.HandleFunc("/healthz", b.Check).Methods(http.MethodGet)

	server := &http.Server{
		Addr:    cfg.RestServerPort,
		Handler: router,
	}

	err = server.ListenAndServe()
	if err != nil {
		logEntry.Warnf("Error encountered: %s", err.Error())
	}
}
