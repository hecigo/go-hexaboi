package core

import (
	"errors"
	"fmt"

	"github.com/goccy/go-json"
)

type JSON json.RawMessage

func UnmarshalNoPanic(origin interface{}, dest interface{}) error {
	if dest == nil {
		return errors.New("failed to unmarshal JSON: unknown destination type")
	}

	var bytes []byte
	switch origin := origin.(type) {
	case []byte:
		bytes = origin
	case string:
		bytes = []byte(origin)
	case map[string]interface{}:
		mbytes, err := json.Marshal(origin)
		if err != nil {
			return err
		}
		bytes = mbytes
	case map[string]string:
		mbytes, err := json.Marshal(origin)
		if err != nil {
			return err
		}
		bytes = mbytes
	default:
		return fmt.Errorf("failed to unmarshal JSON: origin type is invalid\n%v", origin)
	}

	err := json.Unmarshal(bytes, &dest)
	if err != nil {
		mbytes, e := json.Marshal(origin)
		if e != nil {
			return fmt.Errorf("failed to unmarshal JSON: %v\n%v", e, origin)
		}
		return fmt.Errorf("failed to unmarshal JSON: %v\n%v", err, string(mbytes))
	}

	return nil
}
