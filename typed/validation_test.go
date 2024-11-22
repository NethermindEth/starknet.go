package typed

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

var invalidExamples = make(map[string]TypedData)

func setup(t *testing.T) {
	t.Helper()

	//TODO: add more examples
	fileNames := []string{
		"singleType",
		"danglingType",
		"invalidTypeName1",
		"invalidTypeName2",
		"invalidTypeName3",
		"invalidTypeName4",
		"invalidTypeName5",
	}

	for _, fileName := range fileNames {
		var ttd TypedData
		content, err := os.ReadFile(fmt.Sprintf("./tests/invalidExamples/%s.json", fileName))
		require.NoError(t, err, "fail to read file: %w", err)

		err = json.Unmarshal(content, &ttd)
		require.NoError(t, err, "fail to unmarshal TypedData")

		invalidExamples[fileName] = ttd
	}
}
func TestValidateTypedData(t *testing.T) {
	setup(t)

	// iterates over valid examples
	for key, validTtd := range typedDataExamples {
		ok, err := validTtd.ValidateTypedData()
		require.NoError(t, err, key)
		require.True(t, ok)
	}

	// iterates over invalid examples
	for key, invalidTtd := range invalidExamples {
		ok, err := invalidTtd.ValidateTypedData()
		require.Error(t, err, key)
		require.False(t, ok)
	}
}
