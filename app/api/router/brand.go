package router

import (
	"github.com/gofiber/fiber/v2"
	"hoangphuc.tech/hercules/app/api/dto"
	"hoangphuc.tech/hercules/app/api/handler"
	"hoangphuc.tech/hercules/app/middleware"
)

type BrandRouter struct{}

func (r BrandRouter) Register(root fiber.Router) {
	group := root.Group("/brand")
	handle := handler.NewBrandHandler()
	valid := middleware.NewValidator()

	// GET
	group.Get("/:id", valid.Params(&dto.EntityID{}), handle.Get)

	// POST
	group.Post("/", valid.Body(&dto.BrandCreated{}), handle.Create)
	group.Post("/:id", valid.Params(&dto.EntityID{}), valid.Body(&dto.BrandUpdated{}), handle.Update)
}
