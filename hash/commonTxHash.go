package hash

import (
	"encoding/binary"
	"errors"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/curve"
	internalUtils "github.com/NethermindEth/starknet.go/internal/utils"
	"github.com/NethermindEth/starknet.go/rpc/types"
	"github.com/NethermindEth/starknet.go/types/constraints"
)

var (
	prefixInvoke        = new(felt.Felt).SetBytes([]byte("invoke"))
	prefixDeclare       = new(felt.Felt).SetBytes([]byte("declare"))
	prefixDeployAccount = new(felt.Felt).SetBytes([]byte("deploy_account"))
)

var (
	// @removed ErrFeltToBigInt
	ErrNotAllParametersSet = errors.New("not all necessary parameters have been set")
)

// CalculateDeprecatedTransactionHashCommon calculates the transaction hash
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
func CalculateDeprecatedTransactionHashCommon(
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

// @new
// CalculateV3TransactionHash calculates the hash of a V3 transaction;
// a common function to be used for all V3 transactions.
// Parameters:
//   - prefix: The prefix of the transaction hash
//   - version: The version of the transaction
//   - contractAddress: An contract address. Its meaning depends on the transaction type.
//   - tip: The tip for the transaction
//   - resourceBounds: The resource bounds for the transaction
//   - paymasterData: The paymaster data for the transaction
//   - chainID: The ID of the blockchain
//   - nonce: The nonce for the transaction
//   - feeMode: The fee mode for the transaction
//   - nonceDataMode: The nonce data mode for the transaction
//   - additionalData: Additional data to be included in the hash. Each tx type has its own additional data.
//
// Returns:
//   - *felt.Felt: the calculated transaction hash.
//   - error: an error if any.
func CalculateV3TransactionHash[
	u64 constraints.U64,
	u128 constraints.U128,
	RB constraints.ResourceBounds[u64, u128],
	RBM constraints.ResourceBoundsMapping[u64, u128, RB],
](
	prefix *felt.Felt,
	version string,
	contractAddress *felt.Felt,
	tip interface{ ToUint64() (uint64, error) },
	resourceBounds *RBM,
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

	innerResourceBounds := constraints.ResourceBoundsMappingImpl[u64, u128, RB](*resourceBounds)
	tipAndResourceHash, err := TipAndResourcesHash(tipUint64, &innerResourceBounds)
	if err != nil {
		return nil, err
	}

	DAUint64, err := DataAvailabilityModeConcat(feeMode, nonceDataMode)
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

// @changed accepts a generic type
// TipAndResourcesHash calculates the hash of the tip and resources.
// Parameters:
//   - tip: The tip for the transaction
//   - rbm: The resource bounds mapping for the transaction
//
// Returns:
//   - *felt.Felt: the calculated tip and resources hash.
//   - error: an error if any.
func TipAndResourcesHash[
	u64 constraints.U64,
	u128 constraints.U128,
	RB constraints.ResourceBounds[u64, u128],
](
	tip uint64,
	rbm *constraints.ResourceBoundsMappingImpl[u64, u128, RB],
) (*felt.Felt, error) {
	l1Bytes, err := ResourceBoundsBytes(&rbm.L1Gas, string(types.ResourceL1Gas))
	if err != nil {
		return nil, err
	}
	l2Bytes, err := ResourceBoundsBytes(&rbm.L2Gas, string(types.ResourceL2Gas))
	if err != nil {
		return nil, err
	}
	l1DataGasBytes, err := ResourceBoundsBytes(&rbm.L1DataGas, string(types.ResourceL1DataGas))
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

// @new
// ResourceBoundsBytes converts the resource bounds to a byte format
// necessary for the hash calculation.
// Parameters:
//   - rb: The resource bounds to convert
//   - resource: The resource type name
//
// Returns:
//   - []byte: the converted resource bounds
//   - error: an error if any.
func ResourceBoundsBytes[
	u64 constraints.U64,
	u128 constraints.U128,
	RB constraints.ResourceBounds[u64, u128],
](rb *RB, resource string) ([]byte, error) {
	if rb == nil {
		return nil, errors.New("resource bounds is nil")
	}
	innerRb := constraints.ResourceBoundsImpl[u64, u128](*rb)

	const eight = 8
	maxAmountBytes := make([]byte, eight)
	maxAmountUint64, err := innerRb.MaxAmount.ToUint64()
	if err != nil {
		return nil, err
	}
	binary.BigEndian.PutUint64(maxAmountBytes, maxAmountUint64)
	maxPricePerUnitFelt, err := new(felt.Felt).SetString(string(innerRb.MaxPricePerUnit))
	if err != nil {
		return nil, err
	}
	maxPriceBytes := maxPricePerUnitFelt.Bytes()

	return internalUtils.Flatten(
		[]byte{0},
		[]byte(resource),
		maxAmountBytes,
		maxPriceBytes[16:], // uint128.
	), nil
}

// @changed accepts an interface
// DataAvailabilityModeConcat concatenates the data availability modes
// into a single uint64.
// Parameters:
//   - feeDAMode: The fee data availability mode
//   - nonceDAMode: The nonce data availability mode
//
// Returns:
//   - uint64: the concatenated data availability modes
//   - error: an error if any.
func DataAvailabilityModeConcat(feeDAMode, nonceDAMode interface{ UInt64() (uint64, error) }) (uint64, error) {
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
