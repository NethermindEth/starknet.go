package main

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/NethermindEth/caigo/types"
)

var knownEntries = []string{
	"get_nonce",
	"__execute__",
	"__validate__",
	"__validate_deploy__",
	"__validate_declare__",
	"initialize",
	"tokenID",
	"balanceOf",
}

func TestDictionaryGenerate(t *testing.T) {
	keys := map[string]string{}
	for _, i := range knownEntries {
		keys[fmt.Sprintf("0x%s", types.GetSelectorFromName(i).Text(16))] = i
	}
	content, err := json.Marshal(keys)
	if err != nil {
		t.Error(err)
	}
	err = os.WriteFile("dictionary.json", content, 0755)
	if err != nil {
		t.Error(err)
	}
}
