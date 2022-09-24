package application

import (
	"github.com/sirupsen/logrus"
	"tinyurl/pkg/datamodel"
)

type applicationManager struct {
	logEntry     *logrus.Entry
	storeManager datamodel.StoreManager

	// List of all singleton applications
	tinyUrlApplication datamodel.TinyUrlApplication
	authApplication   datamodel.AuthApplication
}

func NewApplicationManager(logEntry *logrus.Entry, storeManager datamodel.StoreManager) datamodel.ApplicationManager {
	am := &applicationManager{
		logEntry:     logEntry,
		storeManager: storeManager,
	}
	return am
}

func (am *applicationManager) TinyUrlApplication() datamodel.TinyUrlApplication {
	if am.tinyUrlApplication == nil {
		am.tinyUrlApplication = NewTinyUrlApplication(am.storeManager.TinyUrlStore(), am.logEntry)
	}
	return am.tinyUrlApplication
}

func (am *applicationManager) AuthApplication() datamodel.AuthApplication {
	if am.authApplication == nil {
		am.authApplication = NewAuthApplication(am.storeManager.AuthStore())
	}
	return am.authApplication
}
