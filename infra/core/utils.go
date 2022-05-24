package core

import (
	"reflect"
	"strconv"
	"strings"

	"github.com/goccy/go-json"
)

var (
	Utils = utils{}
)

type utils struct {
}

// Convert a map to struct type
func (u utils) MapToStruct(source map[string]interface{}, dest interface{}) (interface{}, error) {
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

// Parse a string to uint
func (u utils) ParseUint(strID string) (id uint, err error) {
	u32, err := strconv.ParseUint(strID, 0, 32)
	if err != nil {
		return 0, err
	}

	return uint(u32), nil
}

//#region Reflect

// Get json field name from StructTag
func (u utils) GetJSONFieldName(t reflect.StructTag) string {
	jsonTag := t.Get("json")
	if jsonTag == "" {
		return ""
	}
	return strings.Split(jsonTag, ",")[0]
}

// Detect the StructField is a number field, isn't it.
func (u utils) IsNumberField(f reflect.StructField) bool {
	switch f.Type.Kind() {
	case reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Float32,
		reflect.Float64:
		return true
	default:
		return false
	}
}

//#endregion
