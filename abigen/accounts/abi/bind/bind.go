package bind

import (
	"fmt"
	"strings"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/abigen/accounts/abi"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/NethermindEth/starknet.go/utils"
)

type CallOpts struct {
	BlockID *rpc.BlockID // BlockID to call the contract at
}

type TransactOpts struct {
	From       *felt.Felt // StarkNet account address
	MaxFee     *felt.Felt // Maximum fee to pay for transaction
	Nonce      *felt.Felt // Nonce for the transaction
	Version    *felt.Felt // Transaction version
	ChainID    string     // Chain ID for the transaction
	WaitForTx  bool       // Wait for transaction to be included in a block
	WaitForTxn bool       // Wait for transaction to be included in a block (deprecated, use WaitForTx)
}

type BoundContract struct {
	address  *felt.Felt       // Deployment address of the contract on the StarkNet network
	abi      abi.ABI          // Reflect based ABI to access the correct StarkNet methods
	caller   ContractCaller   // Read interface to interact with the blockchain
	transact ContractTransact // Write interface to interact with the blockchain
	filterer ContractFilterer // Event filtering to interact with the blockchain
}

func NewBoundContract(address *felt.Felt, abi abi.ABI, caller ContractCaller, transact ContractTransact, filterer ContractFilterer) *BoundContract {
	return &BoundContract{
		address:  address,
		abi:      abi,
		caller:   caller,
		transact: transact,
		filterer: filterer,
	}
}

func (c *BoundContract) Call(opts *CallOpts, results *[]interface{}, method string, params ...interface{}) error {
	if opts == nil {
		opts = &CallOpts{}
	}

	m, exist := c.abi.Methods[method]
	if !exist {
		return fmt.Errorf("method '%s' not found in ABI", method)
	}

	calldata, err := abi.PackArguments(m.Inputs, params)
	if err != nil {
		return err
	}

	call := rpc.FunctionCall{
		ContractAddress:    c.address,
		EntryPointSelector: utils.GetSelectorFromNameFelt(method),
		Calldata:           calldata,
	}

	blockID := rpc.BlockID{Tag: "latest"}
	if opts.BlockID != nil {
		blockID = *opts.BlockID
	}

	output, err := c.caller.Call(call, blockID)
	if err != nil {
		return err
	}

	if results != nil {
		values, err := abi.UnpackValues(m.Outputs, output)
		if err != nil {
			return err
		}
		*results = values
	}

	return nil
}

func (c *BoundContract) Transact(opts *TransactOpts, method string, params ...interface{}) (*InvokeTxnResponse, error) {
	if opts == nil {
		opts = &TransactOpts{}
	}

	m, exist := c.abi.Methods[method]
	if !exist {
		return nil, fmt.Errorf("method '%s' not found in ABI", method)
	}

	calldata, err := abi.PackArguments(m.Inputs, params)
	if err != nil {
		return nil, err
	}

	call := rpc.FunctionCall{
		ContractAddress:    c.address,
		EntryPointSelector: utils.GetSelectorFromNameFelt(method),
		Calldata:           calldata,
	}

	return c.transact.Invoke(opts, call)
}

func DeployContract(opts *TransactOpts, abi abi.ABI, bytecode []byte, backend ContractBackend, params ...interface{}) (*felt.Felt, *AddTxnResponse, *BoundContract, error) {
	if opts == nil {
		opts = &TransactOpts{}
	}

	if abi.Constructor.Type != "constructor" {
		return nil, nil, nil, fmt.Errorf("constructor not found in ABI")
	}

	constructorArgs, err := abi.PackArguments(abi.Constructor.Inputs, params)
	if err != nil {
		return nil, nil, nil, err
	}

	address, txn, err := backend.DeployContract(opts, bytecode, constructorArgs)
	if err != nil {
		return nil, nil, nil, err
	}

	contract := NewBoundContract(address, abi, backend, backend, backend)

	return address, txn, contract, nil
}

func ToCairoType(goType string) string {
	switch goType {
	case "*felt.Felt":
		return "core::felt252"
	case "uint32":
		return "core::integer::u32"
	case "uint64":
		return "core::integer::u64"
	case "*big.Int":
		return "core::integer::u256"
	case "bool":
		return "core::bool"
	default:
		if strings.HasPrefix(goType, "[]") {
			elementType := goType[2:]
			return fmt.Sprintf("core::array::Array<%s>", ToCairoType(elementType))
		}
		return "core::felt252"
	}
}

func ToGoType(cairoType string) string {
	parts := strings.Split(cairoType, "::")
	baseType := parts[len(parts)-1]

	switch baseType {
	case "felt252":
		return "*felt.Felt"
	case "u8", "u16", "u32":
		return "uint32"
	case "u64":
		return "uint64"
	case "u128":
		return "uint64" // Go doesn't have uint128, use uint64 or big.Int
	case "u256":
		return "*big.Int"
	case "bool":
		return "bool"
	case "ContractAddress":
		return "*felt.Felt"
	default:
		if strings.HasPrefix(baseType, "Array<") && strings.HasSuffix(baseType, ">") {
			elementType := strings.TrimSuffix(strings.TrimPrefix(baseType, "Array<"), ">")
			return "[]" + ToGoType(elementType)
		}
		return "*felt.Felt"
	}
}
