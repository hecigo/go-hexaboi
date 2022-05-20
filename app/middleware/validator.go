package middleware

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"hoangphuc.tech/hercules/infra/core"
)

type Validator struct {
	validate *validator.Validate
}

// Validator constructor
func NewValidator() *Validator {
	return &Validator{
		validate: validator.New(),
	}
}

// Validate body
func (v *Validator) ValidateBody(dto interface{}) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return v.validateBody(c, dto)
	}
}

func (v *Validator) validateBody(c *fiber.Ctx, dto interface{}) error {
	if dto == nil || c.Body() == nil {
		return c.Next()
	}

	if err := c.BodyParser(dto); err != nil {
		return fmt.Errorf("CAN_NOT_PARSE_BODY\r\n%v", err)
	}

	err := v.validate.Struct(dto)
	if err != nil {
		return fmt.Errorf("BODY_INVALID\r\n%v", err)
	}

	return c.Next()
}

// Validate path params
func (v *Validator) ValidateParams(dto interface{}) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return v.validateParams(c, dto)
	}
}

func (v *Validator) validateParams(c *fiber.Ctx, dto interface{}) error {
	if dto == nil {
		return c.Next()
	}

	params := c.AllParams()
	if len(params) == 0 {
		return c.Next()
	}

	dto, err := core.Utils{}.MapToStruct(params, dto)
	if err != nil {
		return fmt.Errorf("CAN_NOT_PARSE_PARAMS\r\n%v", err)
	}

	err = v.validate.Struct(dto)
	if err != nil {
		return fmt.Errorf("PATH_PARAMS_INVALID\r\n%v", err)
	}

	return c.Next()
}

// Validate query string
func (v *Validator) ValidateQuery(dto interface{}) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return v.validateQuery(c, dto)
	}
}

func (v *Validator) validateQuery(c *fiber.Ctx, dto interface{}) error {
	if dto == nil {
		return c.Next()
	}

	if err := c.QueryParser(dto); err != nil {
		return fmt.Errorf("CAN_NOT_PARSE_QUERY_STRING\r\n%v", err)
	}

	err := v.validate.Struct(dto)
	if err != nil {
		return fmt.Errorf("QUERY_STRING_INVALID\r\n%v", err)
	}

	return c.Next()
}

// Validate all route: Params > Query string > Body
func (v *Validator) ValidateAll(dto ...interface{}) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		if err := v.validateParams(c, dto); err != nil {
			return err
		}

		if err := v.validateQuery(c, dto); err != nil {
			return err
		}

		if err := v.validateBody(c, dto); err != nil {
			return err
		}

		return c.Next()
	}
}
