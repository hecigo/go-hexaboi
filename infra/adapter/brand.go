package adapter

import (
	"hoangphuc.tech/go-hexaboi/domain/model"
	"hoangphuc.tech/go-hexaboi/infra/bigquery"
	"hoangphuc.tech/go-hexaboi/infra/orm"
	"hoangphuc.tech/go-hexaboi/infra/postgres"
)

type BrandRepository struct{}

var (
	repoBrand   postgres.BrandRepository = postgres.BrandRepository{}
	bqRepoBrand bigquery.BrandRepository = bigquery.BrandRepository{}
)

func (*BrandRepository) Create(m *model.Brand) error {
	// Convert payload to orm.Item
	o := orm.NewBrand(m)
	if err := repoBrand.Create(o); err != nil {
		return err
	}
	o.ToModel(m)
	return nil
}

func (*BrandRepository) BatchCreate(m []*model.Brand) (int64, error) {
	var ormBrands []*orm.Brand
	for _, mi := range m {
		o := orm.NewBrand(mi)
		ormBrands = append(ormBrands, o)
	}

	if count, err := repoBrand.BatchCreate(ormBrands); err != nil {
		return count, err
	}

	for i, o := range ormBrands {
		o.ToModel(m[i])
	}

	return int64(len(ormBrands)), nil
}

func (*BrandRepository) Update(id uint, m *model.Brand) error {
	o := orm.NewBrand(m)
	o.ID = id
	if err := repoBrand.Update(o); err != nil {
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

func (*BrandRepository) GetByCode(code string) (*model.Brand, error) {
	o, err := repoBrand.GetByCode(code)
	if err != nil {
		return nil, err
	}

	var m model.Brand
	o.ToModel(&m)
	return &m, nil
}

// Query all brand from BigQuery
func (*BrandRepository) BQFindAll() ([]*model.Brand, error) {
	brands, err := bqRepoBrand.FindAll()
	if err != nil {
		return nil, err
	}

	var r []*model.Brand
	for _, b := range brands {
		var m model.Brand
		b.ToModel(&m)
		r = append(r, &m)
	}
	return r, nil
}
