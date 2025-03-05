package utils

import "testing"

// TestGetSelectorFromName tests the GetSelectorFromName function.
//
// It checks if the GetSelectorFromName function returns the expected values
// for different input names.
// The expected values are hard-coded and compared against the actual values.
// If any of the actual values do not match the expected values, an error is
// reported.
//
// Parameters:
// - t: The testing.T object used for reporting test failures and logging test output
// Returns:
//
//	none
func TestGetSelectorFromName(t *testing.T) {
	sel1 := BigToHex(GetSelectorFromName("initialize"))
	sel2 := BigToHex(GetSelectorFromName("mint"))
	sel3 := BigToHex(GetSelectorFromName("test"))

	exp1 := "0x79dc0da7c54b95f10aa182ad0a46400db63156920adb65eca2654c0945a463"
	exp2 := "0x2f0b3c5710379609eb5495f1ecd348cb28167711b73609fe565a72734550354"
	exp3 := "0x22ff5f21f0b81b113e63f7db6da94fedef11b2119b4088b89664fb9a3cb658"

	if sel1 != exp1 || sel2 != exp2 || sel3 != exp3 {
		t.Errorf("invalid Keccak256 encoding: %v %v %v\n", sel1, sel2, sel3)
	}
}
