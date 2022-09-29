package application

import (
	"crypto/md5"
	"github.com/keys-pub/keys/encoding"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	neturl "net/url"
	"strings"
	"tinyurl/pkg/datamodel"
	"tinyurl/pkg/messages"
)

const (
	prefixUrl = "http://www.tinyurl/"
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

	// Convert the longUrl to md5 hash, base62 encode the hash, and use the 1st
	// 15 characters as the tinyUrl for the longUrl
	hash := md5.Sum([]byte(longUrl))
	tinyUrl := encoding.EncodeBase62(hash[:])[:15]
	url := datamodel.Url{
		TinyUrl: tinyUrl,
		LongUrl: longUrl,
	}

	err = w.store.Create(url)
	if err != nil {
		return "", errors.WithStack(err)
	}

	w.logEntry.Infof("Generated tiny url %s", tinyUrl)
	return prefixUrl + tinyUrl, nil
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

	tinyUrl = strings.Replace(tinyUrl, prefixUrl, "", 1)
	longUrl, err := w.store.Fetch(tinyUrl)
	if err != nil {
		return "", errors.WithStack(err)
	}

	w.logEntry.Infof("Fetched long url %s", longUrl)
	return longUrl, nil
}
