package router

import (
	"github.com/gofiber/fiber/v2"
	"hoangphuc.tech/hercules/app/api/dto"
	"hoangphuc.tech/hercules/app/api/handler"
	"hoangphuc.tech/hercules/app/middleware"
)

type ItemRouter struct{}

func (r ItemRouter) Register(root fiber.Router) {
	group := root.Group("/item")
	handle := handler.NewItemHandler()
	valid := middleware.NewValidator()

	// GET
	group.Get("/id.:id", valid.Params(&dto.EntityID{}), handle.Get)
	group.Get("/:code", valid.Params(&dto.EntityCode{}), handle.GetByCode)

	// POST
	group.Post("/", valid.Body(&dto.ItemCreated{}), handle.Create)
}
