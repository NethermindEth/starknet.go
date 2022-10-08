package caigo

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/dontpanicdao/caigo/rpc"

	 "github.com/dontpanicdao/caigo/types"
)

type account interface {
	Sign(msgHash *big.Int) (*big.Int, *big.Int, error)
	TransactionHash(calls []types.FunctionCall, details types.ExecuteDetails) (*big.Int, error)
	Call(ctx context.Context, call types.FunctionCall) ([]string, error)
	Nonce(ctx context.Context) (*big.Int, error)
	EstimateFee(ctx context.Context, calls []types.FunctionCall, details types.ExecuteDetails) (*types.FeeEstimate, error)
	Execute(ctx context.Context, calls []types.FunctionCall, details types.ExecuteDetails) (*types.AddInvokeTransactionOutput, error)
}

var _ account = &RPCAccount{}

type RPCAccountPlugin interface {
	PluginCall(calls []types.FunctionCall) (types.FunctionCall, error)
}

type RPCAccount struct {
	Provider *rpc.Provider
	Address  string
	private  *big.Int
	version  uint64
	plugin   RPCAccountPlugin
}

type RPCAccountOption struct {
	RPCAccountPlugin RPCAccountPlugin
	version          uint64
}

type AccountOptionFunc func(string, string) (RPCAccountOption, error)

func AccountVersion0(string, string) (RPCAccountOption, error) {
	return RPCAccountOption{
		version: uint64(0),
	}, nil
}

func AccountVersion1(string, string) (RPCAccountOption, error) {
	return RPCAccountOption{
		version: uint64(1),
	}, nil
}

func NewRPCAccount(private, address string, provider *rpc.Provider, options ...AccountOptionFunc) (*RPCAccount, error) {
	var accountPlugin RPCAccountPlugin
	version := uint64(0)
	for _, o := range options {
		opt, err := o(private, address)
		if err != nil {
			return nil, err
		}
		if opt.version != 0 {
			version = opt.version
		}
		if opt.RPCAccountPlugin != nil {
			if accountPlugin != nil {
				return nil, errors.New("multiple plugins not supported")
			}
			accountPlugin = opt.RPCAccountPlugin
		}
	}
	if version != 0 {
		return nil, errors.New("account v1 not yet supported")
	}
	priv := types.SNValToBN(private)

	return &RPCAccount{
		Provider: provider,
		Address:  address,
		private:  priv,
		version:  version,
		plugin:   accountPlugin,
	}, nil
}

func (account *RPCAccount) Call(ctx context.Context, call types.FunctionCall) ([]string, error) {
	return account.Provider.Call(ctx, call, rpc.WithBlockTag("latest"))
}

func (account *RPCAccount) Sign(msgHash *big.Int) (*big.Int, *big.Int, error) {
	return Curve.Sign(msgHash, account.private)
}

func (account *RPCAccount) TransactionHash(calls []types.FunctionCall, details types.ExecuteDetails) (*big.Int, error) {
	chainID, err := account.Provider.ChainID(context.Background())
	if err != nil {
		return nil, err
	}

	var callArray []*big.Int
	switch {
	case account.version == 0:
		callArray = fmtV0Calldata(details.Nonce, calls)
	case account.version == 1:
		callArray = fmtCalldata(calls)
	default:
		return nil, fmt.Errorf("version %d unsupported", account.version)
	}
	cdHash, err := Curve.ComputeHashOnElements(callArray)
	if err != nil {
		return nil, err
	}

	var multiHashData []*big.Int
	switch {
	case account.version == 0:
		multiHashData = []*big.Int{
			types.UTF8StrToBig(TRANSACTION_PREFIX),
			big.NewInt(int64(account.version)),
			types.SNValToBN(account.Address),
			types.GetSelectorFromName(EXECUTE_SELECTOR),
			cdHash,
			details.MaxFee,
			types.UTF8StrToBig(chainID),
		}
	case account.version == 1:
		multiHashData = []*big.Int{
			types.UTF8StrToBig(TRANSACTION_PREFIX),
			big.NewInt(int64(account.version)),
			types.SNValToBN(account.Address),
			big.NewInt(0),
			cdHash,
			details.MaxFee,
			types.UTF8StrToBig(chainID),
			details.Nonce,
		}
	default:
		return nil, fmt.Errorf("version %d unsupported", account.version)
	}
	return Curve.ComputeHashOnElements(multiHashData)
}

func (account *RPCAccount) Nonce(ctx context.Context) (*big.Int, error) {
	switch {
	case account.version == 0:
		nonce, err := account.Provider.Call(
			ctx,
			types.FunctionCall{
				ContractAddress:    types.HexToHash(account.Address),
				EntryPointSelector: "get_nonce",
				Calldata:           []string{},
			},
			rpc.WithBlockTag("latest"),
		)
		if err != nil {
			return nil, err
		}
		if len(nonce) == 0 {
			return nil, errors.New("nonce error")
		}
		n, ok := big.NewInt(0).SetString(nonce[0], 0)
		if !ok {
			return nil, errors.New("nonce error")
		}
		return n, nil
	case account.version == 1:
		nonce, err := account.Provider.Nonce(
			ctx,
			types.HexToHash(account.Address),
		)
		if err != nil {
			return nil, err
		}
		if nonce == nil {
			return nil, errors.New("nonce is nil")
		}
		n, ok := big.NewInt(0).SetString(*nonce, 0)
		if !ok {
			return nil, errors.New("nonce error")
		}
		return n, nil
	}
	return nil, fmt.Errorf("version %d unsupported", account.version)
}

func (account *RPCAccount) EstimateFee(ctx context.Context, calls []types.FunctionCall, details types.ExecuteDetails) (*types.FeeEstimate, error) {
	var err error
	nonce := details.Nonce
	if details.Nonce == nil {
		nonce, err = account.Nonce(ctx)
		if err != nil {
			return nil, err
		}
	}
	maxFee, _ := big.NewInt(0).SetString("0x200000000", 0)
	if details.MaxFee != nil {
		maxFee = details.MaxFee
	}
	version := account.version
	if account.plugin != nil {
		call, err := account.plugin.PluginCall(calls)
		if err != nil {
			return nil, err
		}
		calls = append([]types.FunctionCall{call}, calls...)
	}
	txHash, err := account.TransactionHash(
		calls,
		types.ExecuteDetails{
			Nonce:  nonce,
			MaxFee: maxFee,
		},
	)
	if err != nil {
		return nil, err
	}
	s1, s2, err := account.Sign(txHash)
	if err != nil {
		return nil, err
	}
	var calldata []string
	switch {
	case account.version == 0:
		calldata = fmtV0CalldataStrings(nonce, calls)
	case account.version == 1:
		calldata = fmtCalldataStrings(calls)
	default:
		return nil, fmt.Errorf("version %d unsupported", account.version)
	}
	accountDefaultV0Entrypoint := "__execute__"
	call := types.FunctionInvoke{
		MaxFee:    maxFee,
		Version:   version,
		Signature: types.Signature{s1, s2},
		FunctionCall: types.FunctionCall{
			ContractAddress:    types.HexToHash(account.Address),
			EntryPointSelector: accountDefaultV0Entrypoint,
			Calldata:           calldata,
		},
	}
	return account.Provider.EstimateFee(ctx, call, rpc.WithBlockTag("latest"))
}

func (account *RPCAccount) Execute(ctx context.Context, calls []types.FunctionCall, details types.ExecuteDetails) (*types.AddInvokeTransactionOutput, error) {
	if account.version != 0 {
		return nil, errors.New("only invoke v0 is implemented")
	}
	var err error
	version := account.version
	nonce := details.Nonce
	if details.Nonce == nil {
		nonce, err = account.Nonce(ctx)
		if err != nil {
			return nil, err
		}
	}
	fmt.Println("Got Nonce", details.Nonce.Text(10))
	maxFee := details.MaxFee
	if details.MaxFee == nil {
		fmt.Println("Will run estimateFee")
		estimate, err := account.EstimateFee(ctx, calls, details)
		if err != nil {
			return nil, err
		}
		fmt.Println("EstimateFee executed")
		v, ok := big.NewInt(0).SetString(string(estimate.OverallFee), 0)
		if !ok {
			return nil, errors.New("could not match OverallFee to big.Int")
		}
		maxFee = v.Mul(v, big.NewInt(2))
	}
	if account.plugin != nil {
		call, err := account.plugin.PluginCall(calls)
		if err != nil {
			return nil, err
		}
		calls = append([]types.FunctionCall{call}, calls...)
	}
	txHash, err := account.TransactionHash(
		calls,
		types.ExecuteDetails{
			Nonce:  nonce,
			MaxFee: maxFee,
		},
	)
	if err != nil {
		return nil, err
	}
	s1, s2, err := account.Sign(txHash)
	if err != nil {
		return nil, err
	}
	var calldata []string
	switch {
	case account.version == 0:
		calldata = fmtV0CalldataStrings(nonce, calls)
	case account.version == 1:
		calldata = fmtCalldataStrings(calls)
	default:
		return nil, fmt.Errorf("version %d unsupported", account.version)
	}
	// TODO: change this payload to manage both V0 and V1
	return account.Provider.AddInvokeTransaction(
		context.Background(),
		types.FunctionCall{
			ContractAddress:    types.HexToHash(account.Address),
			EntryPointSelector: "__execute__",
			Calldata:           calldata,
		},
		[]string{fmt.Sprintf("0x%s", s1.Text(16)), fmt.Sprintf("0x%s", s2.Text(16))},
		fmt.Sprintf("0x%s", maxFee.Text(16)),
		fmt.Sprintf("0x%d", version),
	)
}
