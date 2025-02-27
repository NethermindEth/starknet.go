package utils

import "github.com/NethermindEth/starknet.go/internal/utils"

// UnmarshallJSONFileToType reads a JSON file at the given path and unmarshals it into the specified type T.
// If any error occurs during file reading or unmarshalling, it returns an error.
//
// Parameters:
// - filePath: string path to the JSON file
// - subfield: string subfield to unmarshal from the JSON file
// Returns:
// - *T: pointer to the unmarshalled data of type T
// - error: error if file reading or unmarshalling fails
func UnmarshallJSONFileToType[T any](filePath string, subfield string) (*T, error) {
	return utils.UnmarshallJSONFileToType[T](filePath, subfield)
}
