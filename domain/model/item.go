package model

import (
	base "hoangphuc.tech/hercules/domain/base"
)

// Items can be goods, products, gifts, services...
type Item struct {
	base.EntityID
	Code      string `json:"code"`
	Name      string `json:"name"`
	ShortName string `json:"short_name"`

	// The key-value collection of variant attributes. Ex: color, size...
	VariantAttributes map[string]string `json:"variant_attributes"`

	// This SKU represents a group of SKU of the same type.
	// It is the parent/configurable SKU as well.
	MasterSKU string `json:"master_sku"`

	// This is the principal category of SKU, defined by HPI
	PrimaryCategory Category `json:"primary_category"`

	// All categories including SKU
	Categories []Category `json:"categories"`

	// Brand
	Brand Brand `json:"brand"`
	base.Entity
}
