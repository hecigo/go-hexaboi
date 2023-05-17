package router

import (
	"github.com/gofiber/fiber/v2"
	"hecigo.com/go-hexaboi/app/api/dto"
	"hecigo.com/go-hexaboi/app/api/handler"
	"hecigo.com/go-hexaboi/app/api/middleware"
)

type ItemRouter struct{}

func (r ItemRouter) Register(root fiber.Router) {
	group := root.Group("/item")
	handle := handler.NewItemHandler()
	valid := middleware.NewValidator()

	// Special item
	group.Get("/id.:id", valid.Params(&dto.EntityID{}), handle.Get)
	group.Get("/:code", valid.Params(&dto.EntityCode{}), handle.GetByCode)

	// Create new
	group.Post("/", valid.Body(&dto.ItemCreated{}), handle.Create)
	group.Post("/id.:id", valid.Params(&dto.EntityID{}), valid.Body(&dto.ItemUpdated{}), handle.Update)

	// Search
	group.Post("/search/_index", handle.SearchIndex)
	group.Post("/search", handle.Search)
}
