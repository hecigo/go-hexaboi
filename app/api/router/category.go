package router

import (
	"github.com/gofiber/fiber/v2"
	"hoangphuc.tech/hercules/app/api/dto"
	"hoangphuc.tech/hercules/app/api/handler"
	"hoangphuc.tech/hercules/app/middleware"
)

type CategoryRouter struct{}

func (r CategoryRouter) Register(root fiber.Router) {
	group := root.Group("/category")
	handle := handler.NewCategoryHandler()
	valid := middleware.NewValidator()

	// POST
	group.Post("/", valid.Body(&dto.CategoryCreated{}), handle.Create)
	group.Get("/:id", valid.Params(&dto.EntityID{}), handle.Get)
}
