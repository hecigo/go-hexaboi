package orm

import (
	model "hoangphuc.tech/hercules/domain/model"
)

type Brand struct {
	EntityID
	Name string `json:"name" gorm:"not null; check:name <> ''"`
	Entity
}

// Brand constructor
func NewBrand(_brand model.Brand) *Brand {
	return &Brand{
		EntityID: *NewEntityID(_brand.ID),
		Name:     _brand.Name,
		Entity:   *NewEntity(_brand.Entity),
	}
}

func (b *Brand) ToModel(brand *model.Brand) {
	if brand == nil {
		brand = new(model.Brand)
	}
	b.Entity.ToModel(&brand.Entity)
	brand.ID = b.ID
	brand.Name = b.Name
}
