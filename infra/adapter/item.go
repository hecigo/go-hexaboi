package adapter

import (
	"hoangphuc.tech/go-hexaboi/domain/model"
	"hoangphuc.tech/go-hexaboi/infra/orientdb"
	"hoangphuc.tech/go-hexaboi/infra/orm"
	"hoangphuc.tech/go-hexaboi/infra/postgres"
)

type ItemRepository struct{}

var (
	repoItem      postgres.ItemRepository = postgres.ItemRepository{}
	graphRepoItem orientdb.ItemRepository = orientdb.ItemRepository{}
)

func (*ItemRepository) Create(m *model.Item) error {
	o := orm.NewItem(m)
	if err := repoItem.Create(o); err != nil {
		return err
	}

	// Fetch primary_category
	cate, err := repoCate.GetByCodeNoJoins(o.PrimaryCategoryCode)
	if err != nil {
		return err
	}
	o.PrimaryCategory = *cate

	// Fetch brand
	brand, err := repoBrand.GetByCode(o.BrandCode)
	if err != nil {
		return err
	}
	o.Brand = *brand

	o.ToModel(m)
	return nil
}

func (*ItemRepository) BatchCreate(m []*model.Item) (int64, error) {
	var ormRecords []*orm.Item
	for _, mi := range m {
		o := orm.NewItem(mi)
		ormRecords = append(ormRecords, o)
	}

	if count, err := repoItem.BatchCreate(ormRecords); err != nil {
		return count, err
	}

	for i, o := range ormRecords {
		o.ToModel(m[i])
	}

	return int64(len(ormRecords)), nil
}

func (*ItemRepository) Update(id uint, m *model.Item) error {
	o := orm.NewItem(m)
	o.ID = id
	if err := repoItem.Update(o); err != nil {
		return err
	}
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
	o, err := graphRepoItem.GetByCode(code)
	if err != nil || o == nil {
		return nil, err
	}

	var m model.Item
	o.ToModel(&m)
	return &m, nil
}

// Query all brand from BigQuery
func (*ItemRepository) BQFindAll(page uint, pageSize uint) ([]*model.Item, error) {
	return nil, nil
}

func (*ItemRepository) BQCount() (count uint64, err error) {
	return 0, nil
}
