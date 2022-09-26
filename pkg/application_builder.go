package pkg

import (
	"context"
	"github.com/go-redis/redis"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"tinyurl/pkg/application"
	"tinyurl/pkg/datamodel"
	"tinyurl/pkg/store"
)

type applicationManagerBuilder struct {
	ctx      context.Context
	logEntry *logrus.Entry
	db       *mongo.Client
	dbName   string
	cache    *redis.Client
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
	if builder.db != nil && len(builder.dbName) > 0 {
		storeManagerOpts = append(storeManagerOpts, store.WithDb(builder.db, builder.dbName))
	}

	if builder.cache != nil {
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

func (builder *applicationManagerBuilder) WithDb(db *mongo.Client) *applicationManagerBuilder {
	builder.db = db
	return builder
}

func (builder *applicationManagerBuilder) WithCache(cache *redis.Client) *applicationManagerBuilder {
	builder.cache = cache
	return builder
}

func (builder *applicationManagerBuilder) WithDbName(dbName string) *applicationManagerBuilder {
	builder.dbName = dbName
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
