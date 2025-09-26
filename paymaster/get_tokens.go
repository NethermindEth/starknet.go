package paymaster

import (
	"context"

	"github.com/NethermindEth/juno/core/felt"
)

// Get a list of the tokens that the paymaster supports, together with their prices in STRK
//
// Parameters:
//   - ctx: The context.Context object for controlling the function call
//
// Returns:
//   - []TokenData: An array of token data
//   - error: An error if any
func (p *Paymaster) GetSupportedTokens(ctx context.Context) ([]TokenData, error) {
	var response []TokenData
	if err := p.c.CallContextWithSliceArgs(ctx, &response, "paymaster_getSupportedTokens"); err != nil {
		return nil, err
	}

	return response, nil
}

// Object containing data about the token: contract address, number of decimals and current price in STRK
type TokenData struct {
	// Token contract address
	TokenAddress *felt.Felt `json:"token_address"`
	// The number of decimals of the token
	Decimals uint8 `json:"decimals"`
	// Price in STRK (in FRI units)
	PriceInStrk string `json:"price_in_strk"` // u256 as a hex string
}
