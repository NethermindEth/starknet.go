package account

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/contracts"
	"github.com/NethermindEth/starknet.go/hash"
	"github.com/NethermindEth/starknet.go/rpc"
)

var (
	ErrTxnTypeUnSupported    = errors.New("unsupported transaction type")
	ErrTxnVersionUnSupported = errors.New("unsupported transaction version")
)

//go:generate mockgen -destination=../mocks/mock_account.go -package=mocks -source=account.go AccountInterface
type AccountInterface interface {
	BuildAndEstimateDeployAccountTxn(
		ctx context.Context,
		salt *felt.Felt,
		classHash *felt.Felt,
		constructorCalldata []*felt.Felt,
		opts *TxnOptions,
	) (*rpc.BroadcastDeployAccountTxnV3, *felt.Felt, error)
	BuildAndSendInvokeTxn(
		ctx context.Context,
		functionCalls []rpc.InvokeFunctionCall,
		opts *TxnOptions,
	) (*rpc.AddInvokeTransactionResponse, error)
	BuildAndSendDeclareTxn(
		ctx context.Context,
		casmClass *contracts.CasmClass,
		contractClass *contracts.ContractClass,
		opts *TxnOptions,
	) (*rpc.AddDeclareTransactionResponse, error)
	Nonce(ctx context.Context) (*felt.Felt, error)
	SendTransaction(ctx context.Context, txn rpc.BroadcastTxn) (*rpc.TransactionResponse, error)
	Sign(ctx context.Context, msg *felt.Felt) ([]*felt.Felt, error)
	SignInvokeTransaction(ctx context.Context, tx rpc.InvokeTxnType) error
	SignDeployAccountTransaction(ctx context.Context, tx rpc.DeployAccountType, precomputeAddress *felt.Felt) error
	SignDeclareTransaction(ctx context.Context, tx rpc.DeclareTxnType) error
	TransactionHashInvoke(invokeTxn rpc.InvokeTxnType) (*felt.Felt, error)
	TransactionHashDeployAccount(tx rpc.DeployAccountType, contractAddress *felt.Felt) (*felt.Felt, error)
	TransactionHashDeclare(tx rpc.DeclareTxnType) (*felt.Felt, error)
	WaitForTransactionReceipt(
		ctx context.Context,
		transactionHash *felt.Felt,
		pollInterval time.Duration,
	) (*rpc.TransactionReceiptWithBlockInfo, error)
}

var _ AccountInterface = &Account{} //nolint:exhaustruct

const BRAAVOS_WARNING_MESSAGE = `WARNING: Currently, Braavos accounts are incompatible with transactions sent via
RPC 0.8.0. Ref: https://community.starknet.io/t/starknet-devtools-for-0-13-5/115495#p-2359168-braavos-compatibility-issues-3`

type Account struct {
	Provider     rpc.RpcProvider
	ChainId      *felt.Felt
	Address      *felt.Felt
	publicKey    string
	CairoVersion int
	ks           Keystore
}

// NewAccount creates a new Account instance.
//
// Parameters:
//   - provider: the provider to use
//   - accountAddress: the account address
//   - publicKey: the public key of the account
//   - keystore: the keystore to use
//   - cairoVersion: the cairo version of the account (0 or 2)
//
// It returns:
//   - *Account: a pointer to newly created Account
//   - error: an error if any
func NewAccount(
	provider rpc.RpcProvider,
	accountAddress *felt.Felt,
	publicKey string,
	keystore Keystore,
	cairoVersion int,
) (*Account, error) {
	// TODO: Remove this temporary check once solved (starknet v0.14.0 should do it)
	// This temporary check is to warn the user that Braavos account restricts transactions to have exactly two resource fields.
	// This makes them incompatible with transactions sent via RPC 0.8.0
	accClassHash, err := provider.ClassHashAt(context.Background(), rpc.WithBlockTag("latest"), accountAddress)
	// ignoring the error to not break mock tests (if the provider is not working, it will return an error in the next ChainID call anyway)
	if err == nil {
		// Since felt.Felt.String() returns a string without leading zeros, we need to remove them from the
		// class hashes for the comparison
		braavosClassHashes := []string{
			// Original class hash: 0x02c8c7e6fbcfb3e8e15a46648e8914c6aa1fc506fc1e7fb3d1e19630716174bc
			"0x2c8c7e6fbcfb3e8e15a46648e8914c6aa1fc506fc1e7fb3d1e19630716174bc",
			// Original class hash: 0x00816dd0297efc55dc1e7559020a3a825e81ef734b558f03c83325d4da7e6253
			"0x816dd0297efc55dc1e7559020a3a825e81ef734b558f03c83325d4da7e6253",
			// Original class hash: 0x041bf1e71792aecb9df3e9d04e1540091c5e13122a731e02bec588f71dc1a5c3
			"0x41bf1e71792aecb9df3e9d04e1540091c5e13122a731e02bec588f71dc1a5c3",
		}
		if slices.Contains(braavosClassHashes, accClassHash.String()) {
			fmt.Print(BRAAVOS_WARNING_MESSAGE + "\n\n")
		}
	}

	chainID, err := provider.ChainID(context.Background())
	if err != nil {
		return nil, err
	}

	account := &Account{
		Provider:     provider,
		Address:      accountAddress,
		publicKey:    publicKey,
		ks:           keystore,
		CairoVersion: cairoVersion,
		ChainId:      new(felt.Felt).SetBytes([]byte(chainID)),
	}

	return account, nil
}

// Nonce retrieves the nonce for the account's contract address.
func (account *Account) Nonce(ctx context.Context) (*felt.Felt, error) {
	return account.Provider.Nonce(context.Background(), rpc.WithBlockTag("pending"), account.Address)
}

// TransactionHashDeployAccount calculates the transaction hash for a deploy account transaction.
//
// Parameters:
//   - tx: The deploy account transaction to calculate the hash for. Can be of type DeployAccountTxn or DeployAccountTxnV3.
//   - contractAddress: The contract address as parameters as a *felt.Felt
//
// Returns:
//   - *felt.Felt: the calculated transaction hash
//   - error: an error if any
func (account *Account) TransactionHashDeployAccount(tx rpc.DeployAccountType, contractAddress *felt.Felt) (*felt.Felt, error) {
	// https://docs.starknet.io/architecture-and-concepts/network-architecture/transactions/#deploy_account_transaction
	switch txn := tx.(type) {
	// deployAccTxn v1, pointer and struct
	case *rpc.DeployAccountTxnV1:
		return hash.TransactionHashDeployAccountV1(txn, contractAddress, account.ChainId)
	case rpc.DeployAccountTxnV1:
		return hash.TransactionHashDeployAccountV1(&txn, contractAddress, account.ChainId)
	// deployAccTxn v3, pointer and struct
	case *rpc.DeployAccountTxnV3:
		return hash.TransactionHashDeployAccountV3(txn, contractAddress, account.ChainId)
	case rpc.DeployAccountTxnV3:
		return hash.TransactionHashDeployAccountV3(&txn, contractAddress, account.ChainId)
	default:
		return nil, fmt.Errorf("%w: got '%T' instead of a valid invoke txn type", ErrTxnTypeUnSupported, txn)
	}
}

// TransactionHashInvoke calculates the transaction hash for the given invoke transaction.
//
// Parameters:
//   - tx: The invoke transaction to calculate the hash for. Can be of type InvokeTxnV0, InvokeTxnV1, or InvokeTxnV3.
//
// Returns:
//   - *felt.Felt: The calculated transaction hash as a *felt.Felt
//   - error: an error, if any
//
// If the transaction type is unsupported, the function returns an error.
func (account *Account) TransactionHashInvoke(tx rpc.InvokeTxnType) (*felt.Felt, error) {
	switch txn := tx.(type) {
	// invoke v0, pointer and struct
	case *rpc.InvokeTxnV0:
		return hash.TransactionHashInvokeV0(txn, account.ChainId)
	case rpc.InvokeTxnV0:
		return hash.TransactionHashInvokeV0(&txn, account.ChainId)
	// invoke v1, pointer and struct
	case *rpc.InvokeTxnV1:
		return hash.TransactionHashInvokeV1(txn, account.ChainId)
	case rpc.InvokeTxnV1:
		return hash.TransactionHashInvokeV1(&txn, account.ChainId)
	// invoke v3, pointer and struct
	case *rpc.InvokeTxnV3:
		return hash.TransactionHashInvokeV3(txn, account.ChainId)
	case rpc.InvokeTxnV3:
		return hash.TransactionHashInvokeV3(&txn, account.ChainId)
	default:
		return nil, fmt.Errorf("%w: got '%T' instead of a valid invoke txn type", ErrTxnTypeUnSupported, txn)
	}
}

// TransactionHashDeclare calculates the transaction hash for declaring a transaction type.
//
// Parameters:
//   - tx: The `tx` parameter of type `rpc.DeclareTxnType`. Can be one of the types DeclareTxnV1/V2/V3, and BroadcastDeclareTxnV3
//
// Returns:
//   - *felt.Felt: the calculated transaction hash as `*felt.Felt` value
//   - error: an error, if any
//
// If the `tx` parameter is not one of the supported types, the function returns an error `ErrTxnTypeUnSupported`.
func (account *Account) TransactionHashDeclare(tx rpc.DeclareTxnType) (*felt.Felt, error) {
	switch txn := tx.(type) {
	// Due to inconsistencies in version 0 hash calculation we don't calculate the hash
	case *rpc.DeclareTxnV0, rpc.DeclareTxnV0:
		return nil, ErrTxnVersionUnSupported
	// declare v1, pointer and struct
	case *rpc.DeclareTxnV1:
		return hash.TransactionHashDeclareV1(txn, account.ChainId)
	case rpc.DeclareTxnV1:
		return hash.TransactionHashDeclareV1(&txn, account.ChainId)
	// declare v2, pointer and struct
	case *rpc.DeclareTxnV2:
		return hash.TransactionHashDeclareV2(txn, account.ChainId)
	case rpc.DeclareTxnV2:
		return hash.TransactionHashDeclareV2(&txn, account.ChainId)
	// declare v3, pointer and struct
	case *rpc.DeclareTxnV3:
		return hash.TransactionHashDeclareV3(txn, account.ChainId)
	case rpc.DeclareTxnV3:
		return hash.TransactionHashDeclareV3(&txn, account.ChainId)
	// broadcast declare v3, pointer and struct
	case *rpc.BroadcastDeclareTxnV3:
		return hash.TransactionHashBroadcastDeclareV3(txn, account.ChainId)
	case rpc.BroadcastDeclareTxnV3:
		return hash.TransactionHashBroadcastDeclareV3(&txn, account.ChainId)
	default:
		return nil, fmt.Errorf("%w: got '%T' instead of a valid declare txn type", ErrTxnTypeUnSupported, txn)
	}
}

// PrecomputeAccountAddress calculates the precomputed address for an account.
// ref: https://docs.starknet.io/architecture-and-concepts/smart-contracts/contract-address/
//
// Parameters:
//   - salt: the salt for the address of the deployed contract
//   - classHash: the class hash of the contract to be deployed
//   - constructorCalldata: the parameters passed to the constructor
//
// Returns:
//   - *felt.Felt: the precomputed address as a *felt.Felt
//   - error: an error if any
func PrecomputeAccountAddress(salt, classHash *felt.Felt, constructorCalldata []*felt.Felt) *felt.Felt {
	return contracts.PrecomputeAddress(&felt.Zero, salt, classHash, constructorCalldata)
}

// FmtCalldata generates the formatted calldata for the given function calls and Cairo version.
//
// Parameters:
//   - fnCalls: a slice of rpc.FunctionCall representing the function calls.
//
// Returns:
//   - a slice of *felt.Felt representing the formatted calldata.
//   - an error if Cairo version is not supported.
func (account *Account) FmtCalldata(fnCalls []rpc.FunctionCall) ([]*felt.Felt, error) {
	switch account.CairoVersion {
	case 0:
		return FmtCallDataCairo0(fnCalls), nil
	case 2:
		return FmtCallDataCairo2(fnCalls), nil
	default:
		return nil, fmt.Errorf("account cairo version '%d' not supported", account.CairoVersion)
	}
}

// FmtCallDataCairo0 generates a slice of *felt.Felt that represents the calldata for the given function calls in Cairo 0 format.
//
// Parameters:
//   - fnCalls: a slice of rpc.FunctionCall containing the function calls.
//
// Returns:
//   - a slice of *felt.Felt representing the generated calldata.
//
// https://github.com/project3fusion/StarkSharp/blob/main/StarkSharp/StarkSharp.Rpc/Modules/Transactions/Hash/TransactionHash.cs#L27
func FmtCallDataCairo0(callArray []rpc.FunctionCall) []*felt.Felt {
	calldata := make([]*felt.Felt, 0, 10) //nolint:mnd
	calls := make([]*felt.Felt, 0, 10)    //nolint:mnd

	calldata = append(calldata, new(felt.Felt).SetUint64(uint64(len(callArray))))

	offset := uint64(0)
	for _, call := range callArray {
		calldata = append(calldata, call.ContractAddress, call.EntryPointSelector, new(felt.Felt).SetUint64(offset))
		callDataLen := uint64(len(call.Calldata))
		calldata = append(calldata, new(felt.Felt).SetUint64(callDataLen))
		offset += callDataLen

		calls = append(calls, call.Calldata...)
	}

	calldata = append(calldata, new(felt.Felt).SetUint64(offset))
	calldata = append(calldata, calls...)

	return calldata
}

// FmtCallDataCairo2 generates the calldata for the given function calls for Cairo 2 contracts.
//
// Parameters:
//   - fnCalls: a slice of rpc.FunctionCall containing the function calls.
//
// Returns:
//   - a slice of *felt.Felt representing the generated calldata.
//
// https://github.com/project3fusion/StarkSharp/blob/main/StarkSharp/StarkSharp.Rpc/Modules/Transactions/Hash/TransactionHash.cs#L22
func FmtCallDataCairo2(callArray []rpc.FunctionCall) []*felt.Felt {
	result := make([]*felt.Felt, 0, 10) //nolint:mnd

	result = append(result, new(felt.Felt).SetUint64(uint64(len(callArray))))

	for _, call := range callArray {
		result = append(result, call.ContractAddress, call.EntryPointSelector)

		callDataLen := uint64(len(call.Calldata))
		result = append(result, new(felt.Felt).SetUint64(callDataLen))

		result = append(result, call.Calldata...)
	}

	return result
}

func makeResourceBoundsMapWithZeroValues() *rpc.ResourceBoundsMapping {
	return &rpc.ResourceBoundsMapping{
		L1Gas: rpc.ResourceBounds{
			MaxAmount:       "0x0",
			MaxPricePerUnit: "0x0",
		},
		L1DataGas: rpc.ResourceBounds{
			MaxAmount:       "0x0",
			MaxPricePerUnit: "0x0",
		},
		L2Gas: rpc.ResourceBounds{
			MaxAmount:       "0x0",
			MaxPricePerUnit: "0x0",
		},
	}
}
