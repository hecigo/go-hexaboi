package adapter

import (
	"hoangphuc.tech/hercules/infra/orm"
	"hoangphuc.tech/hercules/infra/postgres"
)

type ItemRepository struct{}

var (
	repoItem postgres.ItemRepository = postgres.ItemRepository{}
)

func (*ItemRepository) Create(item *orm.Item) (err error) {
	return repoItem.Create(item)
}
