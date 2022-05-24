package base

import (
	"time"
)

type Entity struct {
	CreatedBy string    `json:"created_by"` // Username of creator
	CreatedAt time.Time `json:"created_at"`
	UpdatedBy string    `json:"updated_by"` // Username of the latest editor
	UpdatedAt time.Time `json:"updated_at"`
}

type EntityID struct {
	ID uint `json:"id"`
}
