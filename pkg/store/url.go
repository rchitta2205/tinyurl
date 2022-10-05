package store

import (
	"context"
	dapr "github.com/dapr/go-sdk/client"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"tinyurl/pkg/datamodel"
	"tinyurl/pkg/messages"
)

type tinyUrlStore struct {
	ctx        context.Context
	daprClient dapr.Client
	logEntry   *logrus.Entry
	db         string
	cache      string
}

func NewTinyUrlStore(ctx context.Context, daprClient dapr.Client,
	logEntry *logrus.Entry, db, cache string) datamodel.TinyUrlStore {
	return &tinyUrlStore{
		ctx:        ctx,
		daprClient: daprClient,
		logEntry:   logEntry,
		db:         db,
		cache:      cache,
	}
}

func (t *tinyUrlStore) Fetch(tinyUrl string) (string, error) {
	var longUrl string

	// Check in cache to see if the tiny url exists
	item, err := t.daprClient.GetState(t.ctx, t.cache, tinyUrl, nil)
	if err == nil {
		// The key exists in the cache so we just return it
		longUrl = string(item.Value)
		if len(longUrl) > 0 {
			t.logEntry.Infof("Fetched url from cache %s", longUrl)
			return longUrl, nil
		}
	}

	// Check in db to see if the tiny url exists, since it doesn't
	// exist in cache.
	item, err = t.daprClient.GetState(t.ctx, t.db, tinyUrl, nil)
	if err != nil {
		// If the url doesn't exist even in db then we return error
		return "", errors.Wrapf(err, messages.ErrTinyUrlNotExist)
	}

	longUrl = string(item.Value)
	if len(longUrl) == 0 {
		// If the url is empty then we return an error
		return "", errors.New(messages.ErrTinyUrlNotExist)
	}

	// Setting it in the cache for later retrievals
	err = t.daprClient.SaveState(t.ctx, t.cache, tinyUrl, item.Value, nil)
	if err != nil {
		return "", errors.WithStack(err)
	}

	t.logEntry.Infof("Fetched url from db %s", longUrl)
	return longUrl, nil
}

func (t *tinyUrlStore) Create(url datamodel.Url) error {
	err := t.daprClient.SaveState(t.ctx, t.db, url.TinyUrl, []byte(url.LongUrl), nil)
	if err != nil {
		return errors.WithStack(err)
	}
	t.logEntry.Infof("Created tinyurl %s in db", url.TinyUrl)
	return nil
}
