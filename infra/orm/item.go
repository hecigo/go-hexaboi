package orm

import (
	model "hoangphuc.tech/hercules/domain/model"
)

type Item struct {
	Entity
	SKU               string            `json:"sku" gorm:"uniqueIndex; not null; type:varchar(100); check:sku <> ''"`
	Name              string            `json:"name" gorm:"not null; check:name <> ''"`
	VariantAttributes map[string]string `json:"variant_attributes" gorm:"type:json"`
	MasterSKU         string            `json:"master_sku" gorm:"type:varchar(100)"`
}

// Item constructor
func NewItem(_item *model.Item) *Item {
	return &Item{
		Entity:            *NewEntity(_item.Entity),
		SKU:               _item.SKU,
		Name:              _item.Name,
		MasterSKU:         _item.MasterSKU,
		VariantAttributes: _item.VariantAttributes,
	}
}
