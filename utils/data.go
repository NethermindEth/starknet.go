package utils

import internalUtils "github.com/NethermindEth/starknet.go/internal/utils"

// UnmarshalJSONFileToType reads a JSON file at the given path and unmarshals it
// into the specified type T.
// If any error occurs during file reading or unmarshalling, it returns an error.
//
// Parameters:
//   - filePath: string path to the JSON file
//   - subfields: string slice of subfields to unmarshal from the JSON file
//
// You can use 'subfields' to specify the path of the type you want to unmarshal.
// Example: if you want to unmarshal the field "bar" that is within the json,
// you can specify the subfields as ["bar"].
// If you want to unmarshal the field "foo" that is within the "bar" field,
// you can specify the subfields as ["bar", "foo"].
//
// Returns:
//   - T: the unmarshalled data of type T
//   - error: error if file reading or unmarshalling fails
func UnmarshalJSONFileToType[T any](filePath string, subfields ...string) (T, error) {
	return internalUtils.UnmarshalJSONFileToType[T](filePath, subfields...)
}
