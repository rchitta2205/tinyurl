package store

import (
	"github.com/pkg/errors"
)

// TODO: To be replaced with actual db
type method struct {
	permissions map[string][]string
}

type user struct {
	role map[string]string
}

var (
	methodDb = method{
		permissions: map[string][]string{
			"/TinyUrlService/Create":  {"admin", "user"},
			"/TinyUrlService/Fetch":   {"admin"},
		},
	}

	userDb = user{
		role: map[string]string{
			"alice": "admin",
			"bob":   "user",
		},
	}
)

type authStore struct {
	methodDb method
	userDb   user
}

func NewAuthStore() *authStore {
	db := &authStore{}
	db.userDb = userDb
	db.methodDb = methodDb
	return db
}

func (a *authStore) Authorize(username, method string) (bool, error) {
	role, ok := a.userDb.role[username]
	if !ok {
		return false, errors.New("User: " + username + ", doesn't have a role in database.")
	}

	acceptableRoles, ok := a.methodDb.permissions[method]
	if !ok {
		return false, errors.New("Method: " + method + " is not accessible.")
	}

	for _, currRole := range acceptableRoles {
		if currRole == role{
			return true, nil
		}
	}

	return false, errors.New("User: " + username + ", doesn't have access to make a call to method: " + method)
}
