package datamodel

type StoreManager interface {
	TinyUrlStore() TinyUrlStore
	AuthStore() AuthStore
}
