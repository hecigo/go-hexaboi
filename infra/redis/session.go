package redis

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"

	"github.com/hecigo/goredis"
	"github.com/hecigo/goutils"
	log "github.com/sirupsen/logrus"
	"hecigo.com/go-hexaboi/domain/base"
)

type Session struct {
	base.EntityID
	UserID        string `json:"userId"`
	UserFullname  string `json:"userFullName"`
	UserEmail     string `json:"userEmail"`
	UserLastLogin string `json:"userLastLogin"`
	PhoneNumber   string `json:"phoneNumber"`
	IsActivated   int    `json:"isActivated"`
	IsDeleted     int    `json:"isDeleted"`
	base.Entity
}

var sessionKeyFormat, randHashKey string

func EnableSession() {
	sessionKeyFormat = goutils.Env("REDIS_SESSION_KEY_FORMAT", "gohexaboi.session:%s/%s")
	randHashKey = goutils.Env("REDIS_SESSION_HASH_KEY", "")
}

// Get session by user info
func GetSession(ctx context.Context, uuid string, args ...string) *Session {
	if sessionKeyFormat == "" {
		log.Errorln("session key format is empty, let call EnableSession() first to load key format from environemnt variable")
		return nil
	}

	// hash session key
	sessionKey := fmt.Sprintf(sessionKeyFormat, uuid, args) // {uuid}/{args}
	hasher := md5.New()
	hasher.Write([]byte(sessionKey + randHashKey))
	hashed := hex.EncodeToString(hasher.Sum(nil))
	sessionKey = fmt.Sprintf(sessionKeyFormat, uuid, hashed) // {uuid}/{hashed({uuid}/{args})}

	val, err := goredis.Get[Session](ctx, sessionKey)
	if err != nil {
		log.Errorln(err)
		return nil
	}
	session := val.(Session)

	return &session
}
