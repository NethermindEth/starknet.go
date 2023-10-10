package account

import (
	"context"
	"errors"
	"time"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/curve"
	"github.com/NethermindEth/starknet.go/hash"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/NethermindEth/starknet.go/utils"
)

var (
	ErrNotAllParametersSet   = errors.New("Not all neccessary parameters have been set")
	ErrTxnTypeUnSupported    = errors.New("Unsupported transction type")
	ErrTxnVersionUnSupported = errors.New("Unsupported transction version")
	ErrFeltToBigInt          = errors.New("Felt to BigInt error")
)

var (
	PREFIX_TRANSACTION      = new(felt.Felt).SetBytes([]byte("invoke"))
	PREFIX_DECLARE          = new(felt.Felt).SetBytes([]byte("declare"))
	PREFIX_CONTRACT_ADDRESS = new(felt.Felt).SetBytes([]byte("STARKNET_CONTRACT_ADDRESS"))
	PREFIX_DEPLOY_ACCOUNT   = new(felt.Felt).SetBytes([]byte("deploy_account"))
)

//go:generate mockgen -destination=../mocks/mock_account.go -package=mocks -source=account.go AccountInterface
type AccountInterface interface {
	Sign(ctx context.Context, msg *felt.Felt) ([]*felt.Felt, error)
	TransactionHashInvoke(invokeTxn rpc.InvokeTxnType) (*felt.Felt, error)
	TransactionHashDeployAccount(tx rpc.DeployAccountTxn, contractAddress *felt.Felt) (*felt.Felt, error)
	TransactionHashDeclare(tx rpc.DeclareTxnType) (*felt.Felt, error)
	SignInvokeTransaction(ctx context.Context, tx *rpc.InvokeTxnV1) error
	SignDeployAccountTransaction(ctx context.Context, tx *rpc.DeployAccountTxn, precomputeAddress *felt.Felt) error
	SignDeclareTransaction(ctx context.Context, tx *rpc.DeclareTxnV2) error
	PrecomputeAddress(deployerAddress *felt.Felt, salt *felt.Felt, classHash *felt.Felt, constructorCalldata []*felt.Felt) (*felt.Felt, error)
	WaitForTransactionReceipt(ctx context.Context, transactionHash *felt.Felt, pollInterval time.Duration) (*rpc.TransactionReceipt, error)
}

var _ AccountInterface = &Account{}
var _ rpc.RpcProvider = &Account{}

type Account struct {
	provider       rpc.RpcProvider
	ChainId        *felt.Felt
	AccountAddress *felt.Felt
	publicKey      string
	ks             Keystore
}

// NewAccount creates a new Account instance.
//
// It takes in a provider of type rpc.RpåProvider, an accountAddress of type *felt.Felt,
// a publicKey of type string, and a keystore of type Keystore.
// It returns an Account pointer and an error.
func NewAccount(provider rpc.RpcProvider, accountAddress *felt.Felt, publicKey string, keystore Keystore) (*Account, error) {
	account := &Account{
		provider:       provider,
		AccountAddress: accountAddress,
		publicKey:      publicKey,
		ks:             keystore,
	}

	chainID, err := provider.ChainID(context.Background())
	if err != nil {
		return nil, err
	}
	account.ChainId = new(felt.Felt).SetBytes([]byte(chainID))

	return account, nil
}

// Sign signs the given felt message using the account's private key.
//
// ctx is the context used for the signing operation.
// msg is the felt message to be signed.
// Returns an array of signed felt messages and an error, if any.
func (account *Account) Sign(ctx context.Context, msg *felt.Felt) ([]*felt.Felt, error) {

	msgBig := utils.FeltToBigInt(msg)

	s1, s2, err := account.ks.Sign(ctx, account.publicKey, msgBig)
	if err != nil {
		return nil, err
	}
	s1Felt := utils.BigIntToFelt(s1)
	s2Felt := utils.BigIntToFelt(s2)

	return []*felt.Felt{s1Felt, s2Felt}, nil
}

// SignInvokeTransaction signs and invokes a transaction.
//
// ctx - the context.Context for the function execution.
// invokeTx - the InvokeTxnV1 struct representing the transaction to be invoked.
// Returns an error if there was an error in the signing or invoking process.
func (account *Account) SignInvokeTransaction(ctx context.Context, invokeTx *rpc.InvokeTxnV1) error {

	txHash, err := account.TransactionHashInvoke(*invokeTx)
	if err != nil {
		return err
	}
	signature, err := account.Sign(ctx, txHash)
	if err != nil {
		return err
	}
	invokeTx.Signature = signature
	return nil
}

// SignDeployAccountTransaction signs a deploy account transaction.
//
// It takes in the context, a pointer to the deploy account transaction, and a precomputed address.
// It returns an error.
func (account *Account) SignDeployAccountTransaction(ctx context.Context, tx *rpc.DeployAccountTxn, precomputeAddress *felt.Felt) error {

	hash, err := account.TransactionHashDeployAccount(*tx, precomputeAddress)
	if err != nil {
		return err
	}
	signature, err := account.Sign(ctx, hash)
	if err != nil {
		return err
	}
	tx.Signature = signature
	return nil
}

// SignDeclareTransaction signs a DeclareTxnV2 transaction using the provided Account.
//
// It takes a context.Context and a *rpc.DeclareTxnV2 as parameters.
// It returns an error.
func (account *Account) SignDeclareTransaction(ctx context.Context, tx *rpc.DeclareTxnV2) error {

	hash, err := account.TransactionHashDeclare(*tx)
	if err != nil {
		return err
	}
	signature, err := account.Sign(ctx, hash)
	if err != nil {
		return err
	}
	tx.Signature = signature
	return nil
}

// TransactionHashDeployAccount calculates the transaction hash for a deploy account transaction.
//
// It takes a DeployAccountTxn and a contract address as parameters.
// It returns the calculated transaction hash and an error.
func (account *Account) TransactionHashDeployAccount(tx rpc.DeployAccountTxn, contractAddress *felt.Felt) (*felt.Felt, error) {

	// https://docs.starknet.io/documentation/architecture_and_concepts/Network_Architecture/transactions/#deploy_account_transaction

	// There is only version 1 of deployAccount txn
	if tx.Version != rpc.TransactionV1 {
		return nil, ErrTxnTypeUnSupported
	}
	calldata := []*felt.Felt{tx.ClassHash, tx.ContractAddressSalt}
	calldata = append(calldata, tx.ConstructorCalldata...)
	calldataHash, err := hash.ComputeHashOnElementsFelt(calldata)
	if err != nil {
		return nil, err
	}

	versionFelt, err := new(felt.Felt).SetString(string(tx.Version))
	if err != nil {
		return nil, err
	}

	// https://docs.starknet.io/documentation/architecture_and_concepts/Network_Architecture/transactions/#deploy_account_hash_calculation
	return hash.CalculateTransactionHashCommon(
		PREFIX_DEPLOY_ACCOUNT,
		versionFelt,
		contractAddress,
		&felt.Zero,
		calldataHash,
		tx.MaxFee,
		account.ChainId,
		[]*felt.Felt{tx.Nonce},
	)
}

// TransactionHashInvoke calculates the transaction hash for the given invoke transaction.
//
// The function takes an invoke transaction as input and returns the calculated transaction hash and an error, if any.
// The invoke transaction can be of type InvokeTxnV0 or InvokeTxnV1.
// For InvokeTxnV0, the function checks if all the required parameters are set and then computes the transaction hash using the provided data.
// For InvokeTxnV1, the function performs similar checks and computes the transaction hash using the provided data.
// If the transaction type is unsupported, the function returns an error.
func (account *Account) TransactionHashInvoke(tx rpc.InvokeTxnType) (*felt.Felt, error) {

	// https://docs.starknet.io/documentation/architecture_and_concepts/Network_Architecture/transactions/#v0_hash_calculation
	switch txn := tx.(type) {
	case rpc.InvokeTxnV0:
		if txn.Version == "" || len(txn.Calldata) == 0 || txn.MaxFee == nil || txn.EntryPointSelector == nil {
			return nil, ErrNotAllParametersSet
		}

		calldataHash, err := hash.ComputeHashOnElementsFelt(txn.Calldata)
		if err != nil {
			return nil, err
		}

		txnVersionFelt, err := new(felt.Felt).SetString(string(txn.Version))
		if err != nil {
			return nil, err
		}
		return hash.CalculateTransactionHashCommon(
			PREFIX_TRANSACTION,
			txnVersionFelt,
			txn.ContractAddress,
			txn.EntryPointSelector,
			calldataHash,
			txn.MaxFee,
			account.ChainId,
			[]*felt.Felt{},
		)

	case rpc.InvokeTxnV1:
		if txn.Version == "" || len(txn.Calldata) == 0 || txn.Nonce == nil || txn.MaxFee == nil || txn.SenderAddress == nil {
			return nil, ErrNotAllParametersSet
		}

		calldataHash, err := hash.ComputeHashOnElementsFelt(txn.Calldata)
		if err != nil {
			return nil, err
		}
		txnVersionFelt, err := new(felt.Felt).SetString(string(txn.Version))
		if err != nil {
			return nil, err
		}
		return hash.CalculateTransactionHashCommon(
			PREFIX_TRANSACTION,
			txnVersionFelt,
			txn.SenderAddress,
			&felt.Zero,
			calldataHash,
			txn.MaxFee,
			account.ChainId,
			[]*felt.Felt{txn.Nonce},
		)
	}
	return nil, ErrTxnTypeUnSupported
}

// TransactionHashDeclare calculates the transaction hash for declaring a transaction type.
//
// It takes a `tx` parameter of type `rpc.DeclareTxnType` and returns a `*felt.Felt` value and an `error`.
// The `rpc.DeclareTxnType` can be one of the following types:
//   - `rpc.DeclareTxnV0`
//   - `rpc.DeclareTxnV1`
//   - `rpc.DeclareTxnV2`
//
// For `rpc.DeclareTxnV0`, the function returns an error `ErrTxnVersionUnSupported`.
//
// For `rpc.DeclareTxnV1` and `rpc.DeclareTxnV2`, the function performs the following steps:
//   - It checks if all the required parameters are set. If any of the required parameters are missing,
//     it returns an error `ErrNotAllParametersSet`.
//   - It calculates the hash of the `txn.ClassHash` using the `hash.ComputeHashOnElementsFelt` function.
//   - It converts the `txn.Version` string to a `*felt.Felt` value using the `new(felt.Felt).SetString` function.
//   - It calls the `hash.CalculateTransactionHashCommon` function with the following parameters:
//       - `PREFIX_DECLARE`
//       - `txnVersionFelt`
//       - `txn.SenderAddress`
//       - `&felt.Zero`
//       - `calldataHash`
//       - `txn.MaxFee`
//       - `account.ChainId`
//       - A slice containing `txn.Nonce`
//       - For `rpc.DeclareTxnV2`, it also includes `txn.CompiledClassHash` in the slice of parameters.
//   - It returns the calculated transaction hash and nil error.
//
// If the `tx` parameter is not one of the supported types, the function returns an error `ErrTxnTypeUnSupported`.
func (account *Account) TransactionHashDeclare(tx rpc.DeclareTxnType) (*felt.Felt, error) {

	switch txn := tx.(type) {
	case rpc.DeclareTxnV0:
		// Due to inconsistencies in version 0 hash calculation we don't calculate the hash
		return nil, ErrTxnVersionUnSupported
	case rpc.DeclareTxnV1:
		if txn.SenderAddress == nil || txn.Version == "" || txn.ClassHash == nil || txn.MaxFee == nil || txn.Nonce == nil {
			return nil, ErrNotAllParametersSet
		}

		calldataHash, err := hash.ComputeHashOnElementsFelt([]*felt.Felt{txn.ClassHash})
		if err != nil {
			return nil, err
		}

		txnVersionFelt, err := new(felt.Felt).SetString(string(txn.Version))
		if err != nil {
			return nil, err
		}
		return hash.CalculateTransactionHashCommon(
			PREFIX_DECLARE,
			txnVersionFelt,
			txn.SenderAddress,
			&felt.Zero,
			calldataHash,
			txn.MaxFee,
			account.ChainId,
			[]*felt.Felt{txn.Nonce},
		)
	case rpc.DeclareTxnV2:
		if txn.CompiledClassHash == nil || txn.SenderAddress == nil || txn.Version == "" || txn.ClassHash == nil || txn.MaxFee == nil || txn.Nonce == nil {
			return nil, ErrNotAllParametersSet
		}

		calldataHash, err := hash.ComputeHashOnElementsFelt([]*felt.Felt{txn.ClassHash})
		if err != nil {
			return nil, err
		}

		txnVersionFelt, err := new(felt.Felt).SetString(string(txn.Version))
		if err != nil {
			return nil, err
		}
		return hash.CalculateTransactionHashCommon(
			PREFIX_DECLARE,
			txnVersionFelt,
			txn.SenderAddress,
			&felt.Zero,
			calldataHash,
			txn.MaxFee,
			account.ChainId,
			[]*felt.Felt{txn.Nonce, txn.CompiledClassHash},
		)
	}

	return nil, ErrTxnTypeUnSupported
}

// PrecomputeAddress calculates the precomputed address for an account.
// ref: https://github.com/starkware-libs/cairo-lang/blob/master/src/starkware/starknet/core/os/contract_address/contract_address.py
//
// It takes the deployer address, salt, class hash, and constructor calldata
// as parameters.
// It returns the precomputed address as a *felt.Felt and an error if any.
func (account *Account) PrecomputeAddress(deployerAddress *felt.Felt, salt *felt.Felt, classHash *felt.Felt, constructorCalldata []*felt.Felt) (*felt.Felt, error) {

	bigIntArr := utils.FeltArrToBigIntArr([]*felt.Felt{
		PREFIX_CONTRACT_ADDRESS,
		deployerAddress,
		salt,
		classHash,
	})

	constructorCalldataBigIntArr := utils.FeltArrToBigIntArr(constructorCalldata)
	constructorCallDataHashInt, _ := curve.Curve.ComputeHashOnElements(constructorCalldataBigIntArr)
	bigIntArr = append(bigIntArr, constructorCallDataHashInt)

	preBigInt, err := curve.Curve.ComputeHashOnElements(bigIntArr)
	if err != nil {
		return nil, err
	}
	return utils.BigIntToFelt(preBigInt), nil

}

// WaitForTransactionReceipt waits for the transaction receipt of the given transaction hash to succeed or fail.
//
// It takes a context, a transaction hash, and a poll interval as parameters.
// It returns the transaction receipt and an error.
func (account *Account) WaitForTransactionReceipt(ctx context.Context, transactionHash *felt.Felt, pollInterval time.Duration) (*rpc.TransactionReceipt, error) {
	t := time.NewTicker(pollInterval)
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-t.C:
			receipt, err := account.TransactionReceipt(ctx, transactionHash)
			if err != nil {
				if err.Error() == rpc.ErrHashNotFound.Error() {
					continue
				} else {
					return nil, err
				}
			}
			return &receipt, nil
		}
	}
}

// AddInvokeTransaction generates an invoke transaction and adds it to the account's provider.
//
// ctx - the context.Context object for the transaction.
// invokeTx - the invoke transaction to be added.
// Returns the AddInvokeTransactionResponse and an error if any.
func (account *Account) AddInvokeTransaction(ctx context.Context, invokeTx rpc.InvokeTxnV1) (*rpc.AddInvokeTransactionResponse, error) {
	return account.provider.AddInvokeTransaction(ctx, invokeTx)
}

// AddDeclareTransaction adds a declare transaction to the account.
//
// ctx: The context.Context for the request.
// declareTransaction: The input for adding a declare transaction.
// Returns: The response for adding a declare transaction and an error, if any.
func (account *Account) AddDeclareTransaction(ctx context.Context, declareTransaction rpc.AddDeclareTxnInput) (*rpc.AddDeclareTransactionResponse, error) {
	return account.provider.AddDeclareTransaction(ctx, declareTransaction)
}

// AddDeployAccountTransaction adds a deploy account transaction to the account.
//
// ctx: The context.Context object for the function.
// deployAccountTransaction: The rpc.DeployAccountTxn object representing the deploy account transaction.
// Returns a pointer to rpc.AddDeployAccountTransactionResponse and an error if any.
func (account *Account) AddDeployAccountTransaction(ctx context.Context, deployAccountTransaction rpc.DeployAccountTxn) (*rpc.AddDeployAccountTransactionResponse, error) {
	return account.provider.AddDeployAccountTransaction(ctx, deployAccountTransaction)
}

// BlockHashAndNumber returns the block hash and number for the account.
//
// ctx - The context in which the function is called.
// Returns the block hash and number as an rpc.BlockHashAndNumberOutput object.
// Returns an error if there was an issue retrieving the block hash and number.
func (account *Account) BlockHashAndNumber(ctx context.Context) (*rpc.BlockHashAndNumberOutput, error) {
	return account.provider.BlockHashAndNumber(ctx)
}

// BlockNumber returns the block number of the account.
//
// ctx - The context in which the function is called.
// Returns the block number as a uint64 and any error encountered.
func (account *Account) BlockNumber(ctx context.Context) (uint64, error) {
	return account.provider.BlockNumber(ctx)
}

// BlockTransactionCount returns the number of transactions in a block.
//
// ctx - The context.Context object for the function.
// blockID - The rpc.BlockID object representing the block.
// Returns the number of transactions in the block and an error, if any.
func (account *Account) BlockTransactionCount(ctx context.Context, blockID rpc.BlockID) (uint64, error) {
	return account.provider.BlockTransactionCount(ctx, blockID)
}

// BlockWithTxHashes retrieves a block with transaction hashes.
//
// ctx - the context.Context object for the request.
// blockID - the rpc.BlockID object specifying the block to retrieve.
// Returns an interface{} representing the retrieved block and an error if there was any.
func (account *Account) BlockWithTxHashes(ctx context.Context, blockID rpc.BlockID) (interface{}, error) {
	return account.provider.BlockWithTxHashes(ctx, blockID)
}

// BlockWithTxs retrieves the specified block along with its transactions.
//
// ctx: The context.Context object for the function.
// blockID: The rpc.BlockID parameter for the function.
// Returns: An interface{} and an error.
func (account *Account) BlockWithTxs(ctx context.Context, blockID rpc.BlockID) (interface{}, error) {
	return account.provider.BlockWithTxs(ctx, blockID)
}

// Call is a function that performs a function call on an Account.
//
// It takes in a context.Context object, a rpc.FunctionCall object, and a rpc.BlockID object as parameters.
// It returns a slice of *felt.Felt objects and an error object.
func (account *Account) Call(ctx context.Context, call rpc.FunctionCall, blockId rpc.BlockID) ([]*felt.Felt, error) {
	return account.provider.Call(ctx, call, blockId)
}

// ChainID returns the chain ID associated with the account.
//
// ctx: the context.Context object for the function.
// Returns:
//   - string: the chain ID.
//   - error: any error encountered while retrieving the chain ID.
func (account *Account) ChainID(ctx context.Context) (string, error) {
	return account.provider.ChainID(ctx)
}

// Class description of the Go function.
//
// Class is a method that calls the `Class` method of the `provider` field of the `account` struct.
// It takes a `context.Context` as the first parameter, a `rpc.BlockID` as the second parameter, and a `*felt.Felt` as the third parameter.
// It returns a `rpc.ClassOutput` and an `error`.
func (account *Account) Class(ctx context.Context, blockID rpc.BlockID, classHash *felt.Felt) (rpc.ClassOutput, error) {
	return account.provider.Class(ctx, blockID, classHash)
}

// ClassAt description of the Go function.
//
// ClassAt retrieves the class at the specified block ID and contract address.
// It takes the following parameters:
// - ctx: The context.Context object for the function.
// - blockID: The rpc.BlockID object representing the block ID.
// - contractAddress: The felt.Felt object representing the contract address.
// 
// It returns the rpc.ClassOutput object and an error.
func (account *Account) ClassAt(ctx context.Context, blockID rpc.BlockID, contractAddress *felt.Felt) (rpc.ClassOutput, error) {
	return account.provider.ClassAt(ctx, blockID, contractAddress)
}

// ClassHashAt returns the class hash at the given block ID for the specified contract address.
//
// ctx - The context to use for the function call.
// blockID - The ID of the block.
// contractAddress - The address of the contract to get the class hash for.
// Returns the class hash as a *felt.Felt and an error if any occurred.
func (account *Account) ClassHashAt(ctx context.Context, blockID rpc.BlockID, contractAddress *felt.Felt) (*felt.Felt, error) {
	return account.provider.ClassHashAt(ctx, blockID, contractAddress)
}

// EstimateFee estimates the fee for a set of requests in the given block ID.
//
// ctx: The context.Context object for the function.
// requests: An array of rpc.EstimateFeeInput objects representing the requests to estimate the fee for.
// blockID: The rpc.BlockID object representing the block ID for which to estimate the fee.
// []rpc.FeeEstimate: An array of rpc.FeeEstimate objects representing the estimated fees.
// error: An error object if any error occurred during the estimation process.
func (account *Account) EstimateFee(ctx context.Context, requests []rpc.EstimateFeeInput, blockID rpc.BlockID) ([]rpc.FeeEstimate, error) {
	return account.provider.EstimateFee(ctx, requests, blockID)
}

// EstimateMessageFee estimates the fee for a given message in the context of an account.
//
// ctx - The context.Context object for the function.
// msg - The rpc.MsgFromL1 object representing the message.
// blockID - The rpc.BlockID object representing the block ID.
// Returns a pointer to rpc.FeeEstimate and an error if any.
func (account *Account) EstimateMessageFee(ctx context.Context, msg rpc.MsgFromL1, blockID rpc.BlockID) (*rpc.FeeEstimate, error) {
	return account.provider.EstimateMessageFee(ctx, msg, blockID)
}

// Events retrieves events for the account.
//
// ctx: the context.Context to use for the request.
// input: the input parameters for retrieving events.
// Returns:
// - *rpc.EventChunk: the chunk of events retrieved.
// - error: an error if the retrieval fails.
func (account *Account) Events(ctx context.Context, input rpc.EventsInput) (*rpc.EventChunk, error) {
	return account.provider.Events(ctx, input)
}

// Nonce returns the nonce for the specified account and contract address.
//
// It takes the following parameters:
// - ctx: The context.Context object for the function.
// - blockID: The rpc.BlockID object for the function.
// - contractAddress: The felt.Felt object for the function.
//
// It returns a string pointer and an error.
func (account *Account) Nonce(ctx context.Context, blockID rpc.BlockID, contractAddress *felt.Felt) (*string, error) {
	return account.provider.Nonce(ctx, blockID, contractAddress)
}

// SimulateTransactions simulates transactions using the provided context, blockID, txns, and simulationFlags.
// It returns a list of simulated transactions and an error, if any.
func (account *Account) SimulateTransactions(ctx context.Context, blockID rpc.BlockID, txns []rpc.Transaction, simulationFlags []rpc.SimulationFlag) ([]rpc.SimulatedTransaction, error) {
	return account.provider.SimulateTransactions(ctx, blockID, txns, simulationFlags)
}

// StorageAt is a function that retrieves the storage value at the given key for a contract address.
//
// ctx: The context.Context object for the function.
// contractAddress: The contract address for which to retrieve the storage value.
// key: The key of the storage value to retrieve.
// blockID: The block ID at which to retrieve the storage value.
// string: The storage value at the given key.
// error: An error if the retrieval fails.
func (account *Account) StorageAt(ctx context.Context, contractAddress *felt.Felt, key string, blockID rpc.BlockID) (string, error) {
	return account.provider.StorageAt(ctx, contractAddress, key, blockID)
}

// StateUpdate updates the state of the Account.
//
// It takes a context.Context and a rpc.BlockID as parameters.
// It returns a *rpc.StateUpdateOutput and an error.
func (account *Account) StateUpdate(ctx context.Context, blockID rpc.BlockID) (*rpc.StateUpdateOutput, error) {
	return account.provider.StateUpdate(ctx, blockID)
}

// Syncing returns the sync status of the account.
//
// It takes a context.Context parameter and returns a *rpc.SyncStatus
// and an error.
func (account *Account) Syncing(ctx context.Context) (*rpc.SyncStatus, error) {
	return account.provider.Syncing(ctx)
}

// TraceBlockTransactions retrieves a list of trace transactions for a given block hash.
//
// ctx: The context.Context object.
// blockHash: The hash of the block to retrieve trace transactions for.
// []rpc.Trace: The list of trace transactions for the given block.
// error: An error if there was a problem retrieving the trace transactions.
func (account *Account) TraceBlockTransactions(ctx context.Context, blockHash *felt.Felt) ([]rpc.Trace, error) {
	return account.provider.TraceBlockTransactions(ctx, blockHash)
}

// TransactionReceipt retrieves the transaction receipt for the given transaction hash.
//
// ctx - The context to use for the request.
// transactionHash - The hash of the transaction.
// Return type: rpc.TransactionReceipt, error.
func (account *Account) TransactionReceipt(ctx context.Context, transactionHash *felt.Felt) (rpc.TransactionReceipt, error) {
	return account.provider.TransactionReceipt(ctx, transactionHash)
}

// TransactionTrace returns the transaction trace for a given transaction hash.
//
// ctx: The context.Context object for the request.
// transactionHash: The transaction hash for which the transaction trace is to be retrieved.
// Returns: The rpc.TxnTrace object representing the transaction trace, and an error if any.
func (account *Account) TransactionTrace(ctx context.Context, transactionHash *felt.Felt) (rpc.TxnTrace, error) {
	return account.provider.TransactionTrace(ctx, transactionHash)
}

// TransactionByBlockIdAndIndex returns a transaction by block ID and index.
//
// ctx - The context for the function.
// blockID - The ID of the block.
// index - The index of the transaction in the block.
// Returns the transaction and an error, if any.
func (account *Account) TransactionByBlockIdAndIndex(ctx context.Context, blockID rpc.BlockID, index uint64) (rpc.Transaction, error) {
	return account.provider.TransactionByBlockIdAndIndex(ctx, blockID, index)
}

// TransactionByHash returns the transaction with the given hash.
//
// It takes a context.Context and a *felt.Felt hash as parameters.
// It returns a rpc.Transaction and an error.
func (account *Account) TransactionByHash(ctx context.Context, hash *felt.Felt) (rpc.Transaction, error) {
	return account.provider.TransactionByHash(ctx, hash)
}

// FmtCalldata generates the formatted calldata for the given function calls and Cairo version.
//
// Parameters:
// - fnCalls: a slice of rpc.FunctionCall representing the function calls.
// - cairoVersion: an integer representing the Cairo version.
//
// Returns:
// - a slice of *felt.Felt representing the formatted calldata.
// - an error if Cairo version is not supported.
func (account *Account) FmtCalldata(fnCalls []rpc.FunctionCall, cairoVersion int) ([]*felt.Felt, error) {
	switch cairoVersion {
	case 0:
		return FmtCalldataCairo0(fnCalls), nil
	case 2:
		return FmtCalldataCairo2(fnCalls), nil
	default:
		return nil, errors.New("Cairo version not supported")
	}
}

// FmtCalldataCairo0 generates a slice of *felt.Felt that represents the calldata for the given function calls in Cairo 0 format.
//
// The function takes in a slice of rpc.FunctionCall named fnCalls as a parameter.
// It returns a slice of *felt.Felt.
func FmtCalldataCairo0(fnCalls []rpc.FunctionCall) []*felt.Felt {
	execCallData := []*felt.Felt{}
	execCallData = append(execCallData, new(felt.Felt).SetUint64(uint64(len(fnCalls))))

	// Legacy : Cairo 0
	concatCallData := []*felt.Felt{}
	for _, fnCall := range fnCalls {
		execCallData = append(
			execCallData,
			fnCall.ContractAddress,
			fnCall.EntryPointSelector,
			new(felt.Felt).SetUint64(uint64(len(concatCallData))),
			new(felt.Felt).SetUint64(uint64(len(fnCall.Calldata))+1),
		)
		concatCallData = append(concatCallData, fnCall.Calldata...)
	}
	execCallData = append(execCallData, new(felt.Felt).SetUint64(uint64(len(concatCallData))+1))
	execCallData = append(execCallData, concatCallData...)
	execCallData = append(execCallData, new(felt.Felt).SetUint64(0))

	return execCallData
}

// FmtCalldataCairo2 generates the calldata for the given function calls for Cairo 2 contracs.
//
// Parameters:
// - fnCalls: a slice of rpc.FunctionCall containing the function calls.
//
// Return type:
// - a slice of *felt.Felt representing the generated calldata.
func FmtCalldataCairo2(fnCalls []rpc.FunctionCall) []*felt.Felt {
	execCallData := []*felt.Felt{}
	execCallData = append(execCallData, new(felt.Felt).SetUint64(uint64(len(fnCalls))))

	concatCallData := []*felt.Felt{}
	for _, fnCall := range fnCalls {
		execCallData = append(
			execCallData,
			fnCall.ContractAddress,
			fnCall.EntryPointSelector,
			new(felt.Felt).SetUint64(uint64(len(concatCallData))),
			new(felt.Felt).SetUint64(uint64(len(fnCall.Calldata))),
		)
		concatCallData = append(concatCallData, fnCall.Calldata...)
	}
	execCallData = append(execCallData, new(felt.Felt).SetUint64(uint64(len(concatCallData))))
	execCallData = append(execCallData, concatCallData...)

	return execCallData
}
