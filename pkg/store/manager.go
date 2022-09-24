package store

import (
	"github.com/sirupsen/logrus"
	"tinyurl/pkg/datamodel"
)

type storeManager struct {
	logEntry *logrus.Entry

	// List of all singleton stores
	tinyUrlStore datamodel.TinyUrlStore
	authStore   datamodel.AuthStore
}

func NewStoreManager(logEntry *logrus.Entry) datamodel.StoreManager {
	return &storeManager{
		logEntry: logEntry,
	}
}

func (sm *storeManager) TinyUrlStore() datamodel.TinyUrlStore {
	if sm.tinyUrlStore == nil {
		sm.tinyUrlStore = NewTinyUrlStore()
	}
	return sm.tinyUrlStore
}

func (sm *storeManager) AuthStore() datamodel.AuthStore {
	if sm.authStore == nil {
		sm.authStore = NewAuthStore()
	}
	return sm.authStore
}
