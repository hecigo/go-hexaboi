package orientdb

import (
	"github.com/goccy/go-json"
)

type Result struct {
	Result []map[string]interface{} `json:"result"`
}

func (r Result) String() string {
	bytes, err := json.Marshal(r)
	if err != nil {
		return err.Error()
	}
	return string(bytes)
}
