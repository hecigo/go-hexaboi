package orm

import (
	model "hoangphuc.tech/hercules/domain/model"
)

type Brand struct {
	EntityID
	Name string `json:"name" gorm:"not null; check:name <> ''"`
	Entity
}

// Initialize orm.Brand from model.Brand
func NewBrand(b *model.Brand) *Brand {
	return &Brand{
		EntityID: *NewEntityID(b.ID),
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
	brand.Name = b.Name
}
