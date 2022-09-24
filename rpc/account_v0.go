package rpc

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/dontpanicdao/caigo"
	"github.com/dontpanicdao/caigo/rpc/types"
)

type AccountV0 struct {
	Provider *Provider
	Address  string
	private  *big.Int
	version  *big.Int
}

func (provider *Provider) NewAccountV0(private, address string) (*AccountV0, error) {
	priv := caigo.SNValToBN(private)

	return &AccountV0{
		Provider: provider,
		Address:  address,
		private:  priv,
		version:  big.NewInt(0),
	}, nil
}

func (account *AccountV0) Sign(msgHash *big.Int) (*big.Int, *big.Int, error) {
	return caigo.Curve.Sign(msgHash, account.private)
}

func (account *AccountV0) HashMultiCall(calls []types.FunctionCall, details types.ExecuteDetails) (*big.Int, error) {
	chainID, err := account.Provider.ChainID(context.Background())
	if err != nil {
		return nil, err
	}

	callArray := fmtExecuteCalldata(details.Nonce, calls)
	cdHash, err := caigo.Curve.ComputeHashOnElements(callArray)
	if err != nil {
		return nil, err
	}

	multiHashData := []*big.Int{
		caigo.UTF8StrToBig(TRANSACTION_PREFIX),
		account.version,
		caigo.SNValToBN(account.Address),
		caigo.GetSelectorFromName(EXECUTE_SELECTOR),
		cdHash,
		details.MaxFee,
		caigo.UTF8StrToBig(chainID),
	}

	return caigo.Curve.ComputeHashOnElements(multiHashData)
}

func (account *AccountV0) Nonce(ctx context.Context) (*big.Int, error) {
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
}

func (account *AccountV0) EstimateFee(ctx context.Context, calls []types.FunctionCall, details types.ExecuteDetails) (*types.FeeEstimate, error) {
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
	txHash, err := account.HashMultiCall(
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
	calldata := fmtExecuteCalldataStrings(nonce, calls)
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

func (account *AccountV0) Execute(ctx context.Context, calls []types.FunctionCall, details types.ExecuteDetails) (*types.AddInvokeTransactionOutput, error) {
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
	txHash, err := account.HashMultiCall(
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
	calldata := fmtExecuteCalldataStrings(nonce, calls)
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
