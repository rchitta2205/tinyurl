package benchmark

import (
	"context"
	errjson "encoding/json"
	"fmt"
	dapr "github.com/dapr/go-sdk/client"
	"github.com/gofrs/uuid"
	"github.com/keys-pub/keys/json"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"time"
	"tinyurl/pkg/config"
	"tinyurl/pkg/datamodel"
)

const (
	appIdEnv     = "TINYURL_APP_ID"
	createMethod = "/v1/longurl"
	fetchMethod  = "/v1/tinyurl"
)

type Benchmark struct {
	ctx        context.Context
	cfg        config.Config
	logEntry   *logrus.Entry
	daprClient dapr.Client
}

func NewBenchMark(ctx context.Context, cfg config.Config, logEntry *logrus.Entry, daprClient dapr.Client) *Benchmark {
	return &Benchmark{
		ctx:        ctx,
		logEntry:   logEntry,
		cfg:        cfg,
		daprClient: daprClient,
	}
}

func (b *Benchmark) Check(w http.ResponseWriter, r *http.Request) {
	go func() {
		time.Sleep(15 * time.Second)
		b.runBenchmark()
	}()
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("ok"))
	if err != nil {
		sendErr(w, http.StatusInternalServerError, err.Error())
	}
}

func (b *Benchmark) runBenchmark() {
	totalCases := 100000

	for i := 0; i < totalCases; i++ {
		b.benchmarkHelper()
	}
}

func (b *Benchmark) benchmarkHelper() {
	urlId, terr := uuid.NewV4()
	if terr != nil {
		b.logEntry.Warn("Couldn't generate uuid")
	}

	if terr == nil {
		longUrl := "http://www.rahul.com" + urlId.String()
		tinyUrl := b.commonCall(longUrl, createMethod)

		// Fetch the tiny url twice to mimic db and cache extractions
		for j := 0; j < 2; j++ {
			expLongUrl := b.commonCall(tinyUrl, fetchMethod)
			if expLongUrl != longUrl {
				b.logEntry.Warnf("Unequal long urls: (%s, %s)", longUrl, expLongUrl)
			} else {
				b.logEntry.Info("Test Successful")
			}
		}
	}
}

func (b *Benchmark) commonCall(url, action string) string {
	data := fmt.Sprintf(`{"url": "%s"}`, url)

	content := &dapr.DataContent{
		ContentType: "application/json",
		Data:        []byte(data),
	}

	appId := os.Getenv(appIdEnv)
	out, err := b.daprClient.InvokeMethodWithContent(b.ctx, appId, action, http.MethodPost, content)
	if err != nil {
		b.logEntry.Warnf("service-to-service invocation failed, due to error: %+v", err.Error())
		return ""
	}

	result := &datamodel.Response{}
	err = json.Unmarshal(out, result)
	if err != nil {
		b.logEntry.Warnf("unmarshalling failed, due to error: %+v", err.Error())
		return ""
	}

	return result.Url
}

func sendErr(w http.ResponseWriter, code int, message string) {
	resp, _ := errjson.Marshal(map[string]string{"error": message})
	http.Error(w, string(resp), code)
}
