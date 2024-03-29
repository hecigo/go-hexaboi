package adapter

import (
	"context"

	"hecigo.com/go-hexaboi/domain/model"
	"hecigo.com/go-hexaboi/infra/bigquery"
	"hecigo.com/go-hexaboi/infra/elasticsearch"
	"hecigo.com/go-hexaboi/infra/orientdb"
	"hecigo.com/go-hexaboi/infra/orm"
	"hecigo.com/go-hexaboi/infra/postgres"
)

type ItemRepository struct{}

var (
	ES_INDEX                                = "hpi-salesorder-002"
	ES_DOC_ID_FIELD                         = "productCode"
	repoItem        postgres.ItemRepository = postgres.ItemRepository{}
	bqRepoItem      bigquery.ItemRepository = bigquery.ItemRepository{}
	graphRepoItem   orientdb.ItemRepository = orientdb.ItemRepository{}
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
	items, err := bqRepoItem.FindAllOnPage(page, pageSize)
	if err != nil {
		return nil, err
	}

	var r []*model.Item
	for _, b := range items {
		var m model.Item
		b.ToModel(&m)
		r = append(r, &m)
	}
	return r, nil
}

func (*ItemRepository) BQCount() (count uint64, err error) {
	count, err = bqRepoItem.Count()
	return count, err
}

func (*ItemRepository) Search(ctx context.Context, query interface{}, result interface{}) (total uint64, aggs map[string]interface{}, err error) {
	return elasticsearch.Search(ctx, ES_INDEX, query, result)
}

func (*ItemRepository) SearchIndex(ctx context.Context, items ...interface{}) error {
	return elasticsearch.Index(ctx, ES_INDEX, ES_DOC_ID_FIELD, items...)
}
