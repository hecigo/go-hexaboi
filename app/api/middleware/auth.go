package middleware

import (
	log "github.com/sirupsen/logrus"
	"hoangphuc.tech/go-hexaboi/infra/core"
	"hoangphuc.tech/go-hexaboi/infra/redis"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"

	"hoangphuc.tech/go-hexaboi/infra/jwt"

	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

type AuthService struct{}

type AuthServiceError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type AuthHeader struct {
	APIKey        string `reqHeader:"api-key"`
	Authorization string `reqHeader:"Authorization"`
}

var bytes = core.BytesCrypt

const MySecret = core.MySecretCrypt

func NewAuthService() *AuthService {
	return &AuthService{}
}

func (*AuthService) AuthCheck(c *fiber.Ctx) error {

	a := new(AuthHeader)

	if err := c.ReqHeaderParser(a); err != nil {
		log.Println(fiber.StatusBadRequest, err)
		return response(c, 401, "Permission denied")
	}

	// client to server
	if a.Authorization != "" {
		// validate access token
		var claims *jwt.AWSCognitoClaims
		claims, err := jwt.ValidateJWTToken[jwt.AWSCognitoClaims](a.Authorization)
		if err != nil {
			return response(c, 401, err.Error())
		}

		// get claims; using phone number as username
		session := redis.GetSession(claims.Username)

		// session is found
		if (session != nil) && (session.PhoneNumber == claims.Username) {
			return c.Next()
		}

		return response(c, 401, "session not found")
	}

	// server to server
	stringDecrypt, err := decrypt(c, a.APIKey)
	if err != nil {
		log.Println("error decrypting your classified ---> ", err)
		return response(c, 401, "Permission denied")
	}

	// TODO: dynamic key
	switch stringDecrypt {
	case core.Getenv("HPI_SECRET_TORII_TO_SAWAN", "HPI_SECRET_TORII_TO_SAWAN"):
		return c.Next()
	case core.Getenv("HPI_SECRET_WEB_V2_TO_SAWAN", "HPI_SECRET_WEB_V2_TO_SAWAN"):
		return c.Next()
	default:
		return response(c, 401, "Permission denied")
	}
}

func response(c *fiber.Ctx, status int, msg string) error {
	data := AuthServiceError{
		Status:  status,
		Message: msg,
	}

	if status <= 0 {
		data.Status = c.Response().StatusCode()
	}
	if msg == "" {
		data.Message = utils.StatusMessage(data.Status)
	}

	c.Response().SetStatusCode(status)

	return c.JSON(data)
}

func decode(c *fiber.Ctx, s string) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil, response(c, 401, "Permission denied")
	}
	return data, nil
}

func decrypt(c *fiber.Ctx, text string) (string, error) {
	block, err := aes.NewCipher([]byte(MySecret))
	if err != nil {
		return "", response(c, 401, "Permission denied")
	}

	cipherText, err := decode(c, text)
	if err != nil {
		return "", response(c, 401, "Permission denied")
	}

	cfb := cipher.NewCFBDecrypter(block, bytes)
	plainText := make([]byte, len(cipherText))
	cfb.XORKeyStream(plainText, cipherText)

	return string(plainText), nil

}
