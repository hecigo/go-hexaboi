package handler

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"hoangphuc.tech/hercules/domain/model"
	"hoangphuc.tech/hercules/infra/adapter"
	"hoangphuc.tech/hercules/infra/orm"
)

type ItemHandler struct {
}

var (
	repoItem adapter.ItemRepository = adapter.ItemRepository{}
)

// @Summary Search order item by code.
// @Description Search order item by code.
// @Tags Item
// @Accept */*
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Param code path int true "Code"
// @Router /v1/item/{code} [get]
func (h ItemHandler) Get(c *fiber.Ctx) error {
	code := c.Params("code")

	return c.SendString(fmt.Sprintf("Warehouse code: %s", code))
}

func (h ItemHandler) Create(c *fiber.Ctx) error {
	// 1: Parse from DTO to Model
	// 2: Model.services()
	// 3: Parse from Model > ORM
	// 4: Save ORM

	// Parse payload as domain.Item
	payload := new(model.Item)
	if err := c.BodyParser(payload); err != nil {
		return err
	}

	// Convert payload to orm.Item
	item := orm.NewItem(payload)

	// Create new item into repository
	err := repoItem.Create(item)
	if err != nil {
		return err
	}

	return c.JSON(item)
}

// @Summary Search order items.
// @Description Search order items.
// @Tags Item
// @Accept */*
// @Produce json
// @Success 200 {string} status "ok"
// @Router /v1/item [get]
func (h ItemHandler) Search(c *fiber.Ctx) error {

	return HJSON(c, APIResult{
		Status:  200,
		Message: "Search",
	})
}
