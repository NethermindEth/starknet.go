package caigo

import (
	"math/big"
)

// struct to catch starknet.js transaction payloads
type JSTransaction struct {
	Calldata           []string `json:"calldata"`
	ContractAddress    string   `json:"contract_address"`
	EntryPointSelector string   `json:"entry_point_selector"`
	EntryPointType     string   `json:"entry_point_type"`
	JSSignature        []string `json:"signature"`
	TransactionHash    string   `json:"transaction_hash"`
	Type               string   `json:"type"`
	Nonce              string   `json:"nonce,omitempty"`
}

type Transaction struct {
	Calldata           []*big.Int `json:"calldata"`
	ContractAddress    *big.Int   `json:"contract_address"`
	EntryPointSelector *big.Int   `json:"entry_point_selector"`
	EntryPointType     string   `json:"entry_point_type"`
	Signature          []*big.Int `json:"signature"`
	TransactionHash    *big.Int   `json:"transaction_hash"`
	Type               string   `json:"type"`
	Nonce              *big.Int   `json:"nonce,omitempty"`
}

func (jtx JSTransaction) ConvertTx() (tx Transaction, err error) {
	tx = Transaction{
		ContractAddress: HexToBN(jtx.ContractAddress),
		EntryPointSelector: HexToBN(jtx.EntryPointSelector),
		EntryPointType: jtx.EntryPointType,
		TransactionHash: HexToBN(jtx.TransactionHash),
		Type: jtx.Type,
		Nonce: HexToBN(jtx.Nonce),
	}
	for _, cd := range jtx.Calldata {
		tx.Calldata = append(tx.Calldata, StrToBig(cd))
	}
	for _, sigElem := range jtx.JSSignature {
		tx.Signature = append(tx.Signature, StrToBig(sigElem))
	}
	return tx, err
}

func (tx Transaction) HashTx(pubkey *big.Int, sc StarkCurve) (hash *big.Int, err error) {
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
