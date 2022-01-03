package caigo

type JSTransaction struct {
	Calldata []string `json:"calldata"`
	ContractAddress string `json:"contract_address"`
	EntryPointSelector string `json:"entry_point_selector"`
	EntryPointType string `json:"entry_point_type"`
	JSSignature []string `json:"signature"`
	TransactionHash string `json:"transaction_hash"`
	Type string `json:"type"`
	Nonce string `json:"nonce,omitempty"`
}

// func(tx Transaction) HashTx(pubkey string) (hash *big.Int, err error) {
// 	var bnCallData []*big.Int
// 	for _, cd := range tx.Calldata {
// 		bnCallData = append(bnCallData, StrToBig(cd))
// 	}
// 	cdHash := hashElements(bnCallData)
// }

