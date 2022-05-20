package router

import (
	"github.com/gofiber/fiber/v2"

	swagger "github.com/arsmn/fiber-swagger/v2"
)

type SwaggerRouter struct{}

func (r SwaggerRouter) Register(root fiber.Router) {
	group := root.Group("/swagger")
	group.Get("*", swagger.HandlerDefault)
}
