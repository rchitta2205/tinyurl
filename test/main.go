package main

import (
	"context"
	"fmt"
	dapr "github.com/dapr/go-sdk/client"
	"github.com/gofrs/uuid"
	"github.com/keys-pub/keys/json"
	"github.com/sirupsen/logrus"
	"net"
	"net/http"
	"os"
	"sync"
	"time"
	"tinyurl/pkg/config"
	"tinyurl/pkg/datamodel"
	"tinyurl/pkg/util"
)

const (
	appIdEnv     = "TINYURL_APP_ID"
	createMethod = "/v1/longurl"
	fetchMethod  = "/v1/tinyurl"
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

	daprClient, err := util.Connect(ctx, net.JoinHostPort(cfg.DaprAddr, cfg.DaprPort))
	if err != nil {
		logEntry.Warnf("dapr not initialized, due to error: %+v", err.Error())
		os.Exit(1)
	} else if daprClient == nil {
		logEntry.Fatal("dapr not initialized, due to internal unknown error")
		os.Exit(1)
	}
	defer daprClient.Close()
	runBenchmark(ctx, daprClient, logEntry)
	for {
		logEntry.Info("Benchmark completed...")
		time.Sleep(10 * time.Minute)
	}
}

func runBenchmark(ctx context.Context, daprClient dapr.Client, logEntry *logrus.Entry) {
	totalCases := 10000
	wg := &sync.WaitGroup{}

	for i := 0; i < totalCases; i++ {
		wg.Add(1)
		go benchmarkHelper(ctx, daprClient, logEntry, wg)
	}

	wg.Wait()
}

func benchmarkHelper(ctx context.Context, daprClient dapr.Client, logEntry *logrus.Entry, wg *sync.WaitGroup) {
	defer wg.Done()
	urlId, err := uuid.NewV4()
	if err != nil {
		logEntry.Warn("Couldn't generate uuid")
	}

	if err == nil {
		longUrl := "http://www.rahul.com" + urlId.String()
		tinyUrl := commonCall(ctx, daprClient, logEntry, longUrl, createMethod)

		// Fetch the tiny url twice to mimic db and cache extractions
		for j := 0; j < 2; j++ {
			expLongUrl := commonCall(ctx, daprClient, logEntry, tinyUrl, fetchMethod)
			if expLongUrl != longUrl {
				logEntry.Warnf("Unequal long urls: (%s, %s)", longUrl, expLongUrl)
			} else {
				logEntry.Info("Test Successful")
			}
		}
	}
}

func commonCall(ctx context.Context, daprClient dapr.Client, logEntry *logrus.Entry, url, action string) string {
	data := fmt.Sprintf(`{"url": "%s"}`, url)

	content := &dapr.DataContent{
		ContentType: "application/json",
		Data:        []byte(data),
	}

	appId := os.Getenv(appIdEnv)
	out, err := daprClient.InvokeMethodWithContent(ctx, appId, action, http.MethodPost, content)
	if err != nil {
		logEntry.Warnf("service-to-service invocation failed, due to error: %+v", err.Error())
		os.Exit(1)
	}

	result := &datamodel.Response{}
	err = json.Unmarshal(out, result)
	if err != nil {
		logEntry.Warnf("unmarshalling failed, due to error: %+v", err.Error())
		os.Exit(1)
	}

	return result.Url
}
