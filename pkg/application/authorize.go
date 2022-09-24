package application

import "tinyurl/pkg/datamodel"

type authApplication struct {
	store datamodel.AuthStore
}

func NewAuthApplication(store datamodel.AuthStore) datamodel.AuthApplication {
	return &authApplication{
		store: store,
	}
}

func (a *authApplication) Authorize(username, method string) (bool, error) {
	return a.store.Authorize(username, method)
}
