package core

type Cache interface {
	Set(key string, value interface{}) error
	Get(key string) interface{}
}
