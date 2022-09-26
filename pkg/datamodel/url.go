package datamodel

type Url struct {
	TinyUrl string `json:"tiny_url" bson:"tiny_url"`
	LongUrl string `json:"long_url" bson:"long_url"`
}

type TinyUrlApplication interface {
	Create(string) (string, error)
	Fetch(string) (string, error)
}

type TinyUrlStore interface {
	Create(Url) error
	Fetch(string) (string, error)
}
