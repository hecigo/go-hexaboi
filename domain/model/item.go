package model

import (
	base "hoangphuc.tech/hercules/domain/base"
)

// Items can be goods, products, gifts, services...
type Item struct {
	base.Entity
	SKU  string `json:"sku"`
	Name string `json:"name"`

	// This SKU represents a group of SKU of the same type.
	// It is the parent/configurable SKU as well.
	MasterSKU string `json:"master_sku"`

	// The key-value collection of variant attributes. Ex: color, size...
	VariantAttributes map[string]string `json:"variant_attributes"`
}
