package paymaster

import "context"

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
