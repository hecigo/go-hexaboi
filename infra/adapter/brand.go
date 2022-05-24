package adapter

import (
	"hoangphuc.tech/hercules/domain/model"
	"hoangphuc.tech/hercules/infra/orm"
	"hoangphuc.tech/hercules/infra/postgres"
)

type BrandRepository struct{}

var (
	repoBrand postgres.BrandRepository = postgres.BrandRepository{}
)

func (*BrandRepository) Create(m *model.Brand) error {
	// Convert payload to orm.Item
	o := orm.NewBrand(*m)
	if err := repoBrand.Create(o); err != nil {
		return err
	}
	o.ToModel(m)
	return nil
}

func (*BrandRepository) GetByID(id uint) (*model.Brand, error) {
	o, err := repoBrand.GetByID(id)
	if err != nil {
		return nil, err
	}

	var m model.Brand
	o.ToModel(&m)
	return &m, nil
}
