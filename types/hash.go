package types

import (
	"math/big"

	"github.com/dontpanicdao/caigo/felt"
)

type Curve felt.StarkCurve

// Adheres to 'starknet.js' hash non typedData
func (sc Curve) HashTx(addr felt.Felt, tx Transaction) (hash *felt.Felt, err error) {
	calldataArray := []felt.Felt{felt.BigToFelt(big.NewInt(int64(len(tx.Calldata))))}
	calldataArray = append(calldataArray, tx.Calldata...)

	cdHash, err := felt.StarkCurve(sc).HashElements(calldataArray)
	if err != nil {
		return hash, err
	}

	txHashData := []felt.Felt{
		tx.ContractAddress,
	}
	if tx.EntryPointSelector != nil {
		txHashData = append(txHashData, *tx.EntryPointSelector)
	}
	if cdHash != nil {
		txHashData = append(txHashData, *cdHash)
	}
	return felt.StarkCurve(sc).ComputeHashOnElements(txHashData)
}

// Adheres to 'starknet.js' hash non typedData
func (sc Curve) HashMsg(addr felt.Felt, tx Transaction) (hash *felt.Felt, err error) {
	calldataArray := []felt.Felt{felt.BigToFelt(big.NewInt(int64(len(tx.Calldata))))}
	calldataArray = append(calldataArray, tx.Calldata...)

	cdHash, err := felt.StarkCurve(sc).HashElements(calldataArray)
	if err != nil {
		return hash, err
	}

	txHashData := []felt.Felt{
		addr,
		tx.ContractAddress,
	}
	if tx.EntryPointSelector != nil {
		txHashData = append(txHashData, *tx.EntryPointSelector)
	}
	if cdHash != nil {
		txHashData = append(txHashData, *cdHash)
	}
	if tx.Nonce != nil {
		txHashData = append(txHashData, *tx.Nonce)

	}
	return felt.StarkCurve(sc).ComputeHashOnElements(txHashData)
}
