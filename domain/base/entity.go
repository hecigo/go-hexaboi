package base

import (
	"time"
)

type Entity struct {
	ID        uint64    `json:"id,omitempty"`
	CreatedBy string    `json:"created_by,omitempty"` // Username of creator
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedBy string    `json:"updated_by,omitempty"` // Username of the latest editor
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}
