package account

import (
	"context"
	"errors"

	"github.com/NethermindEth/juno/core/felt"
	starknetgo "github.com/NethermindEth/starknet.go"
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
	TRANSACTION_PREFIX      = new(felt.Felt).SetBytes([]byte("invoke"))
	DECLARE_PREFIX          = new(felt.Felt).SetBytes([]byte("declare"))
	CONTRACT_ADDRESS_PREFIX = new(felt.Felt).SetBytes([]byte("STARKNET_CONTRACT_ADDRESS"))
)

//go:generate mockgen -destination=../mocks/mock_account.go -package=mocks -source=account.go AccountInterface
type AccountInterface interface {
	Sign(ctx context.Context, msg *felt.Felt) ([]*felt.Felt, error)
	BuildInvokeTx(ctx context.Context, invokeTx *rpc.InvokeTxnV1, fnCall *[]rpc.FunctionCall) error
	TransactionHashInvoke(invokeTxn rpc.InvokeTxnType) (*felt.Felt, error)
	TransactionHashDeployAccount(tx rpc.DeployAccountTxn, contractAddress *felt.Felt) (*felt.Felt, error)
	TransactionHashDeclare(tx rpc.DeclareTxnType) (*felt.Felt, error)
	SignInvokeTransaction(ctx context.Context, tx *rpc.InvokeTxnV1) error
	SignDeployAccountTransaction(ctx context.Context, tx *rpc.DeployAccountTxn, precomputeAddress *felt.Felt) error
	SignDeclareTransaction(ctx context.Context, tx *rpc.DeclareTxnV2) error
	PrecomputeAddress(deployerAddress *felt.Felt, salt *felt.Felt, classHash *felt.Felt, constructorCalldata []*felt.Felt) (*felt.Felt, error)
}

var _ AccountInterface = &Account{}
var _ rpc.RpcProvider = &Account{}

type Account struct {
	provider       rpc.RpcProvider
	ChainId        *felt.Felt
	AccountAddress *felt.Felt
	publicKey      string
	ks             starknetgo.Keystore
}

func NewAccount(provider rpc.RpcProvider, accountAddress *felt.Felt, publicKey string, keystore starknetgo.Keystore) (*Account, error) {
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

func (account *Account) Sign(ctx context.Context, msg *felt.Felt) ([]*felt.Felt, error) {

	msgBig, ok := utils.FeltToBigInt(msg)
	if ok != true {
		return nil, ErrFeltToBigInt
	}
	s1, s2, err := account.ks.Sign(ctx, account.publicKey, msgBig)
	if err != nil {
		return nil, err
	}
	s1Felt, _ := utils.BigIntToFelt(s1)
	s2Felt, _ := utils.BigIntToFelt(s2)

	return []*felt.Felt{s1Felt, s2Felt}, nil
}

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

// TransactionHashDeployAccount computes the transaction hash for deployAccount transactions
func (account *Account) TransactionHashDeployAccount(tx rpc.DeployAccountTxn, contractAddress *felt.Felt) (*felt.Felt, error) {

	// https://docs.starknet.io/documentation/architecture_and_concepts/Network_Architecture/transactions/#deploy_account_transaction

	// There is only version 1 of deployAccount txn
	if tx.Version != rpc.TransactionV1 {
		return nil, ErrTxnTypeUnSupported
	}

	Prefix_DEPLOY_ACCOUNT := new(felt.Felt).SetBytes([]byte("deploy_account"))

	calldata := []*felt.Felt{tx.ClassHash, tx.ContractAddressSalt}
	calldata = append(calldata, tx.ConstructorCalldata...)
	calldataHash, err := computeHashOnElementsFelt(calldata)
	if err != nil {
		return nil, err
	}

	versionFelt, err := new(felt.Felt).SetString(string(tx.Version))
	if err != nil {
		return nil, err
	}

	// https://docs.starknet.io/documentation/architecture_and_concepts/Network_Architecture/transactions/#deploy_account_hash_calculation
	return calculateTransactionHashCommon(
		Prefix_DEPLOY_ACCOUNT,
		versionFelt,
		contractAddress,
		&felt.Zero,
		calldataHash,
		tx.MaxFee,
		account.ChainId,
		[]*felt.Felt{tx.Nonce},
	)
}

func (account *Account) TransactionHashInvoke(tx rpc.InvokeTxnType) (*felt.Felt, error) {

	// https://docs.starknet.io/documentation/architecture_and_concepts/Network_Architecture/transactions/#deploy_account_hash_calculation
	switch txn := tx.(type) {
	case rpc.InvokeTxnV0:
		if txn.Version == "" || len(txn.Calldata) == 0 || txn.MaxFee == nil || txn.EntryPointSelector == nil {
			return nil, ErrNotAllParametersSet
		}

		calldataHash, err := computeHashOnElementsFelt(txn.Calldata)
		if err != nil {
			return nil, err
		}

		txnVersionFelt, err := new(felt.Felt).SetString(string(txn.Version))
		if err != nil {
			return nil, err
		}
		return calculateTransactionHashCommon(
			TRANSACTION_PREFIX,
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

		calldataHash, err := computeHashOnElementsFelt(txn.Calldata)
		if err != nil {
			return nil, err
		}
		txnVersionFelt, err := new(felt.Felt).SetString(string(txn.Version))
		if err != nil {
			return nil, err
		}
		return calculateTransactionHashCommon(
			TRANSACTION_PREFIX,
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

func (account *Account) TransactionHashDeclare(tx rpc.DeclareTxnType) (*felt.Felt, error) {

	Prefix_DECLARE := new(felt.Felt).SetBytes([]byte("declare"))

	switch txn := tx.(type) {
	case rpc.DeclareTxnV0:
		// Due to inconsistencies in version 0 hash calculation we don't calculate the hash
		return nil, ErrTxnVersionUnSupported
	case rpc.DeclareTxnV1:
		if txn.SenderAddress == nil || txn.Version == "" || txn.ClassHash == nil || txn.MaxFee == nil || txn.Nonce == nil {
			return nil, ErrNotAllParametersSet
		}

		calldataHash, err := computeHashOnElementsFelt([]*felt.Felt{txn.ClassHash})
		if err != nil {
			return nil, err
		}

		txnVersionFelt, err := new(felt.Felt).SetString(string(txn.Version))
		if err != nil {
			return nil, err
		}
		return calculateTransactionHashCommon(
			Prefix_DECLARE,
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

		calldataHash, err := computeHashOnElementsFelt([]*felt.Felt{txn.ClassHash})
		if err != nil {
			return nil, err
		}

		txnVersionFelt, err := new(felt.Felt).SetString(string(txn.Version))
		if err != nil {
			return nil, err
		}
		return calculateTransactionHashCommon(
			Prefix_DECLARE,
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

// precomputeAddress precomputes the accounts address
// ref: https://github.com/starkware-libs/cairo-lang/blob/master/src/starkware/starknet/core/os/contract_address/contract_address.py
func (account *Account) PrecomputeAddress(deployerAddress *felt.Felt, salt *felt.Felt, classHash *felt.Felt, constructorCalldata []*felt.Felt) (*felt.Felt, error) {

	bigIntArr, err := utils.FeltArrToBigIntArr([]*felt.Felt{
		CONTRACT_ADDRESS_PREFIX,
		deployerAddress,
		salt,
		classHash,
	})
	if err != nil {
		return nil, err
	}

	constructorCalldataBigIntArr, err := utils.FeltArrToBigIntArr(constructorCalldata)
	constructorCallDataHashInt, _ := starknetgo.Curve.ComputeHashOnElements(*constructorCalldataBigIntArr)
	*bigIntArr = append(*bigIntArr, constructorCallDataHashInt)

	preBigInt, err := starknetgo.Curve.ComputeHashOnElements(*bigIntArr)
	if err != nil {
		return nil, err
	}
	return utils.BigIntToFelt(preBigInt)

}

// BuildInvokeTx formats the calldata and signs the transaction
func (account *Account) BuildInvokeTx(ctx context.Context, invokeTx *rpc.InvokeTxnV1, fnCall *[]rpc.FunctionCall) error {
	invokeTx.Calldata = FmtCalldata(*fnCall)
	return account.SignInvokeTransaction(ctx, invokeTx)
}

func (account *Account) AddInvokeTransaction(ctx context.Context, invokeTx rpc.InvokeTxnV1) (*rpc.AddInvokeTransactionResponse, error) {
	return account.provider.AddInvokeTransaction(ctx, invokeTx)
}

func (account *Account) AddDeclareTransaction(ctx context.Context, declareTransaction rpc.AddDeclareTxnInput) (*rpc.AddDeclareTransactionResponse, error) {
	return account.provider.AddDeclareTransaction(ctx, declareTransaction)
}

func (account *Account) AddDeployAccountTransaction(ctx context.Context, deployAccountTransaction rpc.DeployAccountTxn) (*rpc.AddDeployAccountTransactionResponse, error) {
	return account.provider.AddDeployAccountTransaction(ctx, deployAccountTransaction)
}

func (account *Account) BlockHashAndNumber(ctx context.Context) (*rpc.BlockHashAndNumberOutput, error) {
	return account.provider.BlockHashAndNumber(ctx)
}

func (account *Account) BlockNumber(ctx context.Context) (uint64, error) {
	return account.provider.BlockNumber(ctx)
}

func (account *Account) BlockTransactionCount(ctx context.Context, blockID rpc.BlockID) (uint64, error) {
	return account.provider.BlockTransactionCount(ctx, blockID)
}

func (account *Account) BlockWithTxHashes(ctx context.Context, blockID rpc.BlockID) (interface{}, error) {
	return account.provider.BlockWithTxHashes(ctx, blockID)
}

func (account *Account) BlockWithTxs(ctx context.Context, blockID rpc.BlockID) (interface{}, error) {
	return account.provider.BlockWithTxs(ctx, blockID)
}

func (account *Account) Call(ctx context.Context, call rpc.FunctionCall, blockId rpc.BlockID) ([]*felt.Felt, error) {
	return account.provider.Call(ctx, call, blockId)
}

func (account *Account) ChainID(ctx context.Context) (string, error) {
	return account.provider.ChainID(ctx)
}
func (account *Account) Class(ctx context.Context, blockID rpc.BlockID, classHash *felt.Felt) (rpc.ClassOutput, error) {
	return account.provider.Class(ctx, blockID, classHash)
}
func (account *Account) ClassAt(ctx context.Context, blockID rpc.BlockID, contractAddress *felt.Felt) (rpc.ClassOutput, error) {
	return account.provider.ClassAt(ctx, blockID, contractAddress)
}

func (account *Account) ClassHashAt(ctx context.Context, blockID rpc.BlockID, contractAddress *felt.Felt) (*felt.Felt, error) {
	return account.provider.ClassHashAt(ctx, blockID, contractAddress)
}

func (account *Account) EstimateFee(ctx context.Context, requests []rpc.EstimateFeeInput, blockID rpc.BlockID) ([]rpc.FeeEstimate, error) {
	return account.provider.EstimateFee(ctx, requests, blockID)
}
func (account *Account) EstimateMessageFee(ctx context.Context, msg rpc.MsgFromL1, blockID rpc.BlockID) (*rpc.FeeEstimate, error) {
	return account.provider.EstimateMessageFee(ctx, msg, blockID)
}

func (account *Account) Events(ctx context.Context, input rpc.EventsInput) (*rpc.EventChunk, error) {
	return account.provider.Events(ctx, input)
}
func (account *Account) Nonce(ctx context.Context, blockID rpc.BlockID, contractAddress *felt.Felt) (*string, error) {
	return account.provider.Nonce(ctx, blockID, contractAddress)
}

func (account *Account) SimulateTransactions(ctx context.Context, blockID rpc.BlockID, txns []rpc.Transaction, simulationFlags []rpc.SimulationFlag) ([]rpc.SimulatedTransaction, error) {
	return account.provider.SimulateTransactions(ctx, blockID, txns, simulationFlags)
}
func (account *Account) StorageAt(ctx context.Context, contractAddress *felt.Felt, key string, blockID rpc.BlockID) (string, error) {
	return account.provider.StorageAt(ctx, contractAddress, key, blockID)
}
func (account *Account) StateUpdate(ctx context.Context, blockID rpc.BlockID) (*rpc.StateUpdateOutput, error) {
	return account.provider.StateUpdate(ctx, blockID)
}
func (account *Account) Syncing(ctx context.Context) (*rpc.SyncStatus, error) {
	return account.provider.Syncing(ctx)
}

func (account *Account) TraceBlockTransactions(ctx context.Context, blockHash *felt.Felt) ([]rpc.Trace, error) {
	return account.provider.TraceBlockTransactions(ctx, blockHash)
}

func (account *Account) TransactionReceipt(ctx context.Context, transactionHash *felt.Felt) (rpc.TransactionReceipt, error) {
	return account.provider.TransactionReceipt(ctx, transactionHash)
}

func (account *Account) TransactionTrace(ctx context.Context, transactionHash *felt.Felt) (rpc.TxnTrace, error) {
	return account.provider.TransactionTrace(ctx, transactionHash)
}

func (account *Account) TransactionByBlockIdAndIndex(ctx context.Context, blockID rpc.BlockID, index uint64) (rpc.Transaction, error) {
	return account.provider.TransactionByBlockIdAndIndex(ctx, blockID, index)
}

func (account *Account) TransactionByHash(ctx context.Context, hash *felt.Felt) (rpc.Transaction, error) {
	return account.provider.TransactionByHash(ctx, hash)
}
