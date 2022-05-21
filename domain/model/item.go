package model

import (
	base "hoangphuc.tech/hercules/domain/base"
)

// Items can be goods, products, gifts, services...
type Item struct {
	base.Entity
	Code string `json:"code,omitempty"`
	Name string `json:"name,omitempty"`

	// This SKU represents a group of SKU of the same type.
	// It is the parent/configurable SKU as well.
	MasterSKU string `json:"master_sku,omitempty"`

	// The key-value collection of variant attributes. Ex: color, size...
	VariantAttributes map[string]string `json:"variant_attributes,omitempty"`
}
