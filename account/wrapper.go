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

//go:generate mockgen -destination=./wrapper_mock_test.go -package=account -source=wrapper.go providerWrapperI
type providerWrapper interface {
	ChainID(ctx context.Context) (string, error)
	EstimateFee(
		ctx context.Context,
		requests []types.BroadcastTxn,
		simulationFlags []types.SimulationFlag,
		blockID types.BlockID,
	) ([]types.FeeEstimation, error)
	Nonce(ctx context.Context, blockID types.BlockID, contractAddress *felt.Felt) (*felt.Felt, error)
	TransactionByHash(ctx context.Context, hash *felt.Felt) (types.BlockTransaction, error)
	TransactionReceipt(ctx context.Context, transactionHash *felt.Felt) (types.TransactionReceiptWithBlockInfo, error)
	EstimateTip(ctx context.Context, multiplier float64) (tip types.U64, err error)
	SendTransaction(ctx context.Context, txn types.BroadcastTxn) (types.TransactionResponse, error)
	AsV9() rpcv9.RPCProvider
	AsV10() rpcv10.RPCProvider
	Version() RPCVersion
}

type wrapper struct {
	chainID string
	version RPCVersion

	rpcv9  rpcv9.RPCProvider
	rpcv10 rpcv10.RPCProvider
}

var _ providerWrapper = (*wrapper)(nil)

// @new
type _RPCProvider interface {
	*rpcv10.Provider | *rpcv9.Provider
}

func newWrapperFrom[P _RPCProvider](provider P) providerWrapper {
	var wrapper wrapper

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

// implementing the methods
func (p *wrapper) ChainID(ctx context.Context) (string, error)
func (p *wrapper) EstimateFee(
	ctx context.Context,
	requests []types.BroadcastTxn,
	simulationFlags []types.SimulationFlag,
	blockID types.BlockID,
) ([]types.FeeEstimation, error)
func (p *wrapper) Nonce(
	ctx context.Context,
	blockID types.BlockID,
	contractAddress *felt.Felt,
) (*felt.Felt, error)
func (p *wrapper) TransactionByHash(ctx context.Context, hash *felt.Felt) (types.BlockTransaction, error)
func (p *wrapper) TransactionReceipt(
	ctx context.Context,
	transactionHash *felt.Felt,
) (types.TransactionReceiptWithBlockInfo, error)

// func (p *providerWrapper) TransactionStatus(
// 	ctx context.Context,
// 	transactionHash *felt.Felt,
// ) (types.TxnStatusResult, error)

func (p *wrapper) EstimateTip(ctx context.Context, multiplier float64) (tip types.U64, err error)
func (p *wrapper) SendTransaction(ctx context.Context, txn types.BroadcastTxn) (types.TransactionResponse, error)
func (p *wrapper) AsV9() rpcv9.RPCProvider {
	return p.rpcv9
}

func (p *wrapper) AsV10() rpcv10.RPCProvider {
	return p.rpcv10
}

func (p *wrapper) Version() RPCVersion {
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
