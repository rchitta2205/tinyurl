package datamodel

type TinyUrlApplication interface {
	Create(string) string
	Fetch(string) (string, error)
}

type TinyUrlStore interface {
	Create(string, string)
	Fetch(string) (string, error)
}
