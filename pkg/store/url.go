package store

import (
	"context"
	"github.com/go-redis/redis"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"tinyurl/pkg/datamodel"
	"tinyurl/pkg/messages"
	"tinyurl/pkg/util"
)

const (
	tinyUrlCollection = "tiny"
	tinyUrlFieldName  = "tiny_url"
	NoExpiration      = 0
)

type tinyUrlStore struct {
	ctx      context.Context
	db       *mongo.Collection
	cache    *redis.Client
	logEntry *logrus.Entry
}

func NewTinyUrlStore(ctx context.Context, db *mongo.Collection, cache *redis.Client, logEntry *logrus.Entry) datamodel.TinyUrlStore {
	err := util.CreateIndex(ctx, db, tinyUrlFieldName)
	if err != nil {
		logEntry.Fatalf(err.Error())
	}
	return &tinyUrlStore{
		ctx:      ctx,
		db:       db,
		cache:    cache,
		logEntry: logEntry,
	}
}

func (t *tinyUrlStore) Fetch(tinyUrl string) (string, error) {
	var longUrl string
	var err error

	// Check in cache to see if the tiny url exists
	longUrl, err = t.cache.Get(tinyUrl).Result()
	if err == nil {
		// The key exists in the cache so we just return it
		t.logEntry.Infof("Fetched url from cache %s", longUrl)
		return longUrl, nil
	}

	// Check in db to see if the tiny url exists, since it doesn't
	// exist in cache.
	var url datamodel.Url
	var query = bson.D{{tinyUrlFieldName, tinyUrl}}
	err = t.db.FindOne(t.ctx, query).Decode(&url)
	if err != nil {
		// If the url doesn't exist even in db then we return error
		return "", errors.Wrapf(err, messages.ErrTinyUrlNotExist)
	}

	// Setting it in the cache for later retrievals
	longUrl = url.LongUrl
	err = t.cache.Set(tinyUrl, longUrl, NoExpiration).Err()
	if err != nil {
		return "", errors.WithStack(err)
	}

	t.logEntry.Infof("Fetched url from db %s", longUrl)
	return longUrl, nil
}

func (t *tinyUrlStore) Create(url datamodel.Url) error {
	_, err := t.db.InsertOne(t.ctx, url)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}
