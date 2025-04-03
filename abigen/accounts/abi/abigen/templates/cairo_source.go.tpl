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
var {{.Type}}MetaData = &bind.MetaData{
	ABI: `{{.InputABI}}`{{if .InputBin}},
	Bin: `{{.InputBin}}`{{end}},
}

const {{.Type}}ABI = `{{.InputABI}}`

{{if .InputBin}}
var {{.Type}}Bin = `{{.InputBin}}`
{{end}}

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

type {{.Type}}Session struct {
	Contract     *{{.Type}}        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction options to use throughout this session
}

type {{.Type}}CallerSession struct {
	Contract *{{.Type}}Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

type {{.Type}}TransactorSession struct {
	Contract     *{{.Type}}Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction options to use throughout this session
}

type {{.Type}}Raw struct {
	Contract *{{.Type}} // Generic contract binding to access the raw methods on
}

type {{.Type}}CallerRaw struct {
	Contract *{{.Type}}Caller // Generic read-only contract binding to access the raw methods on
}

type {{.Type}}TransactorRaw struct {
	Contract *{{.Type}}Transactor // Generic write-only contract binding to access the raw methods on
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

func New{{.Type}}Caller(address *felt.Felt, caller bind.ContractCaller) (*{{.Type}}Caller, error) {
	contract, err := bind{{.Type}}(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &{{.Type}}Caller{contract: contract}, nil
}

func New{{.Type}}Transactor(address *felt.Felt, transactor bind.ContractTransact) (*{{.Type}}Transactor, error) {
	contract, err := bind{{.Type}}(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &{{.Type}}Transactor{contract: contract}, nil
}

func New{{.Type}}Filterer(address *felt.Felt, filterer bind.ContractFilterer) (*{{.Type}}Filterer, error) {
	contract, err := bind{{.Type}}(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &{{.Type}}Filterer{contract: contract}, nil
}

func bind{{.Type}}(address *felt.Felt, caller bind.ContractCaller, transactor bind.ContractTransact, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader({{.Type}}ABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

{{if .InputBin}}
func Deploy{{.Type}}(auth *bind.TransactOpts, backend bind.ContractBackend {{range .Constructor.Inputs}}, {{.Name}} {{bindtype .Type}}{{end}}) (common.Address, *rpc.InvokeTxnResponse, *{{.Type}}, error) {
	parsed, err := abi.JSON(strings.NewReader({{.Type}}ABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex({{.Type}}Bin), backend {{range .Constructor.Inputs}}, {{.Name}}{{end}})
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &{{.Type}}{
		{{.Type}}Caller:     {{.Type}}Caller{contract: contract},
		{{.Type}}Transactor: {{.Type}}Transactor{contract: contract},
		{{.Type}}Filterer:   {{.Type}}Filterer{contract: contract},
	}, nil
}
{{end}}

func (_{{.Type}} *{{.Type}}Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _{{.Type}}.Contract.{{.Type}}Caller.contract.Call(opts, result, method, params...)
}

func (_{{.Type}} *{{.Type}}Raw) Transfer(opts *bind.TransactOpts) (*rpc.InvokeTxnResponse, error) {
	return _{{.Type}}.Contract.{{.Type}}Transactor.contract.Transfer(opts)
}

func (_{{.Type}} *{{.Type}}Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*rpc.InvokeTxnResponse, error) {
	return _{{.Type}}.Contract.{{.Type}}Transactor.contract.Transact(opts, method, params...)
}

func (_{{.Type}} *{{.Type}}CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _{{.Type}}.Contract.contract.Call(opts, result, method, params...)
}

func (_{{.Type}} *{{.Type}}TransactorRaw) Transfer(opts *bind.TransactOpts) (*rpc.InvokeTxnResponse, error) {
	return _{{.Type}}.Contract.contract.Transfer(opts)
}

func (_{{.Type}} *{{.Type}}TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*rpc.InvokeTxnResponse, error) {
	return _{{.Type}}.Contract.contract.Transact(opts, method, params...)
}

{{range $method := .Methods}}
{{if .Const}}
//
func (_{{$contract.Type}} *{{$contract.Type}}Caller) {{.Normalized.Name}}(opts *bind.CallOpts {{range .Normalized.Inputs}}, {{.Name}} {{bindtype .Type}}{{end}}) ({{if .Structured}}struct{ {{range .Normalized.Outputs}}{{.Name}} {{bindtype .Type}};{{end}} }{{else}}{{range .Normalized.Outputs}}{{bindtype .Type}}{{end}}{{end}}, error) {
	var out []interface{}
	err := _{{$contract.Type}}.contract.Call(opts, &out, "{{.Original.Name}}" {{range .Normalized.Inputs}}, {{.Name}}{{end}})
	
	if err != nil {
		return {{if .Structured}}*new(struct{ {{range .Normalized.Outputs}}{{.Name}} {{bindtype .Type}};{{end}} }){{else}}{{range .Normalized.Outputs}}*new({{bindtype .Type}}){{end}}{{end}}, err
	}
	
	{{if .Structured}}
	out0 := *abi.ConvertType(out[0], new(struct{ {{range .Normalized.Outputs}}{{.Name}} {{bindtype .Type}};{{end}} })).(*struct{ {{range .Normalized.Outputs}}{{.Name}} {{bindtype .Type}};{{end}} })
	return out0, err
	{{else}}
	{{range $i, $_ := .Normalized.Outputs}}
	out{{$i}} := {{if eq (bindtype .Type) "*big.Int"}}*{{end}}abi.ConvertType(out[{{$i}}], new({{bindtype .Type}})).(*{{bindtype .Type}}){{end}}
	return {{range $i, $_ := .Normalized.Outputs}}{{if $i}}, {{end}}out{{$i}}{{end}}, err
	{{end}}
}

//
func (_{{$contract.Type}} *{{$contract.Type}}Session) {{.Normalized.Name}}({{range $i, $_ := .Normalized.Inputs}}{{if $i}}, {{end}}{{.Name}} {{bindtype .Type}}{{end}}) ({{if .Structured}}struct{ {{range .Normalized.Outputs}}{{.Name}} {{bindtype .Type}};{{end}} }{{else}}{{range .Normalized.Outputs}}{{bindtype .Type}}{{end}}{{end}}, error) {
	return _{{$contract.Type}}.Contract.{{.Normalized.Name}}(&_{{$contract.Type}}.CallOpts {{range .Normalized.Inputs}}, {{.Name}}{{end}})
}

//
func (_{{$contract.Type}} *{{$contract.Type}}CallerSession) {{.Normalized.Name}}({{range $i, $_ := .Normalized.Inputs}}{{if $i}}, {{end}}{{.Name}} {{bindtype .Type}}{{end}}) ({{if .Structured}}struct{ {{range .Normalized.Outputs}}{{.Name}} {{bindtype .Type}};{{end}} }{{else}}{{range .Normalized.Outputs}}{{bindtype .Type}}{{end}}{{end}}, error) {
	return _{{$contract.Type}}.Contract.{{.Normalized.Name}}(&_{{$contract.Type}}.CallOpts {{range .Normalized.Inputs}}, {{.Name}}{{end}})
}
{{else}}
//
func (_{{$contract.Type}} *{{$contract.Type}}Transactor) {{.Normalized.Name}}(opts *bind.TransactOpts {{range .Normalized.Inputs}}, {{.Name}} {{bindtype .Type}}{{end}}) (*rpc.InvokeTxnResponse, error) {
	return _{{$contract.Type}}.contract.Transact(opts, "{{.Original.Name}}" {{range .Normalized.Inputs}}, {{.Name}}{{end}})
}

//
func (_{{$contract.Type}} *{{$contract.Type}}Session) {{.Normalized.Name}}({{range $i, $_ := .Normalized.Inputs}}{{if $i}}, {{end}}{{.Name}} {{bindtype .Type}}{{end}}) (*rpc.InvokeTxnResponse, error) {
	return _{{$contract.Type}}.Contract.{{.Normalized.Name}}(&_{{$contract.Type}}.TransactOpts {{range .Normalized.Inputs}}, {{.Name}}{{end}})
}

//
func (_{{$contract.Type}} *{{$contract.Type}}TransactorSession) {{.Normalized.Name}}({{range $i, $_ := .Normalized.Inputs}}{{if $i}}, {{end}}{{.Name}} {{bindtype .Type}}{{end}}) (*rpc.InvokeTxnResponse, error) {
	return _{{$contract.Type}}.Contract.{{.Normalized.Name}}(&_{{$contract.Type}}.TransactOpts {{range .Normalized.Inputs}}, {{.Name}}{{end}})
}
{{end}}
{{end}}

{{range $event := .Events}}
type {{$contract.Type}}{{.Normalized.Name}}Iterator struct {
	Event *{{$contract.Type}}{{.Normalized.Name}} // Event containing the contract specifics and raw log

	cursor   int                 // Current reading index of the iterator
	logs     []rpc.Event         // The logs matching the name and signature
	contract *bind.BoundContract // Subscription for logs from the contract
}

func (it *{{$contract.Type}}{{.Normalized.Name}}Iterator) Next() bool {
	if it.cursor >= len(it.logs) {
		return false
	}
	it.Event.Raw = it.logs[it.cursor]
	it.cursor += 1
	return true
}

func (it *{{$contract.Type}}{{.Normalized.Name}}Iterator) Error() error {
	return nil
}

func (it *{{$contract.Type}}{{.Normalized.Name}}Iterator) Close() error {
	return nil
}

type {{$contract.Type}}{{.Normalized.Name}} struct {
	{{range $i, $key := .Normalized.Keys}}{{$key.Name}} {{bindtype $key.Type}}
	{{end}}{{range $i, $data := .Normalized.Data}}{{$data.Name}} {{bindtype $data.Type}}
	{{end}}
	Raw rpc.Event // Blockchain specific contextual infos
}

//
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

`
