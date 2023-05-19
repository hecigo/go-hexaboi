package orientdb

import (
	"github.com/hecigo/goutils"
	log "github.com/sirupsen/logrus"

	"hecigo.com/go-hexaboi/domain/base"
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
	item, err := goutils.Unmarshal[orm.Item](funcResult.Result[0])
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return &item, nil
}
