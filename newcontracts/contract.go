package newcontract

import (
	"encoding/json"

	"github.com/NethermindEth/starknet.go/rpc"
)

func UnmarshalContractClass(compiledClass []byte) (*rpc.ContractClass, error) {
	var class rpc.ContractClass
	err := json.Unmarshal(compiledClass, &class)
	if err != nil {
		return nil, err
	}
	return &class, nil
}
