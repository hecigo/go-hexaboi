package postgres

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"hoangphuc.tech/go-hexaboi/domain/base"
	"hoangphuc.tech/go-hexaboi/infra/orm"
)

type BrandRepository struct {
	base.Repository
}

// Create brand
func (r *BrandRepository) Create(brand *orm.Brand) error {
	result := DB().Omit(clause.Associations).Create(&brand)
	return result.Error
}

// Batch create brand
func (r *BrandRepository) BatchCreate(brand []*orm.Brand) (int64, error) {
	result := DB().Omit(clause.Associations).Create(&brand)
	return result.RowsAffected, result.Error
}

// Get brand by ID
func (r *BrandRepository) GetByID(id uint) (*orm.Brand, error) {
	var o orm.Brand
	result := DB().Take(&o, id)
	return &o, result.Error
}

func (r *BrandRepository) GetByCode(code string) (*orm.Brand, error) {
	if code == "" {
		return nil, nil
	}
	var brand orm.Brand
	result := DB().Where(&orm.Brand{Code: code}).Take(&brand)
	return &brand, result.Error
}

func (r *BrandRepository) Update(brand *orm.Brand) error {
	result := DB().Clauses(clause.Returning{}).Updates(brand)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected <= 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
