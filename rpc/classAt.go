package rpc

import (
	"context"
	"encoding/json"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/client/rpcerr"
	"github.com/NethermindEth/starknet.go/contracts"
)

// ClassAt returns the class at the specified blockID and contractAddress.
//
// Parameters:
//   - ctx: The context.Context object for the function
//   - blockID: The BlockID of the class
//   - contractAddress: The address of the contract
//
// Returns:
//   - ClassOutput: The output of the class
//   - error: An error if any occurred during the execution
func (provider *Provider) ClassAt(
	ctx context.Context,
	blockID BlockID,
	contractAddress *felt.Felt,
) (ClassOutput, error) {
	var rawClass map[string]any
	if err := do(
		ctx, provider.c, "starknet_getClassAt", &rawClass, blockID, contractAddress,
	); err != nil {
		return nil, rpcerr.UnwrapToRPCErr(err, ErrContractNotFound, ErrBlockNotFound)
	}

	return typecastClassOutput(rawClass)
}

// typecastClassOutput typecasts the rawClass output to the appropriate ClassOutput type.
//
// Parameters:
// rawClass - A pointer to a map[string]any containing the raw class data.
// Returns:
//   - ClassOutput: a ClassOutput interface
//   - error: an error if any
func typecastClassOutput(rawClass map[string]any) (ClassOutput, error) {
	rawClassByte, err := json.Marshal(rawClass)
	if err != nil {
		return nil, rpcerr.Err(rpcerr.InternalError, StringErrData(err.Error()))
	}

	// if contract_class_version exists, then it's a ContractClass type
	if _, exists := (rawClass)["contract_class_version"]; exists {
		var contractClass contracts.ContractClass
		err = json.Unmarshal(rawClassByte, &contractClass)
		if err != nil {
			return nil, rpcerr.Err(rpcerr.InternalError, StringErrData(err.Error()))
		}

		return &contractClass, nil
	}
	var depContractClass contracts.DeprecatedContractClass
	err = json.Unmarshal(rawClassByte, &depContractClass)
	if err != nil {
		return nil, rpcerr.Err(rpcerr.InternalError, StringErrData(err.Error()))
	}

	return &depContractClass, nil
}
