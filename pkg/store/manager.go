package store

import (
	"context"
	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"tinyurl/pkg/datamodel"
)

type storeManager struct {
	ctx      context.Context
	logEntry *logrus.Entry
	dbName   string
	db       *mongo.Client
	cache    *redis.Client

	// List of all singleton stores
	tinyUrlStore datamodel.TinyUrlStore
	authStore    datamodel.AuthStore
}

type StoreManagerOption func(am *storeManager)

func NewStoreManager(ctx context.Context, logEntry *logrus.Entry, opts ...StoreManagerOption) datamodel.StoreManager {
	sm := &storeManager{
		ctx:      ctx,
		logEntry: logEntry,
	}
	for _, opt := range opts {
		opt(sm)
	}
	return sm
}

func WithDb(db *mongo.Client, dbName string) StoreManagerOption {
	return func(sm *storeManager) {
		if db != nil && len(dbName) > 0 {
			sm.db = db
			sm.dbName = dbName
		}
	}
}

func WithCache(cache *redis.Client) StoreManagerOption {
	return func(sm *storeManager) {
		if cache != nil {
			sm.cache = cache
		}
	}
}

func (sm *storeManager) TinyUrlStore() datamodel.TinyUrlStore {
	if sm.tinyUrlStore == nil {
		if sm.db == nil || sm.cache == nil || len(sm.dbName) == 0 {
			sm.tinyUrlStore = NewMockTinyUrlStore()
		} else {
			tinyUrlDb := sm.db.Database(sm.dbName).Collection(tinyUrlCollection)
			sm.tinyUrlStore = NewTinyUrlStore(sm.ctx, tinyUrlDb, sm.cache, sm.logEntry)
		}
	}
	return sm.tinyUrlStore
}

func (sm *storeManager) AuthStore() datamodel.AuthStore {
	if sm.authStore == nil {
		sm.authStore = NewAuthStore()
	}
	return sm.authStore
}
