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
	itemHandler := handler.ItemHandler{}
	validator := middleware.NewValidator()

	// GET
	group.Get("/", validator.ValidateBody(&dto.Item1{}), itemHandler.Search)
	group.Get("/:code", validator.ValidateParams(&dto.Item2{}), itemHandler.Search)

	// POST
	group.Post("/", itemHandler.Create)
}
