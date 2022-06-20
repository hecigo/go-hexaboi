package orientdb

import (
	"bytes"

	"github.com/goccy/go-json"

	log "github.com/sirupsen/logrus"

	"hoangphuc.tech/go-hexaboi/domain/base"
	"hoangphuc.tech/go-hexaboi/infra/orm"
)

type ItemRepository struct {
	base.Repository
}

func (r *ItemRepository) GetByCode(code string) (*orm.Item, error) {
	var queryErr Errors
	resp, err := Client().R().SetError(&queryErr).
		SetPathParam("func_name", "get_item").
		SetBody(map[string]interface{}{"code": code}).
		Post(CMD_FUNCTION)
	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		log.Error(queryErr.Errors[0])
		return nil, queryErr.Errors[0].ToHPIError()
	}

	log.Debugln(resp)
	return nil, nil
}

func JSONMarshal(t interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(t)
	return buffer.Bytes(), err
}
