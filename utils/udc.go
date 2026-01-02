package utils

import (
	"errors"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/contracts"
	"github.com/NethermindEth/starknet.go/curve"
	"github.com/NethermindEth/starknet.go/rpc"
)

// TODO: migrate this to contracts package (hard to do now due to circular imports errors,
// a new types pkg would solve this)

//nolint:lll // The links would be unclickable if we break the line.
var (
	// https://voyager.online/contract/0x041a78e741e5af2fec34b695679bc6891742439f7afb8484ecd7766661ad02bf
	udcAddressCairoV0, _ = new(
		felt.Felt,
	).SetString("0x041a78e741e5af2fec34b695679bc6891742439f7afb8484ecd7766661ad02bf")
	// https://docs.openzeppelin.com/contracts-cairo/1.0.0/udc#udc_contract_address
	udcAddressCairoV2, _ = new(
		felt.Felt,
	).SetString("0x02ceed65a4bd731034c01113685c831b01c15d7d432f71afb1cf1634b53a2125")

	errInvalidUDCVersion    = errors.New("invalid UDC version")
	errClassHashNotProvided = errors.New("classHash not provided")
)

// The options for building the UDC calldata
//

type UDCOptions struct {
	// The salt to be used for the UDC deployment. If not provided, a random value will be used.
	Salt *felt.Felt
	// This parameter is used to determine if the deployer’s address will be included in
	// the contract address calculation.
	// By making deployments dependent upon the origin address, users can reserve a whole address
	// space to prevent someone else from taking ownership of the address. Keep it `false` to include
	// the deployer’s address, and `true` to make it origin independent.
	//
	// This parameter is agnostic to the UDC version. That means that,
	// with `OriginIndependent` set to `true`:
	//   - UDCCairoV0: `unique` will be set to `false`.
	// See more at: https://github.com/starknet-io/starknet-docs/blob/aa1772da8eb42dbc8e6b26ebc37cf898c207f54e/components/Starknet/modules/architecture_and_concepts/pages/Smart_Contracts/universal-deployer.adoc#deployment-types
	//   - UDCCairoV2: `from_zero` will be set to `true`.
	// See more at: https://docs.openzeppelin.com/contracts-cairo/1.0.0/udc#deployment_types
	//nolint:lll // The links would be unclickable if we break the line.
	OriginIndependent bool
	// The UDC version to be used. If not provided, UDCCairoV0 will be used.
	UDCVersion UDCVersion;
}

// Creates a new UDCOptions instance
// UDCCairoV2 will be used as the default UDC version
func NewUDCOptions() *UDCOptions {
	return &UDCOptions{
		UDCVersion: UDCCairoV2,
	}
}

// Enum representing the UDC version to be used
type UDCVersion int

const (
	// Represents the UDC version with Cairo v0 code, with the
	// address 0x041a78e741e5af2fec34b695679bc6891742439f7afb8484ecd7766661ad02bf
	UDCCairoV0 UDCVersion = iota
	// Represents the UDC version with Cairo v2 code, with the
	// address 0x02ceed65a4bd731034c01113685c831b01c15d7d432f71afb1cf1634b53a2125
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
//   - the INVOKE txn function call to deploy the contract, including the UDC
//     address and the calldata
//   - the salt used for the UDC deployment (either the provided one or the random one)
//   - an error if any
func BuildUDCCalldata(
	classHash *felt.Felt,
	constructorCalldata []*felt.Felt,
	opts *UDCOptions,
) (rpc.InvokeFunctionCall, *felt.Felt, error) {
	var result rpc.InvokeFunctionCall
	if classHash == nil {
		return result, nil, errClassHashNotProvided
	}

	if opts == nil {
		opts = NewUDCOptions()
	}

	// salt
	if opts.Salt == nil {
		opts.Salt = felt.NewRandom[felt.Felt]()
	}

	// response
	var udcCallData []*felt.Felt
	var udcAddress *felt.Felt
	var methodName string
	var originIndFelt *felt.Felt

	switch opts.UDCVersion {
	case UDCCairoV0:
		originIndFelt = new(felt.Felt).SetUint64(1)
		if opts.OriginIndependent {
			originIndFelt.SetUint64(0)
		}

		calldataLen := new(felt.Felt).SetUint64(uint64(len(constructorCalldata)))
		udcCallData = append(
			[]*felt.Felt{classHash, opts.Salt, originIndFelt, calldataLen},
			constructorCalldata...)
		udcAddress = udcAddressCairoV0
		methodName = "deployContract"
	case UDCCairoV2:
		originIndFelt = new(felt.Felt).SetUint64(0)
		if opts.OriginIndependent {
			originIndFelt.SetUint64(1)
		}

		// The UDCCairoV2 `calldata` constructor parameter is of type `Span<felt>`, so
		// we need to pass the length of the array as the first element.
		// ref: https://book.cairo-lang.org/ch102-04-serialization-of-cairo-types.html#serialization-of-arrays-and-spans
		//nolint:lll // The link would be unclickable if we break the line.
		constructorCalldata = append(
			[]*felt.Felt{new(felt.Felt).SetUint64(uint64(len(constructorCalldata)))},
			constructorCalldata...)

		udcCallData = append(
			[]*felt.Felt{classHash, opts.Salt, originIndFelt},
			constructorCalldata...)
		udcAddress = udcAddressCairoV2
		methodName = "deploy_contract"
	default:
		return result, nil, errInvalidUDCVersion
	}

	result = rpc.InvokeFunctionCall{
		ContractAddress: udcAddress,
		FunctionName:    methodName,
		CallData:        udcCallData,
	}

	return result, opts.Salt, nil
}

// Precomputes the address for a UDC deployment.
//
// Parameters:
//   - classHash: the class hash of the contract to deploy
//   - salt: the salt to be used for the UDC deployment
//   - constructorCalldata: the calldata to pass to the constructor of the contract
//   - udcVersion: the UDC version to be used
//   - originAccAddress: the address of the account that will deploy the contract. It
//     must be `nil` if `OriginIndependent` is `true`.
//
// Returns:
//   - the precomputed address for the UDC deployment
func PrecomputeAddressForUDC(
	classHash *felt.Felt,
	salt *felt.Felt,
	constructorCalldata []*felt.Felt,
	udcVersion UDCVersion,
	originAccAddress *felt.Felt,
) *felt.Felt {
	// Origin-independent deployments (deployed from zero)
	if originAccAddress == nil {
		return contracts.PrecomputeAddress(
			&felt.Zero,
			salt,
			classHash,
			constructorCalldata,
		)
	}

	// Origin-dependent deployments (deployed from origin)
	var hashedSalt *felt.Felt
	var finalOriginAddress *felt.Felt

	switch udcVersion {
	case UDCCairoV0:
		hashedSalt = curve.Pedersen(originAccAddress, salt)
		finalOriginAddress = udcAddressCairoV0
	case UDCCairoV2:
		hashedSalt = curve.PoseidonArray(originAccAddress, salt)
		finalOriginAddress = udcAddressCairoV2
	default:
		return nil
	}

	return contracts.PrecomputeAddress(
		finalOriginAddress,
		hashedSalt,
		classHash,
		constructorCalldata,
	)
}
