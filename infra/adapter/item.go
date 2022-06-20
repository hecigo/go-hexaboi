package adapter

import (
	"hoangphuc.tech/go-hexaboi/domain/model"
	"hoangphuc.tech/go-hexaboi/infra/orientdb"
	"hoangphuc.tech/go-hexaboi/infra/postgres"
)

type ItemRepository struct{}

var (
	repoItem      postgres.ItemRepository = postgres.ItemRepository{}
	graphRepoItem orientdb.ItemRepository = orientdb.ItemRepository{}
)

func (*ItemRepository) GetByID(id uint) (*model.Item, error) {
	// o, err := repoItem.GetByID(id)
	// if err != nil {
	// 	return nil, err
	// }

	// var m model.Item
	// o.ToModel(&m)
	// return &m, nil
	return nil, nil
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
