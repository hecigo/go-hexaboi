package model

import "hoangphuc.tech/go-hexaboi/domain/base"

// Geographic location
type GEOLocation struct {
	base.EntityID
	Code       string          `json:"code"`
	Name       string          `json:"name"`
	ShortNames string          `json:"short_names"`
	Type       GEOLocationType `json:"type"`
	Parent     *GEOLocation    `json:"parent"`
	base.Entity
}

type GEOLocationType string

const (
	PROVINCE = "province"
	DISTRICT = "district"
	WARD     = "ward"
)
