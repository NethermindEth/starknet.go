package methods

import (
	"context"
	"encoding/json"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/client/rpcerr"
	"github.com/NethermindEth/starknet.go/contracts"
	"github.com/NethermindEth/starknet.go/rpc/callers"
	"github.com/NethermindEth/starknet.go/rpc/internal"
	"github.com/NethermindEth/starknet.go/rpc/rpcv10"
)

// Class retrieves the class information from the Provider with the given hash.
//
// Parameters:
//   - ctx: The context.Context object
//   - blockID: The BlockID object
//   - classHash: The *felt.Felt object
//
// Returns:
//   - ClassOutput: The output of the class.
//   - error: An error if any occurred during the execution.
func Class(
	ctx context.Context,
	c callers.Caller,
	blockID rpcv10.BlockID,
	classHash *felt.Felt,
) (rpcv10.ClassOutput, error) {
	var rawClass map[string]any
	if err := internal.Do(ctx, c, "starknet_getClass", &rawClass, blockID, classHash); err != nil {
		return nil, rpcerr.UnwrapToRPCErr(err, rpcv10.ErrClassHashNotFound, rpcv10.ErrBlockNotFound)
	}

	return typecastClassOutput(rawClass)
}

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
func ClassAt(
	ctx context.Context,
	c callers.Caller,
	blockID rpcv10.BlockID,
	contractAddress *felt.Felt,
) (rpcv10.ClassOutput, error) {
	var rawClass map[string]any
	if err := internal.Do(
		ctx, c, "starknet_getClassAt", &rawClass, blockID, contractAddress,
	); err != nil {
		return nil, rpcerr.UnwrapToRPCErr(err, rpcv10.ErrContractNotFound, rpcv10.ErrBlockNotFound)
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
func typecastClassOutput(rawClass map[string]any) (rpcv10.ClassOutput, error) {
	rawClassByte, err := json.Marshal(rawClass)
	if err != nil {
		return nil, rpcerr.Err(rpcerr.InternalError, rpcv10.StringErrData(err.Error()))
	}

	// if contract_class_version exists, then it's a ContractClass type
	if _, exists := (rawClass)["contract_class_version"]; exists {
		var contractClass contracts.ContractClass
		err = json.Unmarshal(rawClassByte, &contractClass)
		if err != nil {
			return nil, rpcerr.Err(rpcerr.InternalError, rpcv10.StringErrData(err.Error()))
		}

		return &contractClass, nil
	}
	var depContractClass contracts.DeprecatedContractClass
	err = json.Unmarshal(rawClassByte, &depContractClass)
	if err != nil {
		return nil, rpcerr.Err(rpcerr.InternalError, rpcv10.StringErrData(err.Error()))
	}

	return &depContractClass, nil
}

// ClassHashAt retrieves the class hash at the given block ID and contract address.
//
// Parameters:
//   - ctx: The context.Context used for the request
//   - blockID: The ID of the block
//   - contractAddress: The address of the contract
//
// Returns:
//   - *felt.Felt: The class hash
//   - error: An error if any occurred during the execution
func ClassHashAt(
	ctx context.Context,
	c callers.Caller,
	blockID rpcv10.BlockID,
	contractAddress *felt.Felt,
) (*felt.Felt, error) {
	var result *felt.Felt
	if err := internal.Do(
		ctx, c, "starknet_getClassHashAt", &result, blockID, contractAddress,
	); err != nil {
		return nil, rpcerr.UnwrapToRPCErr(err, rpcv10.ErrContractNotFound, rpcv10.ErrBlockNotFound)
	}

	return result, nil
}

// Get the CASM code resulting from compiling a given class
//
// Parameters:
//   - ctx: The context.Context used for the request
//   - classHash: The hash of the contract class whose CASM will be returned
//
// Returns:
//   - CasmCompiledContractClass: The compiled contract class
//   - error: An error if any occurred during the execution
func CompiledCasm(
	ctx context.Context,
	c callers.Caller,
	classHash *felt.Felt,
) (*contracts.CasmClass, error) {
	var result contracts.CasmClass
	if err := internal.Do(ctx, c, "starknet_getCompiledCasm", &result, classHash); err != nil {
		return nil, rpcerr.UnwrapToRPCErr(
			err,
			rpcv10.ErrClassHashNotFound,
			rpcv10.ErrCompilationError,
		)
	}

	return &result, nil
}
