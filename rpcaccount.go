package caigo

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/dontpanicdao/caigo/rpc"
	"github.com/dontpanicdao/caigo/rpc/types"

	ctypes "github.com/dontpanicdao/caigo/types"
)

type account interface {
	Sign(msgHash *big.Int) (*big.Int, *big.Int, error)
	TransactionHash(calls []ctypes.FunctionCall, details ctypes.ExecuteDetails) (*big.Int, error)
	Call(ctx context.Context, call ctypes.FunctionCall) ([]string, error)
	Nonce(ctx context.Context) (*big.Int, error)
	EstimateFee(ctx context.Context, calls []ctypes.FunctionCall, details ctypes.ExecuteDetails) (*ctypes.FeeEstimate, error)
	Execute(ctx context.Context, calls []ctypes.FunctionCall, details ctypes.ExecuteDetails) (*ctypes.AddInvokeTransactionOutput, error)
}

var _ account = &RPCAccount{}

type RPCAccountPlugin interface {
	PluginCall(calls []ctypes.FunctionCall) (ctypes.FunctionCall, error)
}

type RPCAccount struct {
	Provider *rpc.Provider
	Address  string
	private  *big.Int
	version  *big.Int
	plugin   RPCAccountPlugin
}

type RPCAccountOption struct {
	RPCAccountPlugin RPCAccountPlugin
	version          *big.Int
}

type AccountOptionFunc func(string, string) (RPCAccountOption, error)

func AccountVersion0(string, string) (RPCAccountOption, error) {
	return RPCAccountOption{
		version: big.NewInt(0),
	}, nil
}

func AccountVersion1(string, string) (RPCAccountOption, error) {
	return RPCAccountOption{
		version: big.NewInt(1),
	}, nil
}

func NewRPCAccount(private, address string, provider *rpc.Provider, options ...AccountOptionFunc) (*RPCAccount, error) {
	var accountPlugin RPCAccountPlugin
	version := big.NewInt(0)
	for _, o := range options {
		opt, err := o(private, address)
		if err != nil {
			return nil, err
		}
		if opt.version != nil {
			version = opt.version
		}
		if opt.RPCAccountPlugin != nil {
			if accountPlugin != nil {
				return nil, errors.New("multiple plugins not supported")
			}
			accountPlugin = opt.RPCAccountPlugin
		}
	}
	if version.Cmp(big.NewInt(0)) != 0 {
		return nil, errors.New("account v1 not yet supported")
	}
	priv := ctypes.SNValToBN(private)

	return &RPCAccount{
		Provider: provider,
		Address:  address,
		private:  priv,
		version:  version,
		plugin:   accountPlugin,
	}, nil
}

func (account *RPCAccount) Call(ctx context.Context, call ctypes.FunctionCall) ([]string, error) {
	return account.Provider.Call(ctx, call, rpc.WithBlockTag("latest"))
}

func (account *RPCAccount) Sign(msgHash *big.Int) (*big.Int, *big.Int, error) {
	return Curve.Sign(msgHash, account.private)
}

func (account *RPCAccount) TransactionHash(calls []ctypes.FunctionCall, details ctypes.ExecuteDetails) (*big.Int, error) {
	chainID, err := account.Provider.ChainID(context.Background())
	if err != nil {
		return nil, err
	}

	var callArray []*big.Int
	switch {
	case account.version.Cmp(big.NewInt(0)) == 0:
		callArray = fmtV0Calldata(details.Nonce, calls)
	case account.version.Cmp(big.NewInt(1)) == 0:
		callArray = fmtCalldata(calls)
	default:
		return nil, fmt.Errorf("version %s unsupported", account.version.Text(10))
	}
	cdHash, err := Curve.ComputeHashOnElements(callArray)
	if err != nil {
		return nil, err
	}

	var multiHashData []*big.Int
	switch {
	case account.version.Cmp(big.NewInt(0)) == 0:
		multiHashData = []*big.Int{
			ctypes.UTF8StrToBig(TRANSACTION_PREFIX),
			account.version,
			ctypes.SNValToBN(account.Address),
			ctypes.GetSelectorFromName(EXECUTE_SELECTOR),
			cdHash,
			details.MaxFee,
			ctypes.UTF8StrToBig(chainID),
		}
	case account.version.Cmp(big.NewInt(1)) == 0:
		multiHashData = []*big.Int{
			ctypes.UTF8StrToBig(TRANSACTION_PREFIX),
			account.version,
			ctypes.SNValToBN(account.Address),
			big.NewInt(0),
			cdHash,
			details.MaxFee,
			ctypes.UTF8StrToBig(chainID),
			details.Nonce,
		}
	default:
		return nil, fmt.Errorf("version %s unsupported", account.version.Text(10))
	}
	return Curve.ComputeHashOnElements(multiHashData)
}

func (account *RPCAccount) Nonce(ctx context.Context) (*big.Int, error) {
	switch {
	case account.version.Cmp(big.NewInt(0)) == 0:
		nonce, err := account.Provider.Call(
			ctx,
			ctypes.FunctionCall{
				ContractAddress:    ctypes.HexToHash(account.Address),
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
	case account.version.Cmp(big.NewInt(1)) == 0:
		nonce, err := account.Provider.Nonce(
			ctx,
			ctypes.HexToHash(account.Address),
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
	return nil, fmt.Errorf("version %s unsupported", account.version.Text(10))
}

func (account *RPCAccount) EstimateFee(ctx context.Context, calls []ctypes.FunctionCall, details ctypes.ExecuteDetails) (*ctypes.FeeEstimate, error) {
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
	version := big.NewInt(0)
	if account.version != nil {
		version = account.version
	}
	if account.plugin != nil {
		call, err := account.plugin.PluginCall(calls)
		if err != nil {
			return nil, err
		}
		calls = append([]ctypes.FunctionCall{call}, calls...)
	}
	txHash, err := account.TransactionHash(
		calls,
		ctypes.ExecuteDetails{
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
	case account.version.Cmp(big.NewInt(0)) == 0:
		calldata = fmtV0CalldataStrings(nonce, calls)
	case account.version.Cmp(big.NewInt(1)) == 0:
		calldata = fmtCalldataStrings(calls)
	default:
		return nil, fmt.Errorf("version %s unsupported", account.version.Text(10))
	}
	accountDefaultV0Entrypoint := "__execute__"
	call := types.Call{
		MaxFee:             fmt.Sprintf("0x%s", maxFee.Text(16)),
		Version:            ctypes.NumAsHex(fmt.Sprintf("0x%s", version.Text(16))),
		Signature:          []string{fmt.Sprintf("0x%s", s1.Text(16)), fmt.Sprintf("0x%s", s2.Text(16))},
		ContractAddress:    ctypes.HexToHash(account.Address),
		EntryPointSelector: &accountDefaultV0Entrypoint,
		Calldata:           calldata,
	}
	return account.Provider.EstimateFee(ctx, call, rpc.WithBlockTag("latest"))
}

func (account *RPCAccount) Execute(ctx context.Context, calls []ctypes.FunctionCall, details ctypes.ExecuteDetails) (*ctypes.AddInvokeTransactionOutput, error) {
	if account.version != nil && account.version.Cmp(big.NewInt(0)) != 0 {
		return nil, errors.New("only invoke v0 is implemented")
	}
	var err error
	version := big.NewInt(0)
	if account.version != nil {
		version = account.version
	}
	nonce := details.Nonce
	if details.Nonce == nil {
		nonce, err = account.Nonce(ctx)
		if err != nil {
			return nil, err
		}
	}
	maxFee := details.MaxFee
	if details.MaxFee == nil {
		estimate, err := account.EstimateFee(ctx, calls, details)
		if err != nil {
			return nil, err
		}
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
		calls = append([]ctypes.FunctionCall{call}, calls...)
	}
	txHash, err := account.TransactionHash(
		calls,
		ctypes.ExecuteDetails{
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
	case account.version.Cmp(big.NewInt(0)) == 0:
		calldata = fmtV0CalldataStrings(nonce, calls)
	case account.version.Cmp(big.NewInt(1)) == 0:
		calldata = fmtCalldataStrings(calls)
	default:
		return nil, fmt.Errorf("version %s unsupported", account.version.Text(10))
	}
	// TODO: change this payload to manage both V0 and V1
	return account.Provider.AddInvokeTransaction(
		context.Background(),
		ctypes.FunctionCall{
			ContractAddress:    ctypes.HexToHash(account.Address),
			EntryPointSelector: "__execute__",
			Calldata:           calldata,
		},
		[]string{fmt.Sprintf("0x%s", s1.Text(16)), fmt.Sprintf("0x%s", s2.Text(16))},
		fmt.Sprintf("0x%s", maxFee.Text(16)),
		fmt.Sprintf("0x%s", version.Text(16)),
	)
}
