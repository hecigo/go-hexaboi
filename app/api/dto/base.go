package dto

import (
	"hoangphuc.tech/hercules/domain/base"
)

type Entity struct {
	CreatedBy string `json:"created_by" validate:"required"`
	UpdatedBy string `json:"updated_by" validate:"required"`
}

// Convert DTO Entity to Base Entity
func (dto *Entity) ToModel() *base.Entity {
	return &base.Entity{
		CreatedBy: dto.CreatedBy,
		UpdatedBy: dto.UpdatedBy,
	}
}

// Used to validate GET request with `id`
type EntityID struct {
	ID uint `json:"id" validate:"gt=0"`
}

// Used to validate GET request with `code`
type EntityCode struct {
	Code string `json:"code" validate:"required,max=50"`
}
