package caigo

import (
	"context"
	"fmt"
	"math/big"

	"github.com/dontpanicdao/caigo/felt"
	"github.com/dontpanicdao/caigo/types"
)

const (
	EXECUTE_SELECTOR    string = "__execute__"
	TRANSACTION_PREFIX  string = "invoke"
	TRANSACTION_VERSION int64  = 0
	FEE_MARGIN          uint64 = 115
)

type Account struct {
	Provider types.Provider
	Address  felt.Felt
	PublicX  *big.Int
	PublicY  *big.Int
	private  *big.Int
}

type ExecuteDetails struct {
	MaxFee  *felt.Felt
	Nonce   *felt.Felt
	Version *felt.Felt
}

/*
Instantiate a new StarkNet Account which includes structures for calling the network and signing transactions:
- private signing key
- stark curve definition
- full provider definition
- public key pair for signature verifications
*/
func NewAccount(private string, address felt.Felt, provider types.Provider) (*Account, error) {
	priv, ok := big.NewInt(0).SetString(private, 0)
	if !ok {
		return nil, fmt.Errorf("wrongPrivate")
	}
	x, y, err := felt.GetCurve().PrivateToPoint(priv)
	if err != nil {
		return nil, err
	}
	return &Account{
		Provider: provider,
		Address:  address,
		PublicX:  x,
		PublicY:  y,
		private:  priv,
	}, nil
}

func (account *Account) Sign(msgHash *big.Int) (*big.Int, *big.Int, error) {
	return felt.GetCurve().Sign(msgHash, account.private)
}

/*
invocation wrapper for StarkNet account calls to '__execute__' contact calls through an account abstraction
- implementation has been tested against OpenZeppelin Account contract as of: https://github.com/OpenZeppelin/cairo-contracts/blob/4116c1ecbed9f821a2aa714c993a35c1682c946e/src/openzeppelin/account/Account.cairo
- accepts a multicall
*/
func (account *Account) Execute(ctx context.Context, calls []types.Transaction, details ExecuteDetails) (*types.AddTxResponse, error) {
	if details.Nonce == nil {
		nonce, err := account.Provider.AccountNonce(ctx, account.Address)
		if err != nil {
			return nil, err
		}
		details.Nonce = nonce
	}

	if details.MaxFee == nil {
		fee, err := account.EstimateFee(ctx, calls, details)
		if err != nil {
			return nil, err
		}
		details.MaxFee = &types.Felt{
			Int: new(big.Int).SetUint64((fee.OverallFee * FEE_MARGIN) / 100),
		}
	}

	req, err := account.fmtExecute(ctx, calls, details)
	if err != nil {
		return nil, err
	}

	return account.Provider.Invoke(ctx, *req)
}

func (account *Account) HashMultiCall(fee *types.Felt, nonce *types.Felt, calls []types.Transaction) (*big.Int, error) {
	chainID, err := account.Provider.ChainID(context.Background())
	if err != nil {
		return nil, err
	}

	callArray := ExecuteCalldata(nonce, calls)

	// convert callArray into a BigInt array
	callArrayBigInt := make([]*big.Int, 0)
	for _, call := range callArray {
		callArrayBigInt = append(callArrayBigInt, call.Int)
	}

	cdHash, err := Curve.ComputeHashOnElements(callArrayBigInt)
	if err != nil {
		return nil, err
	}

	multiHashData := []big.Felt{
		felt.UTF8StrToFelt(TRANSACTION_PREFIX),
		big.NewInt(TRANSACTION_VERSION),
		SNValToBN(account.Address.String()),
		GetSelectorFromName(EXECUTE_SELECTOR),
		cdHash,
		fee.Int,
		UTF8StrToBig(chainID),
	}

	return Curve.ComputeHashOnElements(multiHashData)
}

func (account *Account) EstimateFee(ctx context.Context, calls []types.Transaction, details ExecuteDetails) (*types.FeeEstimate, error) {
	if details.Nonce == nil {
		nonce, err := account.Provider.AccountNonce(ctx, account.Address)
		if err != nil {
			return nil, err
		}
		details.Nonce = nonce
	}

	if details.MaxFee == nil {
		details.MaxFee = &types.Felt{Int: big.NewInt(0)}
	}

	req, err := account.fmtExecute(ctx, calls, details)
	if err != nil {
		return nil, err
	}

	return account.Provider.EstimateFee(ctx, *req, "")
}

func (account *Account) fmtExecute(ctx context.Context, calls []types.Transaction, details ExecuteDetails) (*types.FunctionInvoke, error) {
	req := types.FunctionInvoke{
		FunctionCall: types.FunctionCall{
			ContractAddress:    account.Address,
			EntryPointSelector: GetSelectorFromName(EXECUTE_SELECTOR),
			Calldata:           ExecuteCalldata(details.Nonce, calls),
		},
		MaxFee: details.MaxFee,
	}

	hash, err := account.HashMultiCall(details.MaxFee, details.Nonce, calls)
	if err != nil {
		return nil, err
	}

	r, s, err := account.Sign(hash)
	if err != nil {
		return nil, err
	}
	req.Signature = types.Signature{types.BigToFelt(r), types.BigToFelt(s)}

	return &req, nil
}

/*
Formats the multicall transactions in a format which can be signed and verified by the network and OpenZeppelin account contracts
*/
func ExecuteCalldata(nonce *types.Felt, calls []types.Transaction) (calldataArray []*types.Felt) {
	callArray := []*types.Felt{types.BigToFelt(big.NewInt(int64(len(calls))))}

	for _, tx := range calls {
		callArray = append(callArray, types.BigToFelt(SNValToBN(tx.ContractAddress.String())), types.BigToFelt(GetSelectorFromName(tx.EntryPointSelector)))

		if len(tx.Calldata) == 0 {
			callArray = append(callArray, types.BigToFelt(big.NewInt(0)), types.BigToFelt(big.NewInt(0)))

			continue
		}

		callArray = append(callArray, types.BigToFelt(big.NewInt(int64(len(calldataArray)))), types.BigToFelt(big.NewInt(int64(len(tx.Calldata)))))
		for _, cd := range tx.Calldata {
			calldataArray = append(calldataArray, types.BigToFelt(SNValToBN(cd.String())))
		}
	}

	callArray = append(callArray, types.BigToFelt(big.NewInt(int64(len(calldataArray)))))
	callArray = append(callArray, calldataArray...)
	callArray = append(callArray, nonce)
	return callArray
}
