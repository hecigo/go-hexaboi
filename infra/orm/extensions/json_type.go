package extensions

import (
	"database/sql/driver"
	"errors"
	"fmt"

	"github.com/goccy/go-json"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type JSON map[string]interface{}

// Scan scan value into Jsonb, implements sql.Scanner interface
func (j *JSON) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSON value:", value))
	}

	if err := json.Unmarshal(bytes, &j); err != nil {
		return err
	}
	return nil
}

// Value return json value, implement driver.Valuer interface
func (j JSON) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}

	str, err := json.Marshal(j)
	return string(str), err
}

func (j *JSON) Marshal(m interface{}) error {
	bytes, err := json.Marshal(m)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(bytes, &j); err != nil {
		return err
	}
	return nil
}

func (j JSON) Unmarshal() map[string]string {
	if len(j) == 0 {
		return nil
	}

	str, err := json.Marshal(j)
	if err != nil {
		return nil
	}

	var result map[string]string
	json.Unmarshal([]byte(str), &result)

	return result
}

func (JSON) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	// use field.Tag, field.TagSettings gets field's tags
	// checkout https://github.com/go-gorm/gorm/blob/master/schema/field.go for all options

	// returns different database type based on driver name
	switch db.Dialector.Name() {
	case "mysql", "sqlite":
		return "JSON"
	case "postgres":
		return "JSONB"
	}
	return ""
}
