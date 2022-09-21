package rpc

import (
	"context"
	"fmt"
	"math/big"

	"github.com/dontpanicdao/caigo"
	"github.com/dontpanicdao/caigo/rpc/types"
)

const (
	EXECUTE_SELECTOR   string = "__execute__"
	TRANSACTION_PREFIX string = "invoke"
)

type Account struct {
	Provider *Client
	Address  string
	private  *big.Int
}

type ExecuteDetails struct {
	MaxFee  *big.Int
	Nonce   *big.Int
	Version *big.Int
}

func (provider *Client) NewAccount(private, address string) (*Account, error) {
	priv := caigo.SNValToBN(private)

	return &Account{
		Provider: provider,
		Address:  address,
		private:  priv,
	}, nil
}

func (account *Account) Sign(msgHash *big.Int) (*big.Int, *big.Int, error) {
	return caigo.Curve.Sign(msgHash, account.private)
}

func (account *Account) HashMultiCall(calls []types.FunctionCall, details ExecuteDetails) (*big.Int, error) {
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
		details.Version,
		caigo.SNValToBN(account.Address),
		caigo.GetSelectorFromName(EXECUTE_SELECTOR),
		cdHash,
		details.MaxFee,
		caigo.UTF8StrToBig(chainID),
	}

	return caigo.Curve.ComputeHashOnElements(multiHashData)
}

func fmtExecuteCalldataStrings(nonce *big.Int, calls []types.FunctionCall) (calldataStrings []string) {
	callArray := fmtExecuteCalldata(nonce, calls)
	for _, data := range callArray {
		calldataStrings = append(calldataStrings, fmt.Sprintf("0x%s", data.Text(16)))
	}
	return calldataStrings
}

/*
Formats the multicall transactions in a format which can be signed and verified by the network and OpenZeppelin account contracts
*/
func fmtExecuteCalldata(nonce *big.Int, calls []types.FunctionCall) (calldataArray []*big.Int) {
	callArray := []*big.Int{big.NewInt(int64(len(calls)))}

	for _, tx := range calls {
		address, _ := big.NewInt(0).SetString(tx.ContractAddress.Hex(), 0)
		callArray = append(callArray, address, caigo.GetSelectorFromName(tx.EntryPointSelector))

		if len(tx.CallData) == 0 {
			callArray = append(callArray, big.NewInt(0), big.NewInt(0))

			continue
		}

		callArray = append(callArray, big.NewInt(int64(len(calldataArray))), big.NewInt(int64(len(tx.CallData))))
		for _, cd := range tx.CallData {
			calldataArray = append(calldataArray, caigo.SNValToBN(cd))
		}
	}

	callArray = append(callArray, big.NewInt(int64(len(calldataArray))))
	callArray = append(callArray, calldataArray...)
	callArray = append(callArray, nonce)
	return callArray
}
