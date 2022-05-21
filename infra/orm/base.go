package orm

import (
	"time"

	domain "hoangphuc.tech/hercules/domain/base"
)

type EntityID struct {
	ID uint64 `json:"id,omitempty" gorm:"primaryKey"` // Always be auto-increment number
}

func NewEntityID(id uint64) *EntityID {
	return &EntityID{
		ID: id,
	}
}

type Entity struct {
	CreatedBy string    `json:"created_by,omitempty"` // Username of creator
	UpdatedBy string    `json:"updated_by,omitempty"` // Username of the latest editor
	CreatedAt time.Time `json:"created_at,omitempty" gorm:"index,sort:desc"`
	UpdatedAt time.Time `json:"updated_at,omitempty" gorm:"index,sort:desc"`
}

func NewEntity(entity domain.Entity) *Entity {
	return &Entity{
		CreatedBy: entity.CreatedBy,
		UpdatedBy: entity.UpdatedBy,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}
}
