package model

import "hoangphuc.tech/go-hexaboi/domain/base"

// Warehouse
type Warehouse struct {
	base.EntityID
	Code     string           `json:"code"`
	Name     string           `json:"name"`
	Types    []*WarehouseType `json:"types"` // Allowing a warehouse with multiple types
	Street   string           `json:"street"`
	Ward     *GEOLocation     `json:"ward"`
	District *GEOLocation     `json:"district"`
	Province *GEOLocation     `json:"province"`
	base.Entity
}

type WarehouseType string

const (
	PHYSICAL = "physical" // Including physical stock
	VIRTUAL  = "virtual"  // A virtual warehouse includes only number of stock without physical term, using to cluster stock data.
)
