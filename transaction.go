package caigo

import (
	"math/big"
)

// Starknet transaction composition
type Transaction struct {
	Calldata           []*big.Int `json:"calldata"`
	ContractAddress    *big.Int   `json:"contract_address"`
	EntryPointSelector *big.Int   `json:"entry_point_selector"`
	EntryPointType     string     `json:"entry_point_type"`
	Signature          []*big.Int `json:"signature"`
	TransactionHash    *big.Int   `json:"transaction_hash"`
	Type               string     `json:"type"`
	Nonce              *big.Int   `json:"nonce,omitempty"`
}

// Adheres to 'starknet.js' hash non typedData
func (sc StarkCurve) HashTx(pubkey *big.Int, tx Transaction) (hash *big.Int, err error) {
	tx.Calldata = append(tx.Calldata, big.NewInt(int64(len(tx.Calldata))))
	cdHash, err := sc.HashElements(tx.Calldata)
	if err != nil {
		return hash, err
	}

	txHashData := []*big.Int{
		pubkey,
		tx.ContractAddress,
		tx.EntryPointSelector,
		cdHash,
		tx.Nonce,
	}
	txHashData = append(txHashData, big.NewInt(int64(len(txHashData))))
	hash, err = sc.HashElements(txHashData)
	return hash, err
}
