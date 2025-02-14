package utils

import (
	"encoding/json"
	"fmt"
	"os"
)

// UnwrapJSON unwraps a JSON value from a map into a new map.
//
// Parameters:
// - data: A map containing JSON raw messages
// - tag: The key to look up in the map

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

// GetAndUnmarshalJSONFromMap retrieves and unmarshals a JSON value from a map into the specified type T.
//
// Parameters:
// - aMap: A map containing JSON raw messages
// - key: The key to look up in the map
// Returns:
// - T: The unmarshaled value of type T
// - error: An error if the key is not found or unmarshaling fails
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

// UnmarshallJSONToType reads a JSON file at the given path and unmarshals it into the specified type T.
// If any error occurs during file reading or unmarshalling, it returns an error.
//
// Parameters:
// - filePath: string path to the JSON file
// - subfield: string subfield to unmarshal from the JSON file
// Returns:
// - *T: pointer to the unmarshalled data of type T
// - error: error if file reading or unmarshalling fails
func UnmarshallJSONToType[T any](filePath string, subfield string) (*T, error) {
	var result T

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	if subfield != "" {
		var tempMap map[string]json.RawMessage
		if err := json.Unmarshal(data, &tempMap); err != nil {
			return nil, err
		}

		if tempData, ok := tempMap[subfield]; ok {
			data = tempData
		} else {
			return nil, fmt.Errorf("invalid subfield: missing field %s", subfield)
		}
	}

	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
