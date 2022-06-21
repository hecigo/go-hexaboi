package orientdb

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"hoangphuc.tech/dora/domain/base"
	"hoangphuc.tech/dora/infra/core"
	"hoangphuc.tech/dora/infra/orm"
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
	var item *orm.Item
	err2 := core.Utils.MapToStruct(funcResult.Result[0], &item)

	if err2 != nil {
		fmt.Printf("MapToStruct falied ", err2)
	}
	if item != nil {
		fmt.Printf("item: ", item.Code)
	}
	if len(funcErr.Errors) > 0 {
		return item, funcErr.Errors[0]
	}
	return item, nil
}
