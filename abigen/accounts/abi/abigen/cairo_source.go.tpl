// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

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

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = context.Background
	_ = bind.DeployContract
	_ = felt.NewFelt
	_ = utils.GetSelectorFromNameFelt
)

{{$structs := .Structs}}
{{range $structs}}
// {{.Name}} is an auto generated low-level Go binding around a Cairo struct.
type {{.Name}} struct {
	{{range $field := .Fields}}
	{{$field.Name}} {{$field.Type}}{{end}}
}
{{end}}

{{range $contract := .Contracts}}
// {{.Type}}MetaData contains all metadata concerning the {{.Type}} contract.
var {{.Type}}MetaData = &bind.MetaData{
	ABI: "{{.InputABI}}",
	{{if .InputBin}}Bin: "{{.InputBin}}",{{end}}
}

// {{.Type}}ABI is the input ABI used to generate the binding from.
const {{.Type}}ABI = "{{.InputABI}}"

{{if .InputBin}}
// {{.Type}}Bin is the compiled bytecode used for deploying new contracts.
var {{.Type}}Bin = "{{.InputBin}}"
{{end}}

// {{.Type}} is an auto generated Go binding around a Cairo contract.
type {{.Type}} struct {
	{{.Type}}Caller     // Read-only binding to the contract
	{{.Type}}Transactor // Write-only binding to the contract
	{{.Type}}Filterer   // Log filterer for contract events
}

// {{.Type}}Caller is an auto generated read-only Go binding around a Cairo contract.
type {{.Type}}Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// {{.Type}}Transactor is an auto generated write-only Go binding around a Cairo contract.
type {{.Type}}Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// {{.Type}}Filterer is an auto generated log filtering Go binding around a Cairo contract events.
type {{.Type}}Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// {{.Type}}Session is an auto generated Go binding around a Cairo contract,
// with pre-set call and transact options.
type {{.Type}}Session struct {
	Contract     *{{.Type}}        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction options to use throughout this session
}

// {{.Type}}CallerSession is an auto generated read-only Go binding around a Cairo contract,
// with pre-set call options.
type {{.Type}}CallerSession struct {
	Contract *{{.Type}}Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// {{.Type}}TransactorSession is an auto generated write-only Go binding around a Cairo contract,
// with pre-set transact options.
type {{.Type}}TransactorSession struct {
	Contract     *{{.Type}}Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction options to use throughout this session
}

// {{.Type}}Raw is an auto generated low-level Go binding around a Cairo contract.
type {{.Type}}Raw struct {
	Contract *{{.Type}} // Generic contract binding to access the raw methods on
}

// {{.Type}}CallerRaw is an auto generated low-level read-only Go binding around a Cairo contract.
type {{.Type}}CallerRaw struct {
	Contract *{{.Type}}Caller // Generic read-only contract binding to access the raw methods on
}

// {{.Type}}TransactorRaw is an auto generated low-level write-only Go binding around a Cairo contract.
type {{.Type}}TransactorRaw struct {
	Contract *{{.Type}}Transactor // Generic write-only contract binding to access the raw methods on
}

// New{{.Type}} creates a new instance of {{.Type}}, bound to a specific deployed contract.
func New{{.Type}}(address *felt.Felt, backend bind.ContractBackend) (*{{.Type}}, error) {
	contract, err := bind{{.Type}}(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &{{.Type}}{{{.Type}}Caller: {{.Type}}Caller{contract: contract}, {{.Type}}Transactor: {{.Type}}Transactor{contract: contract}, {{.Type}}Filterer: {{.Type}}Filterer{contract: contract}}, nil
}

// New{{.Type}}Caller creates a new read-only instance of {{.Type}}, bound to a specific deployed contract.
func New{{.Type}}Caller(address *felt.Felt, caller bind.ContractCaller) (*{{.Type}}Caller, error) {
	contract, err := bind{{.Type}}(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &{{.Type}}Caller{contract: contract}, nil
}

// New{{.Type}}Transactor creates a new write-only instance of {{.Type}}, bound to a specific deployed contract.
func New{{.Type}}Transactor(address *felt.Felt, transactor bind.ContractTransact) (*{{.Type}}Transactor, error) {
	contract, err := bind{{.Type}}(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &{{.Type}}Transactor{contract: contract}, nil
}

// New{{.Type}}Filterer creates a new log filterer instance of {{.Type}}, bound to a specific deployed contract.
func New{{.Type}}Filterer(address *felt.Felt, filterer bind.ContractFilterer) (*{{.Type}}Filterer, error) {
	contract, err := bind{{.Type}}(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &{{.Type}}Filterer{contract: contract}, nil
}

// bind{{.Type}} binds a generic wrapper to an already deployed contract.
func bind{{.Type}}(address *felt.Felt, caller bind.ContractCaller, transactor bind.ContractTransact, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader({{.Type}}ABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

{{if .InputBin}}
// Deploy{{.Type}} deploys a new Cairo contract, binding an instance of {{.Type}} to it.
func Deploy{{.Type}}(auth *bind.TransactOpts, backend bind.ContractBackend {{range .Constructor.Inputs}}, {{.Name}} {{bindtype .Type}}{{end}}) (common.Address, *types.Transaction, *{{.Type}}, error) {
	parsed, err := abi.JSON(strings.NewReader({{.Type}}ABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	{{range $pattern, $name := .Libraries}}
	{{$pattern}}Addr, err := deployments.GetAddressByName("{{$name}}")
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	{{.Type}}Bin = strings.ReplaceAll({{.Type}}Bin, "{{$pattern}}", strings.TrimPrefix({{$pattern}}Addr.String(), "0x"))
	{{end}}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex({{.Type}}Bin), backend {{range .Constructor.Inputs}}, {{.Name}}{{end}})
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &{{.Type}}{{{.Type}}Caller: {{.Type}}Caller{contract: contract}, {{.Type}}Transactor: {{.Type}}Transactor{contract: contract}, {{.Type}}Filterer: {{.Type}}Filterer{contract: contract}}, nil
}
{{end}}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_{{$contract.Type}} *{{$contract.Type}}Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _{{$contract.Type}}.Contract.{{$contract.Type}}Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_{{$contract.Type}} *{{$contract.Type}}Raw) Transfer(opts *bind.TransactOpts) (*rpc.InvokeTxnResponse, error) {
	return _{{$contract.Type}}.Contract.{{$contract.Type}}Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_{{$contract.Type}} *{{$contract.Type}}Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*rpc.InvokeTxnResponse, error) {
	return _{{$contract.Type}}.Contract.{{$contract.Type}}Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_{{$contract.Type}} *{{$contract.Type}}CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _{{$contract.Type}}.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_{{$contract.Type}} *{{$contract.Type}}TransactorRaw) Transfer(opts *bind.TransactOpts) (*rpc.InvokeTxnResponse, error) {
	return _{{$contract.Type}}.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_{{$contract.Type}} *{{$contract.Type}}TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*rpc.InvokeTxnResponse, error) {
	return _{{$contract.Type}}.Contract.contract.Transact(opts, method, params...)
}

{{range .Methods}}
{{if .Const}}
// {{.Normalized.Name}} is a free data retrieval call binding the contract method 0x{{printf "%x" .Original.Selector}}.
//
// Solidity: {{.Original.String}}
func (_{{$contract.Type}} *{{$contract.Type}}Caller) {{.Normalized.Name}}(opts *bind.CallOpts {{range .Normalized.Inputs}}, {{.Name}} {{bindtype .Type}}{{end}}) ({{if .Structured}}struct{ {{range .Normalized.Outputs}}{{.Name}} {{bindtype .Type}};{{end}} }{{else}}{{range .Normalized.Outputs}}{{bindtype .Type}}{{end}}{{end}}, error) {
	var out []interface{}
	err := _{{$contract.Type}}.contract.Call(opts, &out, "{{.Original.Name}}" {{range .Normalized.Inputs}}, {{.Name}}{{end}})
	{{if .Structured}}
	outstruct := new(struct{ {{range .Normalized.Outputs}}{{.Name}} {{bindtype .Type}}
	{{end}} })
	if err != nil {
		return *outstruct, err
	}
	{{range $i, $t := .Normalized.Outputs}}
	outstruct.{{.Name}} = out[{{$i}}].({{bindtype .Type}}){{end}}
	return *outstruct, err
	{{else}}
	if err != nil {
		return {{range $i, $t := .Normalized.Outputs}}{{if $i}}, {{end}}{{if eq (bindtype .Type) "*big.Int"}}new(big.Int){{else if eq (bindtype .Type) "*felt.Felt"}}new(felt.Felt){{else}}{{bindtype .Type}}{}{{end}}{{end}}, err
	}
	{{range $i, $t := .Normalized.Outputs}}
	out{{$i}} := out[{{$i}}].({{bindtype .Type}}){{end}}
	return {{range $i, $t := .Normalized.Outputs}}{{if $i}}, {{end}}out{{$i}}{{end}}, err
	{{end}}
}

// {{.Normalized.Name}} is a free data retrieval call binding the contract method 0x{{printf "%x" .Original.Selector}}.
//
// Solidity: {{.Original.String}}
func (_{{$contract.Type}} *{{$contract.Type}}Session) {{.Normalized.Name}}({{range $i, $_ := .Normalized.Inputs}}{{if $i}}, {{end}}{{.Name}} {{bindtype .Type}}{{end}}) ({{if .Structured}}struct{ {{range .Normalized.Outputs}}{{.Name}} {{bindtype .Type}};{{end}} }{{else}}{{range .Normalized.Outputs}}{{bindtype .Type}}{{end}}{{end}}, error) {
	return _{{$contract.Type}}.Contract.{{.Normalized.Name}}(&_{{$contract.Type}}.CallOpts {{range .Normalized.Inputs}}, {{.Name}}{{end}})
}

// {{.Normalized.Name}} is a free data retrieval call binding the contract method 0x{{printf "%x" .Original.Selector}}.
//
// Solidity: {{.Original.String}}
func (_{{$contract.Type}} *{{$contract.Type}}CallerSession) {{.Normalized.Name}}({{range $i, $_ := .Normalized.Inputs}}{{if $i}}, {{end}}{{.Name}} {{bindtype .Type}}{{end}}) ({{if .Structured}}struct{ {{range .Normalized.Outputs}}{{.Name}} {{bindtype .Type}};{{end}} }{{else}}{{range .Normalized.Outputs}}{{bindtype .Type}}{{end}}{{end}}, error) {
	return _{{$contract.Type}}.Contract.{{.Normalized.Name}}(&_{{$contract.Type}}.CallOpts {{range .Normalized.Inputs}}, {{.Name}}{{end}})
}
{{else}}
// {{.Normalized.Name}} is a paid mutator transaction binding the contract method 0x{{printf "%x" .Original.Selector}}.
//
// Solidity: {{.Original.String}}
func (_{{$contract.Type}} *{{$contract.Type}}Transactor) {{.Normalized.Name}}(opts *bind.TransactOpts {{range .Normalized.Inputs}}, {{.Name}} {{bindtype .Type}}{{end}}) (*rpc.InvokeTxnResponse, error) {
	return _{{$contract.Type}}.contract.Transact(opts, "{{.Original.Name}}" {{range .Normalized.Inputs}}, {{.Name}}{{end}})
}

// {{.Normalized.Name}} is a paid mutator transaction binding the contract method 0x{{printf "%x" .Original.Selector}}.
//
// Solidity: {{.Original.String}}
func (_{{$contract.Type}} *{{$contract.Type}}Session) {{.Normalized.Name}}({{range $i, $_ := .Normalized.Inputs}}{{if $i}}, {{end}}{{.Name}} {{bindtype .Type}}{{end}}) (*rpc.InvokeTxnResponse, error) {
	return _{{$contract.Type}}.Contract.{{.Normalized.Name}}(&_{{$contract.Type}}.TransactOpts {{range .Normalized.Inputs}}, {{.Name}}{{end}})
}

// {{.Normalized.Name}} is a paid mutator transaction binding the contract method 0x{{printf "%x" .Original.Selector}}.
//
// Solidity: {{.Original.String}}
func (_{{$contract.Type}} *{{$contract.Type}}TransactorSession) {{.Normalized.Name}}({{range $i, $_ := .Normalized.Inputs}}{{if $i}}, {{end}}{{.Name}} {{bindtype .Type}}{{end}}) (*rpc.InvokeTxnResponse, error) {
	return _{{$contract.Type}}.Contract.{{.Normalized.Name}}(&_{{$contract.Type}}.TransactOpts {{range .Normalized.Inputs}}, {{.Name}}{{end}})
}
{{end}}
{{end}}

{{range .Events}}
// {{$contract.Type}}{{.Normalized.Name}}Iterator is returned from Filter{{.Normalized.Name}} and is used to iterate over the raw logs and unpacked data for {{.Normalized.Name}} events raised by the {{$contract.Type}} contract.
type {{$contract.Type}}{{.Normalized.Name}}Iterator struct {
	Event *{{$contract.Type}}{{.Normalized.Name}} // Event containing the contract specifics and raw log

	cursor   int                 // Current reading index of the iterator
	logs     []rpc.Event         // The logs matching the name and signature
	contract *bind.BoundContract // Subscription for logs from the contract
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *{{$contract.Type}}{{.Normalized.Name}}Iterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.cursor >= len(it.logs) {
		return false
	}
	it.Event.Raw = it.logs[it.cursor]
	it.cursor += 1
	return true
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *{{$contract.Type}}{{.Normalized.Name}}Iterator) Error() error {
	return nil
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *{{$contract.Type}}{{.Normalized.Name}}Iterator) Close() error {
	return nil
}

// {{$contract.Type}}{{.Normalized.Name}} represents a {{.Normalized.Name}} event raised by the {{$contract.Type}} contract.
type {{$contract.Type}}{{.Normalized.Name}} struct { {{range .Normalized.Keys}}
	{{.Name}} {{bindtype .Type}} {{if .Indexed}}indexed{{end}}; {{end}} {{range .Normalized.Data}}
	{{.Name}} {{bindtype .Type}}; {{end}}
	Raw rpc.Event // Blockchain specific contextual infos
}

// Filter{{.Normalized.Name}} is a free log retrieval operation binding the contract event 0x{{printf "%x" .Original.ID}}.
//
// Solidity: {{.Original.String}}
func (_{{$contract.Type}} *{{$contract.Type}}Filterer) Filter{{.Normalized.Name}}(ctx context.Context, fromBlock *felt.Felt, toBlock *felt.Felt{{range .Normalized.Keys}}{{if .Indexed}}, {{.Name}} []{{bindtype .Type}}{{end}}{{end}}) (*{{$contract.Type}}{{.Normalized.Name}}Iterator, error) {
	{{range .Normalized.Keys}}{{if .Indexed}}var {{.Name}}Rule []interface{}
	for _, {{.Name}}Item := range {{.Name}} {
		{{.Name}}Rule = append({{.Name}}Rule, {{.Name}}Item)
	}{{end}}{{end}}

	logs, err := _{{$contract.Type}}.contract.FilterEvents(ctx, rpc.EventFilter{
		FromBlock: fromBlock,
		ToBlock:   toBlock,
		Address:   []felt.Felt{*_{{$contract.Type}}.contract.Address()},
		Keys: [][]felt.Felt{
			{*utils.GetSelectorFromNameFelt("{{.Original.Name}}")},
			{{range .Normalized.Keys}}{{if .Indexed}}{*utils.GetSelectorFromNameFelt("{{.Name}}")},{{end}}{{end}}
		},
	})
	if err != nil {
		return nil, err
	}
	return &{{$contract.Type}}{{.Normalized.Name}}Iterator{
		Event: &{{$contract.Type}}{{.Normalized.Name}}{
			Raw: logs[0],
		},
		cursor:   0,
		logs:     logs,
		contract: _{{$contract.Type}}.contract,
	}, nil
}

{{end}}
{{end}}

// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.
