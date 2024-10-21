package typedData

import "math/big"

type (
	Felt            string
	Bool            bool
	String          string
	Selector        string
	U128            *big.Int
	I128            *big.Int
	ContractAddress string
	ClassHash       string
	Timestamp       U128
	Shortstring     string
)

type U256 struct {
	Low  U128 `json:"low"`
	High U128 `json:"high"`
}

type TokenAmount struct {
	TokenAddress ContractAddress `json:"token_address"`
	Amount       U256            `json:"amount"`
}

type NftId struct {
	CollectionAddress ContractAddress `json:"collection_address"`
	TokenID           U256            `json:"token_id"`
}
