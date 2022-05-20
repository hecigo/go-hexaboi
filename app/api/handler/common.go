package handler

import (
	"github.com/gofiber/fiber/v2"
)

// NotFound returns custom 404 page
func NotFound(c *fiber.Ctx) error {
	return c.SendStatus(404)
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
