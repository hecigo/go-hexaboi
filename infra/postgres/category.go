package postgres

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"hoangphuc.tech/go-hexaboi/domain/base"
	"hoangphuc.tech/go-hexaboi/infra/orm"
)

type CategoryRepository struct {
	base.Repository
}

// Create category
func (r *CategoryRepository) Create(category *orm.Category) error {
	result := DB().Omit(clause.Associations).Create(&category)
	return result.Error
}

// Batch create category
func (r *CategoryRepository) BatchCreate(categories []*orm.Category) (int64, error) {
	result := DB().Omit(clause.Associations).Create(&categories)
	return result.RowsAffected, result.Error
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
func (r *CategoryRepository) GetByIDNoJoins(id uint) (*orm.Category, error) {
	if id < 1 {
		return nil, nil
	}
	var cate orm.Category
	result := DB().Take(&cate, id)
	return &cate, result.Error
}

func (r *CategoryRepository) GetByCode(code string) (*orm.Category, error) {
	if code == "" {
		return nil, nil
	}
	var cate orm.Category
	result := DB().Joins("Parent").Where(&orm.Category{Code: code}).Take(&cate)
	return &cate, result.Error
}

func (r *CategoryRepository) GetByCodeNoJoins(code string) (*orm.Category, error) {
	if code == "" {
		return nil, nil
	}
	var cate orm.Category
	result := DB().Where(&orm.Category{Code: code}).Take(&cate)
	return &cate, result.Error
}

func (r *CategoryRepository) Update(cate *orm.Category) error {
	result := DB().Clauses(clause.Returning{}).Omit("Parent").Updates(cate)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected <= 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
