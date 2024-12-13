package utils

import (
	"encoding/json"
	"fmt"
)

func UnwrapJSON(data map[string]interface{}, tag string) (map[string]interface{}, error) {
	if data[tag] != nil {
		var unwrappedData map[string]interface{}
		dataInner, err := json.Marshal(data[tag])
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(dataInner, &unwrappedData); err != nil {
			return nil, err
		}
		return unwrappedData, nil
	}
	return data, nil
}

func GetAndUnmarshalJSONFromMap[T any](aMap map[string]json.RawMessage, key string) (result T, err error) {
	value, ok := aMap[key]
	if !ok {
		return result, fmt.Errorf("invalid json: missing field %s", key)
	}

	err = json.Unmarshal(value, &result)
	if err != nil {
		return result, err
	}

	return result, nil
}
