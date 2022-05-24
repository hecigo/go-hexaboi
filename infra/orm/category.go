package orm

import (
	model "hoangphuc.tech/hercules/domain/model"
)

type Category struct {
	EntityID
	Name string `json:"name" gorm:"not null; check:name <> ''"`

	// Products belonging to this category are grouped by the specified division.
	DivisionBy model.Division `json:"division_by" gorm:"type:varchar(50); not null; check:division_by IN ('category','campaign','custom')"`
	ParentID   *uint          `json:"parent_id"`
	Parent     *Category      `json:"parent" gorm:"foreignKey:ParentID; references:ID; constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Entity
}

// Initialize orm.Category from model.Category
func NewCategory(_cate model.Category) *Category {
	cate := &Category{
		EntityID:   *NewEntityID(_cate.ID),
		Name:       _cate.Name,
		DivisionBy: _cate.DivisionBy,
		Entity:     *NewEntity(_cate.Entity),
	}
	if _cate.Parent != nil {
		cate.Parent = NewCategory(*_cate.Parent)
		cate.ParentID = new(uint)
		*cate.ParentID = _cate.Parent.ID
	}

	return cate
}

// Scan orm.Category into model.Category
func (c *Category) ToModel(cate *model.Category) {
	if cate == nil {
		cate = new(model.Category)
	}
	c.Entity.ToModel(&cate.Entity)
	cate.ID = c.ID
	cate.Name = c.Name
	cate.DivisionBy = c.DivisionBy
	if c.Parent != nil {
		cate.Parent = new(model.Category)
		c.Parent.ToModel(cate.Parent)
	}
}
