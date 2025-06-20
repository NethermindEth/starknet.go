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

// TODO: add docs
type UDCOptions struct {
	Salt       *felt.Felt
	Unique     bool
	UDCVersion UDCVersion
}

// TODO: add docs
type UDCVersion int

const (
	UDCCairoV0 UDCVersion = iota
	UDCCairoV2
)

// TODO: add docs
func BuildUDCCalldata(
	classHash *felt.Felt,
	constructorCalldata []*felt.Felt,
	deployOpts *UDCOptions,
) (*rpc.InvokeFunctionCall, error) {
	if classHash == nil {
		return nil, errors.New("classHash not provided")
	}

	if deployOpts == nil {
		deployOpts = new(UDCOptions)
	}

	// salt
	if deployOpts.Salt == nil {
		randFelt, err := new(felt.Felt).SetRandom()
		if err != nil {
			return nil, err
		}
		deployOpts.Salt = randFelt
	}

	// unique
	uniqueFelt := new(felt.Felt).SetUint64(0)
	if deployOpts.Unique {
		uniqueFelt = new(felt.Felt).SetUint64(1)
	}

	// response
	var udcCallData []*felt.Felt
	var udcAddress *felt.Felt
	var methodName string

	switch deployOpts.UDCVersion {
	case UDCCairoV0:
		calldataLen := new(felt.Felt).SetUint64(uint64(len(constructorCalldata)))
		udcCallData = append([]*felt.Felt{classHash, deployOpts.Salt, uniqueFelt, calldataLen}, constructorCalldata...)
		udcAddress = udcAddressCairoV0
		methodName = "deployContract"
	case UDCCairoV2:
		udcCallData = append([]*felt.Felt{classHash, deployOpts.Salt, uniqueFelt}, constructorCalldata...)
		udcAddress = udcAddressCairoV2
		methodName = "deploy_contract"
	}

	fnCall := rpc.InvokeFunctionCall{
		ContractAddress: udcAddress,
		FunctionName:    methodName,
		CallData:        udcCallData,
	}

	return &fnCall, nil
}
