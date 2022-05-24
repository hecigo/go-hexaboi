package adapter

import (
	"hoangphuc.tech/hercules/domain/model"
	"hoangphuc.tech/hercules/infra/orm"
	"hoangphuc.tech/hercules/infra/postgres"
)

type CategoryRepository struct{}

var (
	repoCate postgres.CategoryRepository = postgres.CategoryRepository{}
)

func (*CategoryRepository) Create(cate *model.Category) error {
	// Convert payload to orm.Item
	ormCate := orm.NewCategory(*cate)
	if err := repoCate.Create(ormCate); err != nil {
		return err
	}

	// Fetch parent
	if ormCate.ParentID != nil {
		parent, err := repoCate.GetByID(*ormCate.ParentID)
		if err != nil {
			return err
		}
		ormCate.Parent = parent
	}

	ormCate.ToModel(cate)
	return nil
}

func (*CategoryRepository) GetByID(id uint) (*model.Category, error) {
	ormCate, err := repoCate.GetByID(id)
	if err != nil {
		return nil, err
	}

	var cate model.Category
	ormCate.ToModel(&cate)
	return &cate, nil
}
