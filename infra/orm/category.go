package orm

import (
	model "hecigo.com/go-hexaboi/domain/model"
)

type Category struct {
	EntityID
	Code string `json:"code" gorm:"uniqueIndex; not null; type:varchar(50);"`
	Name string `json:"name" gorm:"not null; check:name <> ''"`

	// Products belonging to this category are grouped by the specified division.
	DivisionBy model.Division `json:"division_by" gorm:"type:varchar(50); not null; check:division_by IN ('merchandise','consumer','campaign','custom')"`
	ParentID   *uint          `json:"parent_id"`
	Parent     *Category      `json:"parent" gorm:"foreignKey:ParentID; references:ID; constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Entity
}

// Initialize orm.Category from model.Category
func NewCategory(cate *model.Category) *Category {
	result := &Category{
		EntityID:   *NewEntityID(cate.ID),
		Code:       cate.Code,
		Name:       cate.Name,
		DivisionBy: cate.DivisionBy,
		Entity:     *NewEntity(cate.Entity),
	}
	if cate.Parent != nil {
		result.Parent = NewCategory(cate.Parent)
		result.ParentID = new(uint)
		*result.ParentID = cate.Parent.ID
	}

	return result
}

// Scan orm.Category into model.Category
func (c *Category) ToModel(cate *model.Category) {
	if cate == nil {
		cate = new(model.Category)
	}
	c.Entity.ToModel(&cate.Entity)
	cate.ID = c.ID
	cate.Code = c.Code
	cate.Name = c.Name
	cate.DivisionBy = c.DivisionBy
	if c.Parent != nil {
		cate.Parent = new(model.Category)
		c.Parent.ToModel(cate.Parent)
	}
}
