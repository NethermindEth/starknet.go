package caigo

import (
	"math/big"

	"github.com/dontpanicdao/caigo/curve"
	"github.com/dontpanicdao/caigo/types"
)

// Adheres to 'starknet.js' hash non typedData
func HashMsg(addr *big.Int, tx types.Transaction) (hash *big.Int, err error) {
	calldataArray := []*big.Int{big.NewInt(int64(len(tx.Calldata)))}
	for _, cd := range tx.Calldata {
		calldataArray = append(calldataArray, HexToBN(cd))
	}

	cdHash, err := HashElements(calldataArray)
	if err != nil {
		return hash, err
	}

	txHashData := []*big.Int{
		addr,
		SNValToBN(tx.ContractAddress),
		GetSelectorFromName(tx.EntryPointSelector),
		cdHash,
		SNValToBN(tx.Nonce),
	}

	hash, err = ComputeHashOnElements(txHashData)
	return hash, err
}

// Adheres to 'starknet.js' hash non typedData
func HashTx(addr *big.Int, tx types.Transaction) (hash *big.Int, err error) {
	calldataArray := []*big.Int{big.NewInt(int64(len(tx.Calldata)))}
	for _, cd := range tx.Calldata {
		calldataArray = append(calldataArray, SNValToBN(cd))
	}

	cdHash, err := HashElements(calldataArray)
	if err != nil {
		return hash, err
	}

	txHashData := []*big.Int{
		SNValToBN(tx.ContractAddress),
		GetSelectorFromName(tx.EntryPointSelector),
		cdHash,
	}

	hash, err = ComputeHashOnElements(txHashData)
	return hash, err
}

/*
	Hashes the contents of a given array with its size using a golang Pedersen Hash implementation.

	(ref: https://github.com/starkware-libs/cairo-lang/blob/13cef109cd811474de114925ee61fd5ac84a25eb/src/starkware/cairo/common/hash_state.py#L6)
*/
func ComputeHashOnElements(elems []*big.Int) (hash *big.Int, err error) {
	elems = append(elems, big.NewInt(int64(len(elems))))
	return HashElements((elems))
}

/*
	Hashes the contents of a given array using a golang Pedersen Hash implementation.

	(ref: https://github.com/seanjameshan/starknet.js/blob/main/src/utils/ellipticCurve.ts)
*/
func HashElements(elems []*big.Int) (hash *big.Int, err error) {
	if len(elems) == 0 {
		elems = append(elems, big.NewInt(0))
	}

	hash = big.NewInt(0)
	for _, h := range elems {
		hash, err = curve.PedersenHash([]*big.Int{hash, h})
		if err != nil {
			return hash, err
		}
	}
	return hash, err
}
