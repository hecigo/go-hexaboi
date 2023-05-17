package router

import (
	"github.com/gofiber/fiber/v2"
	"hecigo.com/go-hexaboi/app/api/dto"
	"hecigo.com/go-hexaboi/app/api/handler"
	"hecigo.com/go-hexaboi/app/api/middleware"
)

type BrandRouter struct{}

func (r BrandRouter) Register(root fiber.Router) {
	group := root.Group("/brand")
	handle := handler.NewBrandHandler()
	valid := middleware.NewValidator()

	// GET
	group.Get("/id.:id", valid.Params(&dto.EntityID{}), handle.Get)
	group.Get("/:code", valid.Params(&dto.EntityCode{}), handle.GetByCode)

	// POST
	group.Post("/", valid.Body(&dto.BrandCreated{}), handle.Create)
	group.Post("/id.:id", valid.Params(&dto.EntityID{}), valid.Body(&dto.BrandUpdated{}), handle.Update)
}
