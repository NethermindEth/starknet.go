package hash

import (
	"errors"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/curve"
	"github.com/NethermindEth/starknet.go/rpc/rpcv10"
	"github.com/NethermindEth/starknet.go/rpc/rpcv9"
)

// @changed it's private now
// calculateDeprecatedTransactionHashCommon calculates the transaction hash
// common to be used in the StarkNet network - a unique identifier of the transaction.
// [specification]: https://github.com/starkware-libs/cairo-lang/blob/8276ac35830148a397e1143389f23253c8b80e93/src/starkware/starknet/core/os/transaction_hash/deprecated_transaction_hash.py#L29
//
// Parameters:
//   - txHashPrefix: The prefix of the transaction hash
//   - version: The version of the transaction
//   - contractAddress: The address of the contract
//   - entryPointSelector: The selector of the entry point
//   - calldataHash: The hashed calldata of the transaction
//   - maxFee: The maximum fee for the transaction
//   - chainID: The ID of the blockchain
//   - additionalData: Additional data to be included in the hash
//
// Returns:
//   - *felt.Felt: the calculated transaction hash
//
//nolint:lll // The link would be unclickable if we break the line.
func calculateDeprecatedTransactionHashCommon(
	txHashPrefix *felt.Felt,
	version string,
	contractAddress *felt.Felt,
	entryPointSelector *felt.Felt,
	calldataHash *felt.Felt,
	maxFee *felt.Felt,
	chainID *felt.Felt,
	additionalData []*felt.Felt,
) (*felt.Felt, error) {
	if txHashPrefix == nil ||
		version == "" ||
		contractAddress == nil ||
		entryPointSelector == nil ||
		calldataHash == nil ||
		maxFee == nil ||
		chainID == nil ||
		additionalData == nil {
		return nil, ErrNotAllParametersSet
	}

	for _, data := range additionalData {
		if data == nil {
			return nil, ErrNotAllParametersSet
		}
	}

	versionFelt, err := new(felt.Felt).SetString(version)
	if err != nil {
		return nil, err
	}

	dataToHash := []*felt.Felt{
		txHashPrefix,
		versionFelt,
		contractAddress,
		entryPointSelector,
		calldataHash,
		maxFee,
		chainID,
	}
	dataToHash = append(dataToHash, additionalData...)

	return curve.PedersenArray(dataToHash...), nil
}

// calculateV3TransactionHash calculates the hash of a V3 transaction;
// a common function to be used for all V3 transactions.
func calculateV3TransactionHash[
	R *rpcv9.ResourceBoundsMapping | *rpcv10.ResourceBoundsMapping,
](
	prefix *felt.Felt,
	version string,
	contractAddress *felt.Felt,
	tip interface{ ToUint64() (uint64, error) },
	resourceBounds R,
	paymasterData []*felt.Felt,
	chainID *felt.Felt,
	nonce *felt.Felt,
	feeMode interface{ UInt64() (uint64, error) },
	nonceDataMode interface{ UInt64() (uint64, error) },
	additionalData []*felt.Felt,
) (*felt.Felt, error) {
	if prefix == nil ||
		version == "" ||
		contractAddress == nil ||
		tip == nil ||
		resourceBounds == nil ||
		paymasterData == nil ||
		chainID == nil ||
		nonce == nil ||
		feeMode == nil ||
		nonceDataMode == nil ||
		additionalData == nil {
		return nil, ErrNotAllParametersSet
	}

	if isOrContainsNil(paymasterData, additionalData) {
		return nil, ErrNotAllParametersSet
	}

	versionFelt, err := new(felt.Felt).SetString(version)
	if err != nil {
		return nil, err
	}

	tipUint64, err := tip.ToUint64()
	if err != nil {
		return nil, err
	}
	tipAndResourceHash, err := tipAndResourcesHash(tipUint64, resourceBounds)
	if err != nil {
		return nil, err
	}

	DAUint64, err := dataAvailabilityModeConcat(feeMode, nonceDataMode)
	if err != nil {
		return nil, err
	}

	dataToHash := []*felt.Felt{
		prefix,
		versionFelt,
		contractAddress,
		tipAndResourceHash,
		curve.PoseidonArray(paymasterData...),
		chainID,
		nonce,
		felt.NewFromUint64[felt.Felt](DAUint64),
	}
	dataToHash = append(dataToHash, additionalData...)

	return curve.PoseidonArray(dataToHash...), nil
}

// @changed it's private now
// tipAndResourcesHash calculates the hash of the tip and resources.
func tipAndResourcesHash[
	R *rpcv9.ResourceBoundsMapping | *rpcv10.ResourceBoundsMapping,
](
	tip uint64,
	resourceBounds R,
) (*felt.Felt, error) {
	switch resourceBounds := any(resourceBounds).(type) {
	case *rpcv9.ResourceBoundsMapping:
		return tipAndResourcesHashInner(
			tip,
			resourceBounds.L1Gas,
			resourceBounds.L2Gas,
			resourceBounds.L1DataGas,
		)
	case *rpcv10.ResourceBoundsMapping:
		return tipAndResourcesHashInner(
			tip,
			resourceBounds.L1Gas,
			resourceBounds.L2Gas,
			resourceBounds.L1DataGas,
		)
	default:
		// should never happen due to generic type constraint
		return nil, errors.New("invalid resource bounds type")
	}
}

// tipAndResourcesHashInner calculates the hash of the tip and resources.
func tipAndResourcesHashInner[
	resource interface{ ~string },
	resourceBounds interface {
		Bytes(resource resource) ([]byte, error)
	},
](
	tip uint64,
	l1Gas, l2Gas, l1DataGas resourceBounds,
) (*felt.Felt, error) {
	l1Bytes, err := l1Gas.Bytes(resource(rpcv10.ResourceL1Gas))
	if err != nil {
		return nil, err
	}
	l2Bytes, err := l2Gas.Bytes(resource(rpcv10.ResourceL2Gas))
	if err != nil {
		return nil, err
	}
	l1DataGasBytes, err := l1DataGas.Bytes(resource(rpcv10.ResourceL1DataGas))
	if err != nil {
		return nil, err
	}

	l1Bounds := new(felt.Felt).SetBytes(l1Bytes)
	l2Bounds := new(felt.Felt).SetBytes(l2Bytes)
	l1DataGasBounds := new(felt.Felt).SetBytes(l1DataGasBytes)

	return curve.PoseidonArray(
		felt.NewFromUint64[felt.Felt](tip),
		l1Bounds,
		l2Bounds,
		l1DataGasBounds,
	), nil
}

// @changed it's private now
// dataAvailabilityModeConcat concatenates the data availability modes
// into a single uint64.
func dataAvailabilityModeConcat(feeDAMode, nonceDAMode interface{ UInt64() (uint64, error) }) (uint64, error) {
	const dataAvailabilityModeBits = 32
	fee64, err := feeDAMode.UInt64()
	if err != nil {
		return 0, err
	}
	nonce64, err := nonceDAMode.UInt64()
	if err != nil {
		return 0, err
	}

	return fee64 + nonce64<<dataAvailabilityModeBits, nil
}

// isOrContainsNil checks if any of the data is nil. If it's not nil and
// it's an array, it checks if any of the elements is nil.
func isOrContainsNil(data ...any) bool {
	for _, data := range data {
		if data == nil {
			return true
		}
		if dataArray, ok := data.([]any); ok {
			for _, d := range dataArray {
				if d == nil {
					return true
				}
			}
		}
	}
	return false
}
