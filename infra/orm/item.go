package orm

import (
	model "hoangphuc.tech/hercules/domain/model"
)

type Item struct {
	EntityID
	SKU               string            `json:"sku" gorm:"uniqueIndex; not null; type:varchar(100); check:sku <> ''"`
	Name              string            `json:"name" gorm:"not null; check:name <> ''"`
	VariantAttributes map[string]string `json:"variant_attributes" gorm:"type:jsonb"`
	MasterSKU         string            `json:"master_sku" gorm:"type:varchar(100)"`
	Entity
}

// Item constructor
func NewItem(_item *model.Item) *Item {
	return &Item{
		EntityID:          *NewEntityID(_item.ID),
		SKU:               _item.SKU,
		Name:              _item.Name,
		MasterSKU:         _item.MasterSKU,
		VariantAttributes: _item.VariantAttributes,
		Entity:            *NewEntity(_item.Entity),
	}
}
