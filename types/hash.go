package types

import (
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
)

// Adheres to 'starknet.js' hash non typedData
func (sc StarkCurve) HashTx(addr *big.Int, tx types.Transaction) (hash *big.Int, err error) {
	calldataArray := []*big.Int{big.NewInt(int64(len(tx.Calldata)))}
	for _, cd := range tx.Calldata {
		calldataArray = append(calldataArray, SNValToBN(cd.String()))
	}

	cdHash, err := sc.HashElements(calldataArray)
	if err != nil {
		return hash, err
	}

	txHashData := []*big.Int{
		SNValToBN(tx.ContractAddress.String()),
		tx.EntryPointSelector.Int,
		cdHash,
	}

	hash, err = sc.ComputeHashOnElements(txHashData)
	return hash, err
}

// Adheres to 'starknet.js' hash non typedData
func (sc StarkCurve) HashMsg(addr *big.Int, tx types.Transaction) (hash *big.Int, err error) {
	calldataArray := []*big.Int{big.NewInt(int64(len(tx.Calldata)))}
	for _, cd := range tx.Calldata {
		calldataArray = append(calldataArray, HexToBN(cd.String()))
	}

	cdHash, err := sc.HashElements(calldataArray)
	if err != nil {
		return hash, err
	}

	txHashData := []*big.Int{
		addr,
		SNValToBN(tx.ContractAddress.String()),
		tx.EntryPointSelector.Int,
		cdHash,
		SNValToBN(tx.Nonce.String()),
	}

	hash, err = sc.ComputeHashOnElements(txHashData)
	return hash, err
}
