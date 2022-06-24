package adapter

import (
	"hoangphuc.tech/dora/domain/model"
	"hoangphuc.tech/dora/infra/postgres"
)

type CategoryRepository struct{}

var (
	repoCate postgres.CategoryRepository = postgres.CategoryRepository{}
)

func (*CategoryRepository) GetByCode(code string) (*model.Category, error) {
	// o, err := repoCate.GetByCode(code)
	// if err != nil {
	// 	return nil, err
	// }

	// var m model.Category
	// o.ToModel(&m)
	// return &m, nil
	return nil, nil
}

func (*CategoryRepository) GetByID(id uint) (*model.Category, error) {
	// ormCate, err := repoCate.GetByID(id)
	// if err != nil {
	// 	return nil, err
	// }

	// var cate model.Category
	// ormCate.ToModel(&cate)
	// return &cate, nil
	return nil, nil
}

// Query all brand from BigQuery
func (*CategoryRepository) BQFindAll() ([]*model.Category, error) {
	return nil, nil
}
