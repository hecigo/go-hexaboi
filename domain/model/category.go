package model

import "hoangphuc.tech/go-hexaboi/domain/base"

// Item grouping, category
type Category struct {
	base.EntityID
	Code string `json:"code"`
	Name string `json:"name"`

	// Products belonging to this category are grouped by the specified division.
	DivisionBy Division `json:"division_by"`

	// Parent category
	Parent *Category `json:"parent"`
	base.Entity
}

func (b Category) String() string {
	return b.Code + "\t" + b.Name
}

// Kind of how to group products.
type Division string

const (
	DIVISION_MERCHANDISE Division = "merchandise"
	DIVISION_CONSUMER    Division = "consumer"
	DIVISION_CAMPAIGN    Division = "campaign"
	DIVISION_CUSTOM      Division = "custom"
)
