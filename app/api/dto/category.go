package dto

import (
	"hoangphuc.tech/go-hexaboi/domain/model"
)

type Category struct {
	model.Category
}

// Used to validate on creation
type CategoryCreated struct {
	Code       string         `json:"code" validate:"required,max=50"`
	Name       string         `json:"name" validate:"required"`
	DivisionBy model.Division `json:"division_by" validate:"required,oneof=category campaign custom"`
	ParentID   uint           `json:"parent_id"`
	Entity
}

func (cc *CategoryCreated) ToModel() *model.Category {
	mc := &model.Category{
		Code: cc.Code,
		Name: cc.Name,
	}

	if cc.ParentID > 0 {

	}

	return mc
}

type CategoryUpdated struct {
	Name       string         `json:"name"`
	DivisionBy model.Division `json:"division_by" validate:"oneof=merchandise consumer campaign custom"`
	ParentID   *uint          `json:"parent_id"`
	UpdatedBy  string         `json:"updated_by" validate:"required"`
}

func (ic *CategoryUpdated) ToModel() *model.Category {
	m := &model.Category{
		Name: ic.Name,
	}

	if ic.ParentID != nil {

	}

	return m
}
