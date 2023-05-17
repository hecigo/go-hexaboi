package orientdb

import (
	log "github.com/sirupsen/logrus"

	"hecigo.com/go-hexaboi/domain/base"
	"hecigo.com/go-hexaboi/infra/core"
	"hecigo.com/go-hexaboi/infra/orm"
)

type ItemRepository struct {
	base.Repository
}

func (r *ItemRepository) GetByCode(code string) (*orm.Item, error) {
	var (
		funcErr    Errors
		funcResult Result
	)
	resp, err := Client().R().
		SetError(&funcErr).
		SetResult(&funcResult).
		SetPathParam("func_name", "get_item").
		SetBody(map[string]interface{}{"code": code}).
		Post(CMD_FUNCTION)
	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		log.Error(funcErr.Errors[0])
		return nil, funcErr.Errors[0].ToHPIError()
	}

	// TODO: Map Result to orm.Item
	var item *orm.Item
	core.UnmarshalNoPanic(funcResult.Result[0], &item)

	return item, nil
}
