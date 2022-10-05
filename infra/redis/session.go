package redis

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"hoangphuc.tech/go-hexaboi/domain/base"
	"hoangphuc.tech/go-hexaboi/infra/core"
)

type Session struct {
	base.EntityID
	UserID        string `json:"userId"`
	UserFullname  string `json:"userFullName"`
	UserEmail     string `json:"userEmail"`
	UserPassword  string `json:"userPassword"`
	UserLastLogin string `json:"userLastLogin"`
	PhoneNumber   string `json:"phoneNumber"`
	IsActivated   int    `json:"isActivated"`
	IsDeleted     int    `json:"isDeleted"`
	base.Entity
}

var SESSION_KEY_FORMAT string

func EnableSession() {
	SESSION_KEY_FORMAT = core.Getenv("REDIS_SESSION_KEY_FORMAT", "session:%s")
}

func GetSession(uuid string) *Session {
	if SESSION_KEY_FORMAT == "" {
		log.Errorln("session key format is empty, let call EnableSession() first to load key format from environemnt variable")
		return nil
	}

	session, err := GetSpecificKey[Session](fmt.Sprintf(SESSION_KEY_FORMAT, uuid))
	if err != nil {
		log.Errorln(err)
		return nil
	}
	return session
}
