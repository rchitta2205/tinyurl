package pkg

import (
	"context"
	dapr "github.com/dapr/go-sdk/client"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"tinyurl/pkg/application"
	"tinyurl/pkg/datamodel"
	"tinyurl/pkg/store"
)

type applicationManagerBuilder struct {
	ctx        context.Context
	logEntry   *logrus.Entry
	daprClient dapr.Client
	db         string
	cache      string
}

func NewApplicationManagerBuilder(ctx context.Context) *applicationManagerBuilder {
	return &applicationManagerBuilder{
		ctx: ctx,
	}
}

func (builder *applicationManagerBuilder) Build() (datamodel.ApplicationManager, error) {
	err := builder.validate()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var storeManagerOpts []store.StoreManagerOption
	if builder.daprClient != nil {
		storeManagerOpts = append(storeManagerOpts, store.WithDaprClient(builder.daprClient))
	}

	if len(builder.db) != 0 {
		storeManagerOpts = append(storeManagerOpts, store.WithDb(builder.db))
	}

	if len(builder.cache) != 0 {
		storeManagerOpts = append(storeManagerOpts, store.WithCache(builder.cache))
	}

	storeMgr := store.NewStoreManager(builder.ctx, builder.logEntry, storeManagerOpts...)
	appManager := application.NewApplicationManager(builder.logEntry, storeMgr)
	return appManager, nil
}

func (builder *applicationManagerBuilder) WithLogEntry(logEntry *logrus.Entry) *applicationManagerBuilder {
	builder.logEntry = logEntry
	return builder
}

func (builder *applicationManagerBuilder) WithDaprClient(daprClient dapr.Client) *applicationManagerBuilder {
	builder.daprClient = daprClient
	return builder
}

func (builder *applicationManagerBuilder) WithCache(cache string) *applicationManagerBuilder {
	builder.cache = cache
	return builder
}

func (builder *applicationManagerBuilder) WithDb(db string) *applicationManagerBuilder {
	builder.db = db
	return builder
}

func (builder *applicationManagerBuilder) validate() error {
	if builder.logEntry == nil {
		return errors.New("cannot instantiate app mgr without log entry")
	}

	if builder.ctx == nil {
		return errors.New("cannot instantiate app mgr without ctx")
	}

	return nil
}
