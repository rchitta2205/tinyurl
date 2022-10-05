package datamodel

type Url struct {
	TinyUrl string
	LongUrl string
}

type TinyUrlApplication interface {
	Create(string) (string, error)
	Fetch(string) (string, error)
}

type TinyUrlStore interface {
	Create(Url) error
	Fetch(string) (string, error)
}
