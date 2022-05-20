package handler

import "github.com/gofiber/fiber/v2"

type APIResult struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func HJSON(c *fiber.Ctx, data APIResult) error {
	return c.Status(data.Status).JSON(data)
}
