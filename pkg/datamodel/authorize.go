package datamodel

type AuthApplication interface {
	Authorize(string, string) (bool, error)
}

type AuthStore interface {
	Authorize(string, string) (bool, error)
}
