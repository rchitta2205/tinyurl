package store

import (
	"context"
	dapr "github.com/dapr/go-sdk/client"
	"github.com/sirupsen/logrus"
	"tinyurl/pkg/datamodel"
)

type storeManager struct {
	ctx        context.Context
	logEntry   *logrus.Entry
	daprClient dapr.Client
	db         string
	cache      string

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

func WithDaprClient(daprClient dapr.Client) StoreManagerOption {
	return func(sm *storeManager) {
		if daprClient != nil {
			sm.daprClient = daprClient
		}
	}
}

func WithDb(db string) StoreManagerOption {
	return func(sm *storeManager) {
		if len(db) != 0 {
			sm.db = db
		}
	}
}

func WithCache(cache string) StoreManagerOption {
	return func(sm *storeManager) {
		if len(cache) != 0 {
			sm.cache = cache
		}
	}
}

func (sm *storeManager) TinyUrlStore() datamodel.TinyUrlStore {
	if sm.tinyUrlStore == nil {
		if sm.daprClient == nil || len(sm.db) == 0 || len(sm.cache) == 0 {
			sm.tinyUrlStore = NewMockTinyUrlStore(sm.logEntry)
		} else {
			sm.tinyUrlStore = NewTinyUrlStore(sm.ctx, sm.daprClient, sm.logEntry, sm.db, sm.cache)
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
