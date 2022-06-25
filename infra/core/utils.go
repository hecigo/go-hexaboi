package core

import (
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/goccy/go-json"
)

var (
	Utils = utils{}
)

type utils struct {
}

// Convert any object to string
func (u utils) ToStr(any interface{}) string {
	bytes, err := json.Marshal(any)
	if err != nil {
		return err.Error()
	}

	return string(bytes)
}

// Convert a map to struct type
func (u utils) MapToStruct(source map[string]interface{}, dest interface{}) error {
	// Convert map to json string
	bytes, err := json.Marshal(source)
	if err != nil {
		return err
	}

	// Convert json string to struct
	if err := json.Unmarshal(bytes, dest); err != nil {
		return err
	}

	return nil
}

func (u utils) StructToMap(source interface{}, dest *map[string]interface{}) error {
	bytes, err := json.Marshal(source)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(bytes, &dest); err != nil {
		return err
	}

	return nil
}

func (u utils) MapToStringMap(source map[string]interface{}) (dest *map[string]string, err error) {
	// Convert map to json string
	jsonStr, err := json.Marshal(source)
	if err != nil {
		return nil, err
	}

	// Convert json string to struct
	dest = new(map[string]string)
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

// Get number of days in a year
func (u utils) NumOfDaysInYear(year int) uint {
	beginOfYear := time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC)
	endOfYear := time.Date(year, time.December, 31, 0, 0, 0, 0, time.UTC)
	diff := endOfYear.Sub(beginOfYear)
	return uint(diff.Hours() / 24)
}

// Get now() timestamp as ISO-8601
func (u utils) NowStr() string {
	return time.Now().UTC().Format("2006-01-02T15:04:05-0700")
}
