package postgres

import (
	"hoangphuc.tech/hercules/domain/base"
	"hoangphuc.tech/hercules/infra/orm"
)

type ItemRepository struct {
	base.Repository
}

// Create item
func (r *ItemRepository) Create(item *orm.Item) error {
	result := DB().Create(&item)
	return result.Error
}
