package application

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"math/rand"
	neturl "net/url"
	"tinyurl/pkg/datamodel"
	"tinyurl/pkg/messages"
)

const (
	keyCombinations = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	prefixUrl       = "http://www.tinyurl/"
)

type tinyUrlApplication struct {
	store    datamodel.TinyUrlStore
	logEntry *logrus.Entry
}

func NewTinyUrlApplication(store datamodel.TinyUrlStore, logEntry *logrus.Entry) datamodel.TinyUrlApplication {
	return &tinyUrlApplication{
		store:    store,
		logEntry: logEntry,
	}
}

func (w *tinyUrlApplication) Create(longUrl string) (string, error) {
	longUrl, err := neturl.PathUnescape(longUrl)
	if err != nil {
		return "", errors.Wrapf(err, messages.ErrInvalidUrl, longUrl)
	}

	_, err = neturl.ParseRequestURI(longUrl)
	if err != nil {
		return "", errors.Wrapf(err, messages.ErrInvalidUrl, longUrl)
	}

	var tinyUrl string
	for {
		tinyUrl = w.generateKey()
		_, terr := w.store.Fetch(tinyUrl)
		if terr != nil {
			break
		}
	}

	tinyUrl = prefixUrl + tinyUrl
	url := datamodel.Url{
		TinyUrl: tinyUrl,
		LongUrl: longUrl,
	}

	err = w.store.Create(url)
	if err != nil {
		return "", errors.WithStack(err)
	}

	w.logEntry.Infof("Generated tiny url %s", tinyUrl)
	return tinyUrl, nil
}

func (w *tinyUrlApplication) Fetch(tinyUrl string) (string, error) {
	tinyUrl, err := neturl.PathUnescape(tinyUrl)
	if err != nil {
		return "", errors.Wrapf(err, messages.ErrInvalidUrl, tinyUrl)
	}

	_, err = neturl.ParseRequestURI(tinyUrl)
	if err != nil {
		return "", errors.Wrapf(err, messages.ErrInvalidUrl, tinyUrl)
	}

	longUrl, err := w.store.Fetch(tinyUrl)
	if err != nil {
		return "", errors.WithStack(err)
	}

	w.logEntry.Infof("Fetched long url %s", longUrl)
	return longUrl, nil
}

func (w *tinyUrlApplication) generateKey() string {
	var key []rune
	for i := 0; i < 8; i++ {
		randomIndex := rand.Intn(len(keyCombinations))
		key = append(key, rune(keyCombinations[randomIndex]))
	}
	return string(key)
}
