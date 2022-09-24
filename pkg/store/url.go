package store

import (
	"github.com/pkg/errors"
	"tinyurl/pkg/messages"
)

type tinyUrlStore struct {
	cache map[string]string
}

func NewTinyUrlStore() *tinyUrlStore {
	return &tinyUrlStore{
		cache: make(map[string]string),
	}
}

func (t *tinyUrlStore) Fetch(tinyUrl string) (string, error) {
	longUrl, exists := t.cache[tinyUrl]
	if !exists {
		return "", errors.New(messages.ErrTinyUrlNotExist)
	}
	return longUrl, nil
}

func (t *tinyUrlStore) Create(tinyUrl string, longUrl string) {
	t.cache[tinyUrl] = longUrl
}


