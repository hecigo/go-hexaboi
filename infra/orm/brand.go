package orm

import (
	model "hecigo.com/go-hexaboi/domain/model"
)

type Brand struct {
	EntityID
	Code string `json:"code" gorm:"uniqueIndex; not null; type:varchar(50); check:code <> ''"`
	Name string `json:"name" gorm:"not null; check:name <> ''"`
	Entity
}

// Initialize orm.Brand from model.Brand
func NewBrand(b *model.Brand) *Brand {
	return &Brand{
		EntityID: *NewEntityID(b.ID),
		Code:     b.Code,
		Name:     b.Name,
		Entity:   *NewEntity(b.Entity),
	}
}

// Scan orm.Brand into model.Brand
func (b *Brand) ToModel(brand *model.Brand) {
	if brand == nil {
		brand = new(model.Brand)
	}
	b.Entity.ToModel(&brand.Entity)
	brand.ID = b.ID
	brand.Code = b.Code
	brand.Name = b.Name
}
