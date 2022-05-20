package middleware

import "github.com/gofiber/fiber/v2"

type Middleware interface {
	Enable(app *fiber.App) error
	Print()
}
