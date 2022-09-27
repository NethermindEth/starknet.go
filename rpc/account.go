package rpc

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/dontpanicdao/caigo"
	"github.com/dontpanicdao/caigo/rpc/types"
)

const (
	EXECUTE_SELECTOR   string = "__execute__"
	TRANSACTION_PREFIX string = "invoke"
)

type account interface {
	Sign(msgHash *big.Int) (*big.Int, *big.Int, error)
	TransactionHash(calls []types.FunctionCall, details types.ExecuteDetails) (*big.Int, error)
	Call(ctx context.Context, call types.FunctionCall) ([]string, error)
	Nonce(ctx context.Context) (*big.Int, error)
	EstimateFee(ctx context.Context, calls []types.FunctionCall, details types.ExecuteDetails) (*types.FeeEstimate, error)
	Execute(ctx context.Context, calls []types.FunctionCall, details types.ExecuteDetails) (*types.AddInvokeTransactionOutput, error)
}

var _ account = &Account{}

type AccountPlugin interface {
	PluginCall(calls []types.FunctionCall) (types.FunctionCall, error)
}

type Account struct {
	Provider *Provider
	Address  string
	private  *big.Int
	version  *big.Int
	plugin   AccountPlugin
}

type AccountOption struct {
	AccountPlugin AccountPlugin
	version       *big.Int
}

type AccountOptionFunc func() (AccountOption, error)

func AccountVersion0() (AccountOption, error) {
	return AccountOption{
		version: big.NewInt(0),
	}, nil
}

func AccountVersion1() (AccountOption, error) {
	return AccountOption{
		version: big.NewInt(1),
	}, nil
}

func (provider *Provider) NewAccount(private, address string, options ...AccountOptionFunc) (*Account, error) {
	var accountPlugin AccountPlugin
	version := big.NewInt(0)
	for _, o := range options {
		opt, err := o()
		if err != nil {
			return nil, err
		}
		if opt.version != nil {
			version = opt.version
		}
		if opt.AccountPlugin != nil {
			if accountPlugin != nil {
				return nil, errors.New("multiple plugins not supported")
			}
			accountPlugin = opt.AccountPlugin
		}
	}
	if version.Cmp(big.NewInt(0)) != 0 {
		return nil, errors.New("account v1 not yet supported")
	}
	priv := caigo.SNValToBN(private)

	return &Account{
		Provider: provider,
		Address:  address,
		private:  priv,
		version:  version,
		plugin:   accountPlugin,
	}, nil
}

func (account *Account) Call(ctx context.Context, call types.FunctionCall) ([]string, error) {
	return account.Provider.Call(ctx, call, WithBlockTag("latest"))
}

func (account *Account) Sign(msgHash *big.Int) (*big.Int, *big.Int, error) {
	return caigo.Curve.Sign(msgHash, account.private)
}

func (account *Account) TransactionHash(calls []types.FunctionCall, details types.ExecuteDetails) (*big.Int, error) {
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
	cdHash, err := caigo.Curve.ComputeHashOnElements(callArray)
	if err != nil {
		return nil, err
	}

	var multiHashData []*big.Int
	switch {
	case account.version.Cmp(big.NewInt(0)) == 0:
		multiHashData = []*big.Int{
			caigo.UTF8StrToBig(TRANSACTION_PREFIX),
			account.version,
			caigo.SNValToBN(account.Address),
			caigo.GetSelectorFromName(EXECUTE_SELECTOR),
			cdHash,
			details.MaxFee,
			caigo.UTF8StrToBig(chainID),
		}
	case account.version.Cmp(big.NewInt(1)) == 0:
		multiHashData = []*big.Int{
			caigo.UTF8StrToBig(TRANSACTION_PREFIX),
			account.version,
			caigo.SNValToBN(account.Address),
			cdHash,
			details.MaxFee,
			details.Nonce,
			caigo.UTF8StrToBig(chainID),
		}
	default:
		return nil, fmt.Errorf("version %s unsupported", account.version.Text(10))
	}
	return caigo.Curve.ComputeHashOnElements(multiHashData)
}

func (account *Account) Nonce(ctx context.Context) (*big.Int, error) {
	switch {
	case account.version.Cmp(big.NewInt(0)) == 0:
		nonce, err := account.Provider.Call(
			ctx,
			types.FunctionCall{
				ContractAddress:    types.HexToHash(account.Address),
				EntryPointSelector: "get_nonce",
				CallData:           []string{},
			},
			WithBlockTag("latest"),
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
	return nil, fmt.Errorf("version %s unsupported", account.version.Text(10))
}

func (account *Account) EstimateFee(ctx context.Context, calls []types.FunctionCall, details types.ExecuteDetails) (*types.FeeEstimate, error) {
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
		Version:            types.NumAsHex(fmt.Sprintf("0x%s", version.Text(16))),
		Signature:          []string{fmt.Sprintf("0x%s", s1.Text(16)), fmt.Sprintf("0x%s", s2.Text(16))},
		ContractAddress:    types.HexToHash(account.Address),
		EntryPointSelector: &accountDefaultV0Entrypoint,
		CallData:           calldata,
	}
	return account.Provider.EstimateFee(ctx, call, WithBlockTag("latest"))
}

func (account *Account) Execute(ctx context.Context, calls []types.FunctionCall, details types.ExecuteDetails) (*types.AddInvokeTransactionOutput, error) {
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
		types.FunctionCall{
			ContractAddress:    types.HexToHash(account.Address),
			EntryPointSelector: "__execute__",
			CallData:           calldata,
		},
		[]string{fmt.Sprintf("0x%s", s1.Text(16)), fmt.Sprintf("0x%s", s2.Text(16))},
		fmt.Sprintf("0x%s", maxFee.Text(16)),
		fmt.Sprintf("0x%s", version.Text(16)),
	)
}
