package middleware

import (
	"reflect"

	"github.com/elliotchance/pie/v2"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"hecigo.com/go-hexaboi/infra/core"
)

type Validator struct {
	validate *validator.Validate
}

type ValidationError struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// Validator constructor
func NewValidator() *Validator {
	return &Validator{
		validate: validator.New(),
	}
}

// Validate body
func (v *Validator) Body(dto interface{}) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return v.validateBody(c, dto)
	}
}

func (v *Validator) validateBody(c *fiber.Ctx, dto interface{}) error {
	if dto == nil || c.Body() == nil {
		return c.Next()
	}

	if err := c.BodyParser(dto); err != nil {
		c.Status(fiber.StatusBadRequest)
		return HError(c, fiber.StatusBadRequest, "CAN_NOT_PARSE_BODY", err)
	}

	err := v.validate.Struct(dto)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return HError(c, fiber.StatusBadRequest, "BODY_INVALID", err, errorDetail(err))
	}

	return c.Next()
}

// Validate path params
func (v *Validator) Params(dto interface{}) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return v.validateParams(c, dto)
	}
}

func (v *Validator) validateParams(c *fiber.Ctx, dto interface{}) error {
	if dto == nil {
		return c.Next()
	}

	fiberParams := c.AllParams()
	if len(fiberParams) == 0 {
		return c.Next()
	}

	// c.AllParams() always return a map[string]string,
	// validate.Struct gets wrong error once a param is the number.
	// We have to convert the string param to a number before validating.

	// Extracts list of numeric field from DTO
	elem := reflect.TypeOf(dto).Elem()
	fieldsLen := elem.NumField()
	var numFields []string
	for i := 0; i < fieldsLen; i++ {
		f := elem.Field(i)
		if core.Utils.IsNumberField(f) {
			jsonFieldName := core.Utils.GetJSONFieldName(f.Tag)
			numFields = append(numFields, jsonFieldName)
		}
	}

	// Convert map[string]string to map[string]interface{}
	params := make(map[string]interface{})
	for k, v := range fiberParams {
		if pie.Contains(numFields, k) {
			vInt, err := core.Utils.ParseUint(v)
			if err != nil {
				return HError(c, fiber.StatusBadRequest, "CAN_NOT_PARSE_PARAMS", err)
			}
			params[k] = vInt
		} else {
			params[k] = v
		}
	}

	// Convert map[string]interface{} to struct
	err := core.Utils.MapToStruct(params, dto)
	if err != nil {
		return HError(c, fiber.StatusBadRequest, "CAN_NOT_PARSE_PARAMS", err)
	}

	err = v.validate.Struct(dto)
	if err != nil {
		return HError(c, fiber.StatusBadRequest, "PATH_PARAMS_INVALID", err, errorDetail(err))
	}

	return c.Next()
}

// Validate query string
func (v *Validator) Query(dto interface{}) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return v.validateQuery(c, dto)
	}
}

func (v *Validator) validateQuery(c *fiber.Ctx, dto interface{}) error {
	if dto == nil {
		return c.Next()
	}

	if err := c.QueryParser(dto); err != nil {
		return HError(c, fiber.StatusBadRequest, "CAN_NOT_PARSE_QUERY_STRING", err)
	}

	err := v.validate.Struct(dto)
	if err != nil {
		return HError(c, fiber.StatusBadRequest, "QUERY_STRING_INVALID", err, errorDetail(err))
	}

	return c.Next()
}

// Validate all route: Params > Query string > Body
func (v *Validator) All(dto ...struct{}) fiber.Handler {
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

func errorDetail(err error) map[string]string {
	var errDetail map[string]string = make(map[string]string)
	for _, err := range err.(validator.ValidationErrors) {
		errDetail["field"] = err.Field()
		errDetail["condition"] = err.ActualTag()
		errDetail["param"] = err.Param()
	}
	return errDetail
}
