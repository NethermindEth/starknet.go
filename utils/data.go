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

func GetTypedFieldFromJSON[T any](data map[string]interface{}, tag string) (T, error) {
	var resp T
	rawResp, ok := data[tag]
	if !ok {
		return resp, fmt.Errorf("missing '%s' field in json object", tag)
	}

	resp, ok = rawResp.(T)
	if !ok {
		return resp, fmt.Errorf("expected type '%T', got '%T'", resp, rawResp)
	}

	return resp, nil
}
