package adapter

import (
	"hoangphuc.tech/hercules/domain/model"
	"hoangphuc.tech/hercules/infra/orm"
	"hoangphuc.tech/hercules/infra/postgres"
)

type ItemRepository struct{}

var (
	repoItem postgres.ItemRepository = postgres.ItemRepository{}
)

func (*ItemRepository) Create(m *model.Item) error {
	o := orm.NewItem(m)
	if err := repoItem.Create(o); err != nil {
		return err
	}

	// Fetch primary_category
	cate, err := repoCate.GetByIDNoRef(o.PrimaryCategoryID)
	if err != nil {
		return err
	}
	o.PrimaryCategory = *cate

	// Fetch brand
	brand, err := repoBrand.GetByID(o.BrandID)
	if err != nil {
		return err
	}
	o.Brand = *brand

	o.ToModel(m)
	return nil
}

func (*ItemRepository) GetByID(id uint) (*model.Item, error) {
	o, err := repoItem.GetByID(id)
	if err != nil {
		return nil, err
	}

	var m model.Item
	o.ToModel(&m)
	return &m, nil
}

func (*ItemRepository) GetByCode(code string) (*model.Item, error) {
	o, err := repoItem.GetByCode(code)
	if err != nil {
		return nil, err
	}

	var m model.Item
	o.ToModel(&m)
	return &m, nil
}
