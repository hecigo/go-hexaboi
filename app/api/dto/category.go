package dto

import (
	"hecigo.com/go-hexaboi/domain/base"
	"hecigo.com/go-hexaboi/domain/model"
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
		Code:       cc.Code,
		Name:       cc.Name,
		DivisionBy: cc.DivisionBy,
		Entity:     *cc.Entity.ToModel(),
	}

	if cc.ParentID > 0 {
		mc.Parent = &model.Category{
			EntityID: base.EntityID{
				ID: cc.ParentID,
			},
		}
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
		Name:       ic.Name,
		DivisionBy: ic.DivisionBy,
	}

	if ic.ParentID != nil {
		m.Parent = &model.Category{
			EntityID: base.EntityID{
				ID: *ic.ParentID,
			},
		}
	}

	return m
}
