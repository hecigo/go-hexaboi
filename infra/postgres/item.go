package postgres

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"hoangphuc.tech/hercules/domain/base"
	"hoangphuc.tech/hercules/infra/orm"
)

type ItemRepository struct {
	base.Repository
}

// Create item
func (r *ItemRepository) Create(item *orm.Item) error {
	// For many2many associations, GORM will upsert the associations before creating the join table references.
	// if you want to skip the upserting of associations, you could skip it like Categories.*
	result := DB().Omit("PrimaryCategory, Brand, Categories.*").Create(&item)
	return result.Error
}

// Update item
func (r *ItemRepository) Update(item *orm.Item) error {
	result := DB().Clauses(clause.Returning{}).Omit("Code, PrimaryCategory, Brand, Categories.*").Updates(item)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected <= 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// Get item by ID
func (r *ItemRepository) GetByID(id uint) (*orm.Item, error) {
	if id == 0 {
		return nil, nil
	}
	var item orm.Item
	result := DB().Joins("PrimaryCategory").Joins("Brand").Preload("Categories").Take(&item, id)
	return &item, result.Error
}

// Get item by code
func (r *ItemRepository) GetByCode(code string) (*orm.Item, error) {
	if code == "" {
		return nil, nil
	}
	var item orm.Item
	result := DB().Joins("PrimaryCategory").Joins("Brand").Preload("Categories").Where(&orm.Item{Code: code}).Take(&item)
	return &item, result.Error
}
