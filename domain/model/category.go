package model

import "hoangphuc.tech/hercules/domain/base"

type Category struct {
	base.EntityID
	Name string `json:"name"`

	// Products belonging to this category are grouped by the specified division.
	DivisionBy Division `json:"division_by"`

	// Parent category
	Parent *Category `json:"parent"`
	base.Entity
}

// Kind of how to group products.
type Division string

const (
	DIVISION_CATEGORY Division = "category"
	DIVISION_CAMPAIGN Division = "campaign"
	DIVISION_CUSTOM   Division = "custom"
)
