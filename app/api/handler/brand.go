package handler

import (
	"github.com/gofiber/fiber/v2"
	"hoangphuc.tech/hercules/app/api/dto"
	"hoangphuc.tech/hercules/infra/adapter"
	"hoangphuc.tech/hercules/infra/core"
)

type BrandHandler struct {
	repoBrand adapter.BrandRepository
}

func NewBrandHandler() *BrandHandler {
	return &BrandHandler{
		repoBrand: adapter.BrandRepository{},
	}
}

func (h BrandHandler) Create(c *fiber.Ctx) error {
	// Parse payload as domain.Item
	d := new(dto.BrandCreated)
	if err := c.BodyParser(d); err != nil {
		return err
	}

	m := d.ToModel()

	// Create new item into repository
	err := h.repoBrand.Create(m)
	if err != nil {
		return err
	}

	return HJSON(c, m)
}

func (h BrandHandler) Get(c *fiber.Ctx) error {
	id, _ := core.Utils.ParseUint(c.Params("id"))
	m, err := h.repoBrand.GetByID(id)
	if err != nil {
		return err
	}
	return HJSON(c, m)
}
