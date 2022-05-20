package core

import (
	"github.com/goccy/go-json"
)

type Utils struct {
}

func (u Utils) MapToStruct(source map[string]string, dest interface{}) (interface{}, error) {
	// Convert map to json string
	jsonStr, err := json.Marshal(source)
	if err != nil {
		return nil, err
	}

	// Convert json string to struct
	if err := json.Unmarshal(jsonStr, dest); err != nil {
		return nil, err
	}

	return dest, nil
}
