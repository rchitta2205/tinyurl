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

func (w *tinyUrlApplication) Create(longUrl string) string {
	var tinyUrl string
	for {
		tinyUrl = w.generateKey()
		_, err := w.store.Fetch(tinyUrl)
		if err == nil {
			break
		}
	}
	w.store.Create(tinyUrl, longUrl)
	return tinyUrl
}

func (w *tinyUrlApplication) Fetch(tinyUrl string) (string, error) {
	_, err := neturl.ParseRequestURI(tinyUrl)
	if err != nil {
		return "", errors.Wrapf(err, messages.ErrInvalidUrl, tinyUrl)
	}

	longUrl, err := w.store.Fetch(tinyUrl)
	if err != nil {
		return "", errors.WithStack(err)
	}

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
