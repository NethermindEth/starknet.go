package account

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/rpc/rpcv10"
	"github.com/NethermindEth/starknet.go/rpc/rpcv9"
	"github.com/NethermindEth/starknet.go/rpc/types"
)

// @todo make everything private
// add tests where needed.

type providerWrapper struct {
	chainID string
	version RPCVersion

	rpcv9  rpcv9.RPCProvider
	rpcv10 rpcv10.RPCProvider
}

// @new
type RPCProvider interface {
	*rpcv10.Provider | *rpcv9.Provider
}

func newWrapperFrom[P RPCProvider](provider P) *providerWrapper {
	var wrapper providerWrapper

	switch p := any(provider).(type) {
	case *rpcv9.Provider:
		wrapper.rpcv9 = p
		wrapper.version = RPCVersion9
	case *rpcv10.Provider:
		wrapper.rpcv10 = p
		wrapper.version = RPCVersion10
	}

	return &wrapper
}

type OtherMethods interface {
	AsV9() rpcv9.RPCProvider
	AsV10() rpcv10.RPCProvider
	EstimateTip(ctx context.Context, multiplier float64) (tip types.U64, err error)
	SendTransaction(ctx context.Context, txn types.BroadcastTxn) (types.TransactionResponse, error)
}

// implementing the methods
func (p *providerWrapper) ChainID(ctx context.Context) (string, error)
func (p *providerWrapper) EstimateFee(
	ctx context.Context,
	requests []types.BroadcastTxn,
	simulationFlags []types.SimulationFlag,
	blockID types.BlockID,
) ([]types.FeeEstimation, error)
func (p *providerWrapper) Nonce(
	ctx context.Context,
	blockID types.BlockID,
	contractAddress *felt.Felt,
) (*felt.Felt, error)
func (p *providerWrapper) TransactionByHash(ctx context.Context, hash *felt.Felt) (types.BlockTransaction, error)
func (p *providerWrapper) TransactionReceipt(
	ctx context.Context,
	transactionHash *felt.Felt,
) (types.TransactionReceiptWithBlockInfo, error)

// func (p *providerWrapper) TransactionStatus(
// 	ctx context.Context,
// 	transactionHash *felt.Felt,
// ) (types.TxnStatusResult, error)

func (p *providerWrapper) EstimateTip(ctx context.Context, multiplier float64) (tip types.U64, err error)
func (p *providerWrapper) SendTransaction(ctx context.Context, txn types.BroadcastTxn) (types.TransactionResponse, error)
func (p *providerWrapper) AsV9() rpcv9.RPCProvider {
	return p.rpcv9
}

func (p *providerWrapper) AsV10() rpcv10.RPCProvider {
	return p.rpcv10
}

func (p *providerWrapper) Version() RPCVersion {
	return p.version
}

// @todo add tests for this type

type RPCVersion int

const (
	RPCVersion9 RPCVersion = iota
	RPCVersion10
)

func (v RPCVersion) String() string {
	switch v {
	case RPCVersion9:
		return "0.9.0"
	case RPCVersion10:
		return "0.10.1"
	}
	return ""
}

func (v RPCVersion) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.String())
}

func (v *RPCVersion) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	semversion, err := semver.NewVersion(s)
	if err != nil {
		return fmt.Errorf("failed to parse version to semver: %w", err)
	}

	switch {
	case semversion.Compare(semver.MustParse(RPCVersion9.String())) == 0:
		*v = RPCVersion9
	case semversion.Compare(semver.MustParse(RPCVersion10.String())) == 0:
		*v = RPCVersion10
	default:
		return fmt.Errorf("invalid RPC version: %s", s)
	}
	return nil
}
