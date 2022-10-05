package store

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"tinyurl/pkg/datamodel"
	"tinyurl/pkg/messages"
)

type mockTinyUrlStore struct {
	cache    map[string]string
	logEntry *logrus.Entry
}

func NewMockTinyUrlStore(logEntry *logrus.Entry) datamodel.TinyUrlStore {
	return &mockTinyUrlStore{
		cache:    make(map[string]string),
		logEntry: logEntry,
	}
}

func (m *mockTinyUrlStore) Fetch(tinyUrl string) (string, error) {
	longUrl, exists := m.cache[tinyUrl]
	if !exists {
		return "", errors.New(messages.ErrTinyUrlNotExist)
	}
	m.logEntry.Infof("Fetched long url from mock db %s", longUrl)
	return longUrl, nil
}

func (m *mockTinyUrlStore) Create(url datamodel.Url) error {
	m.cache[url.TinyUrl] = url.LongUrl
	m.logEntry.Infof("Stored tiny url in mock db %s", url.TinyUrl)
	return nil
}
