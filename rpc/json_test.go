package rpc

import (
	"encoding/json"
	"errors"
	"testing"
)

// TestJSONMultiTypeArrayWithDuck tests the JSONMultiTypeArrayWithDuck function and shows how you can guess what a type is and apply it.
//
// The function tests the unmarshaling of a JSON string into a struct with a multi-type array field.
// It checks if the unmarshaling is successful and if the correct types are assigned to the elements of the array.
// The function takes a *testing.T parameter and does not return anything.
func TestJSONMultiTypeArrayWithDuck(t *testing.T) {
	type V1 struct {
		Label1 string
	}
	type V2 struct {
		Label2 string
	}
	type V interface{}
	type MyType struct {
		Data string
		Tx   []V
	}
	var my MyType

	jsonContent := `{"data": "data", "tx": [{"label2": "yes"}, {"label1": "no"}]}`
	err := json.Unmarshal([]byte(jsonContent), &my)
	if err != nil {
		t.Fatal("should succeed, instead", err)
	}

	for key, value := range my.Tx {
		local, ok := value.(map[string]interface{})
		if !ok {
			t.Fatalf("you should have found map[string]interface{}, instead %T", value)
		}
		labelOutput, err := json.Marshal(local)
		if err != nil {
			t.Fatal("label1Output should succeed, instead", err)
		}
		if _, ok := local["label1"]; ok {
			localType := V1{}
			err = json.Unmarshal(labelOutput, &localType)
			if err != nil {
				t.Fatal("V1 should succeed, instead", err)
			}
			my.Tx[key] = localType
			continue
		}
		if _, ok := local["label2"]; ok {
			localType := V2{}
			err = json.Unmarshal(labelOutput, &localType)
			if err != nil {
				t.Fatal("V1 should succeed, instead", err)
			}
			my.Tx[key] = localType
			continue
		}
		t.Fatal("you should have found a type", errors.New("missing type"))
	}
	if _, ok := my.Tx[0].(V2); !ok {
		t.Fatalf("Tx[0] should be a V2, instead, %T", my.Tx[0])
	}
	if _, ok := my.Tx[1].(V1); !ok {
		t.Fatalf("Tx[0] should be a V1, instead, %T", my.Tx[1])
	}
}

// TestJSONMixingStructWithMap tests the JSON unmarshaling of a struct that mixes fields and a map and shows how 2 embedded type are loaded but not map[string]interface{}.
//
// This function verifies that the JSON unmarshaling of a struct, which includes fields and a map,
// correctly assigns the values from the JSON content.
// It creates a struct, `MyType`, that contains fields from `V1`, `V2`, and `V3`.
// Then, it attempts to unmarshal a JSON string into the `MyType` struct and checks if the values
// were correctly assigned.
//
// Parameters:
// - t: A testing.T object used for reporting test failures and logs.
//
// Return type: None.
func TestJSONMixingStructWithMap(t *testing.T) {
	type V1 struct {
		Label1 string
	}
	type V2 struct {
		Label2 string
		Label3 string
	}
	type V3 map[string]interface{}
	type MyType struct {
		V1
		V2
		V3
	}
	var my MyType
	jsonContent := `{"label2": "label2", "label1": "label1", "label3": "label3", "label4": "label4"}`

	err := json.Unmarshal([]byte(jsonContent), &my)
	if err != nil {
		t.Fatal("should succeed, instead", err)
	}
	if my.V1.Label1 != "label1" {
		t.Fatalf("V1.Label1 should be \"label1\", instead %q", my.V1.Label1)
	}
	if len(my.V3) != 0 {
		t.Fatal("Unfortunately, nothing should be loaded in this map, yet", len(my.V3))
	}
}
