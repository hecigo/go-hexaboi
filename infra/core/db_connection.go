package core

import "gorm.io/gorm"

type DBConnection interface {
	DB() *gorm.DB
	DBByName(name string) *gorm.DB
	OpenDefaultConnection() error
	OpenConnectionByName(connName string) error
	Close(db *gorm.DB) error
	CloseAll() error
}
