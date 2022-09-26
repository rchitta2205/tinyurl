package store

import (
	"github.com/pkg/errors"
	"tinyurl/pkg/datamodel"
	"tinyurl/pkg/messages"
)

type mockTinyUrlStore struct {
	cache map[string]string
}

func NewMockTinyUrlStore() datamodel.TinyUrlStore {
	return &mockTinyUrlStore{
		cache: make(map[string]string),
	}
}

func (m *mockTinyUrlStore) Fetch(tinyUrl string) (string, error) {
	longUrl, exists := m.cache[tinyUrl]
	if !exists {
		return "", errors.New(messages.ErrTinyUrlNotExist)
	}
	return longUrl, nil
}

func (m *mockTinyUrlStore) Create(url datamodel.Url) error {
	m.cache[url.TinyUrl] = url.LongUrl
	return nil
}
