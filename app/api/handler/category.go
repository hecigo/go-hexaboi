package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hecigo/goutils"
	"hecigo.com/go-hexaboi/app/api/dto"
	"hecigo.com/go-hexaboi/infra/adapter"
)

type CategoryHandler struct {
	repoCate adapter.CategoryRepository
}

func NewCategoryHandler() *CategoryHandler {
	return &CategoryHandler{
		repoCate: adapter.CategoryRepository{},
	}
}

func (h CategoryHandler) GetByCode(c *fiber.Ctx) error {
	code := c.Params("code")
	cate, err := h.repoCate.GetByCode(code)
	if err != nil {
		return err
	}
	return HJSON(c, cate)
}

func (h CategoryHandler) Get(c *fiber.Ctx) error {
	id, _ := goutils.StrConv[uint](c.Params("id"))
	cate, err := h.repoCate.GetByID(id)
	if err != nil {
		return err
	}
	return HJSON(c, cate)
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

func (h CategoryHandler) Update(c *fiber.Ctx) error {
	id, _ := goutils.StrConv[uint](c.Params("id"))

	// Parse payload as domain
	d := new(dto.CategoryUpdated)
	if err := c.BodyParser(d); err != nil {
		return err
	}
	m := d.ToModel()
	err := h.repoCate.Update(id, m)
	if err != nil {
		return err
	}
	return HJSON(c, m)
}
