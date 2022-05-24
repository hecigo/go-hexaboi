package handler

import (
	"github.com/gofiber/fiber/v2"
	"hoangphuc.tech/hercules/app/api/dto"
	"hoangphuc.tech/hercules/infra/adapter"
	"hoangphuc.tech/hercules/infra/core"
)

type CategoryHandler struct {
	repoCate adapter.CategoryRepository
}

func NewCategoryHandler() *CategoryHandler {
	return &CategoryHandler{
		repoCate: adapter.CategoryRepository{},
	}
}

func (h CategoryHandler) Create(c *fiber.Ctx) error {
	// Parse payload as domain.Item
	dtoCate := new(dto.CategoryCreated)
	if err := c.BodyParser(dtoCate); err != nil {
		return err
	}

	cate := dtoCate.ToModel()

	// Create new item into repository
	err := h.repoCate.Create(cate)
	if err != nil {
		return err
	}

	return HJSON(c, cate)
}

func (h CategoryHandler) Get(c *fiber.Ctx) error {
	id, _ := core.Utils.ParseUint(c.Params("id"))
	cate, err := h.repoCate.GetByID(id)
	if err != nil {
		return err
	}
	return HJSON(c, cate)
}
