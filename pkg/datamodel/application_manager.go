package datamodel

type ApplicationManager interface {
	TinyUrlApplication() TinyUrlApplication
	AuthApplication() AuthApplication
}
