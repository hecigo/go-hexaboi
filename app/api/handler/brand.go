package handler

import (
	"github.com/gofiber/fiber/v2"
	"hoangphuc.tech/go-hexaboi/app/api/dto"
	"hoangphuc.tech/go-hexaboi/infra/adapter"
	"hoangphuc.tech/go-hexaboi/infra/core"
)

type BrandHandler struct {
	repoBrand adapter.BrandRepository
}

func NewBrandHandler() *BrandHandler {
	return &BrandHandler{
		repoBrand: adapter.BrandRepository{},
	}
}

func (h BrandHandler) GetByCode(c *fiber.Ctx) error {
	code := c.Params("code")
	brand, err := h.repoBrand.GetByCode(code)
	if err != nil {
		return err
	}
	return HJSON(c, brand)
}

func (h BrandHandler) Get(c *fiber.Ctx) error {
	id, _ := core.Utils.ParseUint(c.Params("id"))
	item, err := h.repoBrand.GetByID(id)
	if err != nil {
		return err
	}
	return HJSON(c, item)
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

func (h BrandHandler) Update(c *fiber.Ctx) error {
	id, _ := core.Utils.ParseUint(c.Params("id"))

	// Parse payload as domain
	d := new(dto.BrandUpdated)
	if err := c.BodyParser(d); err != nil {
		return err
	}
	m := d.ToModel()
	err := h.repoBrand.Update(id, m)
	if err != nil {
		return err
	}
	return HJSON(c, m)
}
