package postgres

import (
	"gorm.io/gorm/clause"
	"hoangphuc.tech/hercules/domain/base"
	"hoangphuc.tech/hercules/infra/orm"
)

type BrandRepository struct {
	base.Repository
}

// Create brand
func (r *BrandRepository) Create(cate *orm.Brand) error {
	result := DB().Omit(clause.Associations).Create(&cate)
	return result.Error
}

// Get brand by ID
func (r *BrandRepository) GetByID(id uint) (*orm.Brand, error) {
	var o orm.Brand
	result := DB().Take(&o, id)
	return &o, result.Error
}
