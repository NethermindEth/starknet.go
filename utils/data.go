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

func GetTypedFieldFromJSONMap[T any](data map[string]json.RawMessage, tag string) (resp T, err error) {
	rawResp, ok := data[tag]
	if !ok {
		return resp, fmt.Errorf("missing '%s' field in json object", tag)
	}

	if err := json.Unmarshal(rawResp, &resp); err != nil {
		return resp, err
	}

	return resp, nil
}
