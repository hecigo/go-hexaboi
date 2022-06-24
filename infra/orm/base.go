package orm

import (
	"time"

	"gorm.io/gorm"
	"hoangphuc.tech/go-hexaboi/domain/base"
)

type EntityToModel interface {
	ToModel(entity interface{})
}

type EntityID struct {
	ID uint `json:"id" gorm:"primaryKey"` // Always be auto-increment number
}

func NewEntityID(id uint) *EntityID {
	return &EntityID{
		ID: id,
	}
}

type Entity struct {
	CreatedBy string    `json:"created_by"` // Username of creator
	UpdatedBy string    `json:"updated_by"` // Username of the latest editor
	CreatedAt time.Time `json:"created_at" gorm:"index,sort:desc"`
	UpdatedAt time.Time `json:"updated_at" gorm:"index,sort:desc"`
}

func NewEntity(entity base.Entity) *Entity {
	return &Entity{
		CreatedBy: entity.CreatedBy,
		UpdatedBy: entity.UpdatedBy,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}
}

func (e *Entity) ToModel(entity *base.Entity) {
	if entity == nil {
		entity = new(base.Entity)
	}
	entity.CreatedBy = e.CreatedBy
	entity.CreatedAt = e.CreatedAt
	entity.UpdatedBy = e.UpdatedBy
	entity.UpdatedAt = e.UpdatedAt
}

type FindInBatchCallback func(tx *gorm.DB, batch int) error
