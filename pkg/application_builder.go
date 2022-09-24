package pkg

import (
	"github.com/sirupsen/logrus"
	"tinyurl/pkg/application"
	"tinyurl/pkg/datamodel"
	"tinyurl/pkg/store"
)

type applicationManagerBuilder struct {
	logEntry *logrus.Entry
}

func NewApplicationManagerBuilder(logEntry *logrus.Entry) *applicationManagerBuilder {
	return &applicationManagerBuilder{
		logEntry: logEntry,
	}
}

func (builder *applicationManagerBuilder) Build() (datamodel.ApplicationManager, error) {
	storeMgr := store.NewStoreManager(builder.logEntry)
	appManager := application.NewApplicationManager(builder.logEntry, storeMgr)
	return appManager, nil
}
