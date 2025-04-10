// Package types provides custom types and utilities for StarkNet type handling.
package types

import (
	"fmt"
	"strings"

	"github.com/NethermindEth/juno/core/felt"
)

// Felt is a wrapper around the Juno felt.Felt type, which represents a field element in the StarkNet network.
// This wrapper allows us to extend the functionality of the original type with StarkNet.go specific methods.
type Felt struct {
	*felt.Felt
}

// NewFelt creates a new Felt from a Juno felt.Felt.
func NewFelt(f *felt.Felt) *Felt {
	return &Felt{Felt: f}
}

// AddressStr returns a fixed-length hexadecimal representation of a Felt value, ensuring it has
// exactly 66 characters (including the "0x" prefix) by padding with leading zeros if necessary.
// This is especially important for addresses and class hashes in StarkNet, which should always
// have a consistent length representation.
//
// The standard String() method from the Juno felt.Felt type does not pad with leading zeros,
// which can lead to inconsistencies when displaying or processing address strings.
//
// Returns:
// - string: a hex string of exactly 66 characters (including 0x prefix)
func (f *Felt) AddressStr() string {
	hexStr := f.String()
	
	// If already 66 characters, return as is
	if len(hexStr) == 66 {
		return hexStr
	}
	
	// Remove "0x" prefix
	hexValue := strings.TrimPrefix(hexStr, "0x")
	
	// Pad with leading zeros to make it exactly 64 characters
	return fmt.Sprintf("0x%064s", hexValue)
}
