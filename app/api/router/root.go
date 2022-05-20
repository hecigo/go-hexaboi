package router

import (
	"github.com/gofiber/fiber/v2"
)

type Router interface {
	Register(fiber.Router) error
}

func RegisterDefaultRouter(router fiber.Router) error {
	(ItemRouter{}).Register(router)
	(SwaggerRouter{}).Register(router)

	return nil
}
