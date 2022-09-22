package rpc

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"
)

// TestJSONMultiTypeArrayWithDuck shows how you can guess what a type is and apply it
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

// TestJSONMixingStructWithMap shows how 2 embedded type are loaded but not map[string]interface{}
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

func TestJSONWithBigIntWithInterface(t *testing.T) {
	// cat account_v0.10.json | grep -e "[0-9]\{20\}
	schema := map[string]interface{}{}
	data := `
{
	"program": { 
		"content": "very complex json schema you won't map", 
		"value": 6219495360805491471215297013070624192820083
	}
}`

	err := json.Unmarshal([]byte(data), &schema)
	if err != nil {
		t.Fatal("error parsing json: ", err)
	}
	v, ok := schema["program"].(map[string]interface{})
	if !ok {
		t.Fatalf("error parsing json: %T", schema["program"])
	}
	if fmt.Sprintf("%T", v["value"]) == "float64" {
		t.Log("the value will be truncated by being turned into an float64")
	}
	if fmt.Sprintf("%v", v["value"]) != "6219495360805491471215297013070624192820083" {
		t.Fatal("the value differs from the input, current: ", fmt.Sprintf("%v", v["value"]))
	}
}

func TestJSONWithBigIntWithRawMessage(t *testing.T) {
	// cat account_v0.10.json | grep -e "[0-9]\{20\}
	schema := map[string]json.RawMessage{}
	data := `
{
	"program": { 
		"content": "very complex json schema you won't map", 
		"value": 6219495360805491471215297013070624192820083
	}
}`

	err := json.Unmarshal([]byte(data), &schema)
	if err != nil {
		t.Fatal("error parsing json: ", err)
	}
	if strings.Contains(string(schema["program"]), "6219495360805491471215297013070624192820083") {
		t.Fatal("I have lost data: ", err)
	}
}
