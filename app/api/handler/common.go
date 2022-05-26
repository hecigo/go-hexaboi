package handler

import (
	"github.com/gofiber/fiber/v2"

	"hoangphuc.tech/hercules/app/middleware"
	ext "hoangphuc.tech/hercules/infra/orm/ext"
)

func HJSON(c *fiber.Ctx, data interface{}) error {
	return middleware.HJSON(c, data)
}

// NotFound response
func NotFound(c *fiber.Ctx) error {
	return c.SendStatus(404)
}

// Error response
func DefaultError(c *fiber.Ctx, err error) error {
	if err, ok := ext.Errorf(err); ok {
		return middleware.HError(c, c.Response().StatusCode(), "ORM_ERROR", err)
	}
	return middleware.HError(c, c.Response().StatusCode(), "HPI_ERROR", err)
}

func Error(c *fiber.Ctx, status int, errCode string, err error, data ...interface{}) error {
	return middleware.HError(c, status, errCode, err, data)
}

// @Summary Show the status of server.
// @Description Get the status of server.
// @Tags Root
// @Accept */*
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router / [get]
func HealthCheck(c *fiber.Ctx) error {
	return c.SendStatus(200)
}
