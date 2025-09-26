package paymaster

import "context"

// IsAvailable returns the status of the paymaster service.
// If the paymaster service is correctly functioning, return true. Else, return false
//
// Parameters:
//   - ctx: The context.Context object for controlling the function call
//
// Returns:
//   - bool: True if the paymaster service is correctly functioning, false otherwise
//   - error: An error if any
func (p *Paymaster) IsAvailable(ctx context.Context) (bool, error) {
	var response bool
	if err := p.c.CallContextWithSliceArgs(ctx, &response, "paymaster_isAvailable"); err != nil {
		return false, err
	}

	return response, nil
}
