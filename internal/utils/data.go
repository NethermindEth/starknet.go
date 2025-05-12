package utils

import (
	"encoding/json"
	"fmt"
	"os"
)

// UnwrapJSON unwraps a JSON value from a map into a new map.
//
// Parameters:
//  - data: A map containing JSON raw messages
//  - tag: The key to look up in the map

func UnwrapJSON(data map[string]any, tag string) (map[string]any, error) {
	if data[tag] != nil {
		var unwrappedData map[string]any
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
//   - aMap: A map containing JSON raw messages
//   - key: The key to look up in the map
//
// Returns:
//   - T: The unmarshaled value of type T
//   - error: An error if the key is not found or unmarshaling fails
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

// UnmarshalJSONFileToType reads a JSON file at the given path and unmarshals it into the specified type T.
// If any error occurs during file reading or unmarshalling, it returns an error.
//
// Parameters:
//   - filePath: string path to the JSON file
//   - subfield: string subfield to unmarshal from the JSON file
//
// Returns:
//   - *T: pointer to the unmarshalled data of type T
//   - error: error if file reading or unmarshalling fails
func UnmarshalJSONFileToType[T any](filePath, subfield string) (*T, error) {
	var result T

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	if subfield != "" {
		var tempMap map[string]json.RawMessage
		err = json.Unmarshal(data, &tempMap)
		if err != nil {
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

// RemoveFieldFromJSON removes a field from a JSON bytes slice.
//
// Parameters:
//   - jsonData: pointer to the JSON data
//   - field: string field to remove, it must be a direct child of the JSON object
//
// Returns:
//   - error: error if any
func RemoveFieldFromJSON(jsonData *[]byte, field string) error {
	var data map[string]any
	if err := json.Unmarshal(*jsonData, &data); err != nil {
		return err
	}

	delete(data, field)

	newJSONData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	*jsonData = newJSONData

	return nil
}
