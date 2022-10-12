package jwt

import "github.com/golang-jwt/jwt/v4"

type AWSCognitoClaims struct {
	ClientID  string `json:"client_id"`
	OriginJTI string `json:"origin_jti"`
	EventID   string `json:"event_id"`
	TokenUse  string `json:"token_use"`
	Scope     string `json:"scope"`
	AuthTime  int64  `json:"auth_time"`
	Username  string `json:"username"`
	jwt.RegisteredClaims
}
