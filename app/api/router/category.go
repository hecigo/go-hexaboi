package router

import (
	"github.com/gofiber/fiber/v2"
	"hoangphuc.tech/hercules/app/api/dto"
	"hoangphuc.tech/hercules/app/api/handler"
	"hoangphuc.tech/hercules/app/api/middleware"
)

type CategoryRouter struct{}

func (r CategoryRouter) Register(root fiber.Router) {
	group := root.Group("/category")
	handle := handler.NewCategoryHandler()
	valid := middleware.NewValidator()

	// GET
	group.Get("/id.:id", valid.Params(&dto.EntityID{}), handle.Get)
	group.Get("/:code", valid.Params(&dto.EntityCode{}), handle.GetByCode)

	// POST
	group.Post("/", valid.Body(&dto.CategoryCreated{}), handle.Create)
	group.Post("/id.:id", valid.Params(&dto.EntityID{}), valid.Body(&dto.CategoryUpdated{}), handle.Update)
}
