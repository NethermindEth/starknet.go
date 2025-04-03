package abigen

const tmplCairoSimple = `
package {{.Package}}

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/abigen/accounts/abi"
	"github.com/NethermindEth/starknet.go/abigen/accounts/abi/bind"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/NethermindEth/starknet.go/utils"
)

var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = context.Background
	_ = bind.DeployContract
	_ = felt.NewFelt
	_ = utils.GetSelectorFromNameFelt
)

{{range $contract := .Contracts}}
type {{.Type}} struct {
	{{.Type}}Caller     // Read-only binding to the contract
	{{.Type}}Transactor // Write-only binding to the contract
	{{.Type}}Filterer   // Log filterer for contract events
}

type {{.Type}}Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

type {{.Type}}Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

type {{.Type}}Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

func New{{.Type}}(address *felt.Felt, backend bind.ContractBackend) (*{{.Type}}, error) {
	contract, err := bind{{.Type}}(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &{{.Type}}{
		{{.Type}}Caller:     {{.Type}}Caller{contract: contract},
		{{.Type}}Transactor: {{.Type}}Transactor{contract: contract},
		{{.Type}}Filterer:   {{.Type}}Filterer{contract: contract},
	}, nil
}

func bind{{.Type}}(address *felt.Felt, caller bind.ContractCaller, transactor bind.ContractTransact, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader({{.Type}}ABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

func (_{{.Type}} *{{.Type}}Caller) GetBalance(opts *bind.CallOpts) (*felt.Felt, error) {
	var out []interface{}
	err := _{{.Type}}.contract.Call(opts, &out, "get_balance")
	if err != nil {
		return nil, err
	}
	return out[0].(*felt.Felt), nil
}

func (_{{.Type}} *{{.Type}}Transactor) IncreaseBalance(opts *bind.TransactOpts, amount *felt.Felt) (*rpc.InvokeTxnResponse, error) {
	return _{{.Type}}.contract.Transact(opts, "increase_balance", amount)
}
{{end}}
`
