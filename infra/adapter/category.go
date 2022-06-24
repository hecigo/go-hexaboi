package adapter

import (
	"hoangphuc.tech/go-hexaboi/domain/model"
	"hoangphuc.tech/go-hexaboi/infra/bigquery"
	"hoangphuc.tech/go-hexaboi/infra/orm"
	"hoangphuc.tech/go-hexaboi/infra/postgres"
)

type CategoryRepository struct{}

var (
	repoCate   postgres.CategoryRepository = postgres.CategoryRepository{}
	bqRepoCate bigquery.CategoryRepository = bigquery.CategoryRepository{}
)

func (*CategoryRepository) Create(cate *model.Category) error {
	// Convert payload to orm.Item
	ormCate := orm.NewCategory(cate)
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

func (*CategoryRepository) BatchCreate(m []*model.Category) (int64, error) {
	var ormRecords []*orm.Category
	for _, mi := range m {
		o := orm.NewCategory(mi)
		ormRecords = append(ormRecords, o)
	}

	if count, err := repoCate.BatchCreate(ormRecords); err != nil {
		return count, err
	}

	for i, o := range ormRecords {
		o.ToModel(m[i])
	}

	return int64(len(ormRecords)), nil
}

func (*CategoryRepository) Update(id uint, m *model.Category) error {
	o := orm.NewCategory(m)
	o.ID = id
	if err := repoCate.Update(o); err != nil {
		return err
	}
	o.ToModel(m)
	return nil
}

func (*CategoryRepository) GetByCode(code string) (*model.Category, error) {
	o, err := repoCate.GetByCode(code)
	if err != nil {
		return nil, err
	}

	var m model.Category
	o.ToModel(&m)
	return &m, nil
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

// Query all brand from BigQuery
func (*CategoryRepository) BQFindAll() ([]*model.Category, error) {
	cates, err := bqRepoCate.FindAll()
	if err != nil {
		return nil, err
	}

	var r []*model.Category
	for _, b := range cates {
		var m model.Category
		b.ToModel(&m)
		r = append(r, &m)
	}
	return r, nil
}
