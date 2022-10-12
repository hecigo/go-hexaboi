package middleware

import (
	log "github.com/sirupsen/logrus"
	"hoangphuc.tech/go-hexaboi/infra/core"
	"hoangphuc.tech/go-hexaboi/infra/jwt"
	"hoangphuc.tech/go-hexaboi/infra/redis"

	"github.com/gofiber/fiber/v2"

	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

type Auth struct{}

type AuthServiceError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type AuthHeader struct {
	APIKey        string `reqHeader:"api-key"`
	Authorization string `reqHeader:"Authorization"`
}

func (_auth *Auth) Enable(app *fiber.App) error {
	app.Use(_auth.AuthCheck)
	return nil
}

func (*Auth) AuthCheck(c *fiber.Ctx) error {

	a := new(AuthHeader)

	if err := c.ReqHeaderParser(a); err != nil {
		log.Error("parse Auth header --> ", err)
		return HError(c, fiber.StatusUnauthorized, "UNAUTHORIZED", nil)
	}

	// client to server
	if a.Authorization != "" {
		// validate access token
		var claims *jwt.AWSCognitoClaims
		claims, err := jwt.ValidateJWTToken[jwt.AWSCognitoClaims](a.Authorization)
		if err != nil {
			log.Error("claim JWT --> ", err)
			return HError(c, fiber.StatusUnauthorized, "UNAUTHORIZED", nil)
		}

		// get claims; using phone number as username
		session := redis.GetSession(claims.Username, a.Authorization)

		// session is found
		if (session != nil) && (session.UserID == claims.Username) {
			return c.Next()
		}

		log.Error("no session --> ", err)
		return HError(c, fiber.StatusUnauthorized, "UNAUTHORIZED", nil)
	}

	// server to server
	stringDecrypt, err := decrypt(c, a.APIKey)
	if err != nil {
		log.Println("error decrypting your classified --> ", err)
		// Should NOT expose error message to client
		return HError(c, fiber.StatusUnauthorized, "UNAUTHORIZED", nil)
	}

	// Check Client API Key
	if core.API_CLIENT_SECRETS != nil {
		for _, v := range core.API_CLIENT_SECRETS {
			if stringDecrypt == v {
				return c.Next()
			}
		}
	}

	return HError(c, fiber.StatusUnauthorized, "UNAUTHORIZED", nil)
}

func decode(c *fiber.Ctx, s string) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func decrypt(c *fiber.Ctx, text string) (string, error) {
	block, err := aes.NewCipher([]byte(core.MySecretCrypt))
	if err != nil {
		return "", err
	}

	cipherText, err := decode(c, text)
	if err != nil {
		return "", err
	}

	cfb := cipher.NewCFBDecrypter(block, core.BytesCrypt)
	plainText := make([]byte, len(cipherText))
	cfb.XORKeyStream(plainText, cipherText)

	return string(plainText), nil

}
