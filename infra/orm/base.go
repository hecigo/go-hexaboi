package orm

import (
	"time"

	domain "hoangphuc.tech/hercules/domain/base"
)

type Entity struct {
	ID        uint64    `json:"id" gorm:"primaryKey"` // Always be auto-increment number
	CreatedBy string    `json:"created_by"`           // Username of creator
	UpdatedBy string    `json:"updated_by"`           // Username of the latest editor
	CreatedAt time.Time `json:"created_at" gorm:"index,sort:desc"`
	UpdatedAt time.Time `json:"updated_at" gorm:"index,sort:desc"`
}

func NewEntity(entity domain.Entity) *Entity {
	return &Entity{
		ID:        entity.ID,
		CreatedBy: entity.CreatedBy,
		UpdatedBy: entity.UpdatedBy,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}
}
