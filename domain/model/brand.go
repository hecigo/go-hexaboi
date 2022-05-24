package model

import "hoangphuc.tech/hercules/domain/base"

type Brand struct {
	base.EntityID
	Name string `json:"name"`
	base.Entity
}
