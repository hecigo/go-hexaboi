package middleware

import (
	"github.com/hecigo/goutils"
	log "github.com/sirupsen/logrus"
	"hecigo.com/go-hexaboi/infra/jwt"
	"hecigo.com/go-hexaboi/infra/redis"

	"github.com/gofiber/fiber/v2"
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
		session := redis.GetSession(c.Context(), claims.Username, a.Authorization)

		// session is found
		if (session != nil) && (session.UserID == claims.Username) {
			return c.Next()
		}

		log.Error("no session --> ", err)
		return HError(c, fiber.StatusUnauthorized, "UNAUTHORIZED", nil)
	}

	err := goutils.CheckAPISecretKey(a.APIKey)
	if err != nil {
		log.Println(err)
		return HError(c, fiber.StatusUnauthorized, "UNAUTHORIZED", nil)
	}

	return c.Next()
}
