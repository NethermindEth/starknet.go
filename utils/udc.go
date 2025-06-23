package utils

import (
	"errors"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/rpc"
)

var (
	// https://voyager.online/contract/0x041a78e741e5af2fec34b695679bc6891742439f7afb8484ecd7766661ad02bf
	udcAddressCairoV0, _ = new(felt.Felt).SetString("0x041a78e741e5af2fec34b695679bc6891742439f7afb8484ecd7766661ad02bf")
	// https://docs.openzeppelin.com/contracts-cairo/1.0.0/udc#udc_contract_address
	udcAddressCairoV2, _ = new(felt.Felt).SetString("0x04a64cd09a853868621d94cae9952b106f2c36a3f81260f85de6696c6b050221")
)

// The options for building the UDC calldata
type UDCOptions struct {
	// The salt to be used for the UDC deployment. Default: a random value.
	Salt *felt.Felt
	// This parameter is used to determine if the deployer’s address will be included in the contract address calculation.
	// Default: the address will be included.
	// It behaves differently depending on the UDC version:
	//   - UDCCairoV0: true to include the deployer’s address in the contract address calculation, false to exclude it
	// See more at: https://github.com/starknet-io/starknet-docs/blob/aa1772da8eb42dbc8e6b26ebc37cf898c207f54e/components/Starknet/modules/architecture_and_concepts/pages/Smart_Contracts/universal-deployer.adoc#deployment-types
	//   - UDCCairoV2: behaves as `from_zero`. True to NOT include the deployer’s address in the contract address calculation, false to include it
	// See more at: https://docs.openzeppelin.com/contracts-cairo/1.0.0/udc#deployment_types
	Unique bool
	// The UDC version to be used. Default: UDCCairoV0.
	UDCVersion UDCVersion
}

// Enum representing the UDC version to be used
type UDCVersion int

const (
	// Represents the UDC version with Cairo v0 code, with the address 0x041a78e741e5af2fec34b695679bc6891742439f7afb8484ecd7766661ad02bf
	UDCCairoV0 UDCVersion = iota
	// Represents the UDC version with Cairo v2 code, with the address 0x04a64cd09a853868621d94cae9952b106f2c36a3f81260f85de6696c6b050221
	UDCCairoV2
)

// Builds the INVOKE txn function call to deploy a contract using the UDC.
//
// Parameters:
//   - classHash: the class hash of the contract to deploy
//   - constructorCalldata: the calldata to pass to the constructor of the contract
//   - opts: the options for the UDC deployment. If nil, the default options will be used.
//
// Returns:
//   - the INVOKE txn function call to deploy the contract, including the UDC address and the calldata
//   - an error if any
func BuildUDCCalldata(
	classHash *felt.Felt,
	constructorCalldata []*felt.Felt,
	opts *UDCOptions,
) (*rpc.InvokeFunctionCall, error) {
	if classHash == nil {
		return nil, errors.New("classHash not provided")
	}

	if opts == nil {
		opts = new(UDCOptions)
	}

	// salt
	if opts.Salt == nil {
		randFelt, err := new(felt.Felt).SetRandom()
		if err != nil {
			return nil, err
		}
		opts.Salt = randFelt
	}

	// unique
	uniqueFelt := new(felt.Felt).SetUint64(0)
	if opts.Unique {
		uniqueFelt = new(felt.Felt).SetUint64(1)
	}

	// response
	var udcCallData []*felt.Felt
	var udcAddress *felt.Felt
	var methodName string

	switch opts.UDCVersion {
	case UDCCairoV0:
		calldataLen := new(felt.Felt).SetUint64(uint64(len(constructorCalldata)))
		udcCallData = append([]*felt.Felt{classHash, opts.Salt, uniqueFelt, calldataLen}, constructorCalldata...)
		udcAddress = udcAddressCairoV0
		methodName = "deployContract"
	case UDCCairoV2:
		udcCallData = append([]*felt.Felt{classHash, opts.Salt, uniqueFelt}, constructorCalldata...)
		udcAddress = udcAddressCairoV2
		methodName = "deploy_contract"
	default:
		return nil, errors.New("invalid UDC version")
	}

	fnCall := rpc.InvokeFunctionCall{
		ContractAddress: udcAddress,
		FunctionName:    methodName,
		CallData:        udcCallData,
	}

	return &fnCall, nil
}
