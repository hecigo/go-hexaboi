package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/hecigo/goutils"
)

// Response anything
func HJSON(c *fiber.Ctx, data interface{}) error {
	if hpiResult, ok := data.(goutils.APIRes); ok {
		return hJSON(c, hpiResult)
	}
	status := c.Response().StatusCode()
	return hJSON(c, goutils.APIRes{
		Status:  status,
		Data:    data,
		Message: utils.StatusMessage(status),
	})
}

func hJSON(c *fiber.Ctx, data goutils.APIRes) error {
	if data.Status <= 0 {
		data.Status = c.Response().StatusCode()
	}

	if data.Message == "" {
		data.Message = utils.StatusMessage(data.Status)
	}

	return c.JSON(data)
}

// Centrialize errors
func HError(c *fiber.Ctx, status int, errCode string, err error, data ...interface{}) error {
	// Status code defaults to 500
	code := status
	if code < 400 {
		code = c.Response().StatusCode()
		if code < 400 {
			code = fiber.StatusInternalServerError
		}
	}

	if data == nil {
		data = make([]interface{}, 0)
	}

	message := ""

	// Retrieve the custom status code if it's an *goutils.APIRes
	if err == nil {
		c.Response().SetStatusCode(code)
	} else {
		if e, ok := err.(*goutils.APIRes); ok {
			code = e.Status
			data = append(data, e.Data)
			if e.ErrorCode != "" {
				errCode = e.ErrorCode
			}
			c.Response().SetStatusCode(code) // Override status code of the context with HPIResult
		} else if e, ok := err.(*fiber.Error); ok {
			// Retrieve the custom status code if it's an fiber.*Error
			code = e.Code
		}

		message = err.Error()
	}

	return HJSON(c, goutils.APIRes{
		Status:    code,
		Message:   message,
		ErrorCode: errCode,
		Data:      data,
	})
}
