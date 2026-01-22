package hash

import (
	"errors"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/curve"
	"github.com/NethermindEth/starknet.go/rpc/rpcv10"
	"github.com/NethermindEth/starknet.go/rpc/rpcv9"
)

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
	accountDeploymentData []*felt.Felt,
	calldata []*felt.Felt,
) (*felt.Felt, error) {
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

	return curve.PoseidonArray(
		prefix,
		versionFelt,
		contractAddress,
		tipAndResourceHash,
		curve.PoseidonArray(paymasterData...),
		chainID,
		nonce,
		// dataAvailabilityMode,
		// curve.PoseidonArray(accountDeploymentData...),
		// curve.PoseidonArray(calldata...),
	), nil
}

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
