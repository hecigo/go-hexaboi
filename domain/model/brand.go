package model

import "hoangphuc.tech/hercules/domain/base"

type Brand struct {
	base.EntityID
	Code string `json:"code"`
	Name string `json:"name"`
	base.Entity
}

func (b Brand) String() string {
	return b.Code + "\t" + b.Name
}
