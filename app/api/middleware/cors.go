package middleware

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/hecigo/goutils"
)

type CORS struct {
	AllowOrigins string
	AllowHeaders string
	AllowMethods string
}

func (_cors *CORS) Enable(app *fiber.App) error {

	_cors.AllowOrigins = goutils.Env("CORS_ALLOW_ORIGINS", "localhost")
	_cors.AllowHeaders = goutils.Env("CORS_ALLOW_HEADERS", "")
	_cors.AllowMethods = goutils.Env("CORS_ALLOW_METHODS", "GET,POST,HEAD,PUT,PATCH")

	app.Use(cors.New(cors.Config{
		AllowOrigins: _cors.AllowOrigins,
		AllowHeaders: _cors.AllowHeaders,
		AllowMethods: _cors.AllowMethods,
	}))

	_cors.Print()

	return nil
}

func (_cors CORS) Print() {
	fmt.Println("\r\n┌─────── Middleware/CORS ────────")
	fmt.Printf("| CORS_ALLOW_ORIGINS: %s\r\n", _cors.AllowOrigins)
	fmt.Printf("| CORS_ALLOW_HEADERS: %s\r\n", _cors.AllowHeaders)
	fmt.Printf("| CORS_ALLOW_METHODS: %s\r\n", _cors.AllowMethods)
	fmt.Println("└────────────────────────────────")
}
