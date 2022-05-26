package postgres

import (
	"gorm.io/gorm/clause"
	"hoangphuc.tech/hercules/domain/base"
	"hoangphuc.tech/hercules/infra/orm"
)

type CategoryRepository struct {
	base.Repository
}

// Create category
func (r *CategoryRepository) Create(cate *orm.Category) error {
	result := DB().Omit(clause.Associations).Create(&cate)
	return result.Error
}

// Get category by ID
func (r *CategoryRepository) GetByID(id uint) (*orm.Category, error) {
	if id < 1 {
		return nil, nil
	}
	var cate orm.Category
	result := DB().Joins("Parent").Take(&cate, id)
	return &cate, result.Error
}

// Get category by ID without associations
func (r *CategoryRepository) GetByIDNoRef(id uint) (*orm.Category, error) {
	if id < 1 {
		return nil, nil
	}
	var cate orm.Category
	result := DB().Take(&cate, id)
	return &cate, result.Error
}

func (r *CategoryRepository) Update(cate *orm.Category) error {
	result := DB().Clauses(clause.Returning{}).Omit("Parent").Updates(cate)
	if result.Error != nil {
		return result.Error
	}
	return result.Error
}
