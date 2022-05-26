package dto

import (
	"hoangphuc.tech/hercules/domain/model"
)

type Brand struct {
	model.Brand
}

// Used to validate on creation
type BrandCreated struct {
	Name string `json:"name" validate:"required"`
	Entity
}

func (c *BrandCreated) ToModel() *model.Brand {
	m := &model.Brand{
		Name:   c.Name,
		Entity: *c.Entity.ToModel(),
	}

	return m
}

type BrandUpdated struct {
	Name      string `json:"name"`
	UpdatedBy string `json:"updated_by" validate:"required"`
}

func (c *BrandUpdated) ToModel() *model.Brand {
	m := &model.Brand{
		Name: c.Name,
	}

	return m
}
