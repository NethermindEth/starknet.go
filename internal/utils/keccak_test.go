package utils

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestGetSelectorFromName tests the GetSelectorFromName function.
//
// It checks if the GetSelectorFromName function returns the expected values
// for different input names.
// The expected values are hard-coded and compared against the actual values.
// If any of the actual values do not match the expected values, an error is
// reported.
//
// Parameters:
//   - t: The testing.T object used for reporting test failures and logging test output
//
// Returns:
//
//	none
func TestGetSelectorFromName(t *testing.T) {
	sel1 := BigToHex(GetSelectorFromName("initialise"))
	sel2 := BigToHex(GetSelectorFromName("mint"))
	sel3 := BigToHex(GetSelectorFromName("test"))

	exp1 := "0xa899ecd8428376a1913bc783548272a3206c3079bcdc339f7e080e8c6ddfae"
	exp2 := "0x2f0b3c5710379609eb5495f1ecd348cb28167711b73609fe565a72734550354"
	exp3 := "0x22ff5f21f0b81b113e63f7db6da94fedef11b2119b4088b89664fb9a3cb658"

	if sel1 != exp1 || sel2 != exp2 || sel3 != exp3 {
		t.Errorf("invalid Keccak256 encoding: %v %v %v\n", sel1, sel2, sel3)
	}
}

// TestComputeFact tests the ComputeFact function.
//
// It tests the ComputeFact function by providing a set of test cases
// and comparing the computed hash with the expected hash.
// The test cases consist of program hashes, program outputs,
// and expected hash values.
//
// Parameters:
//   - t: The testing.T object for running the test
//
// Returns:
//
//	none
func TestComputeFact(t *testing.T) {
	testFacts := []struct {
		programHash   *big.Int
		programOutput []*big.Int
		expected      *big.Int
	}{
		{
			programHash: HexToBN(
				"0x114952172aed91e59f870a314e75de0a437ff550e4618068cec2d832e48b0c7",
			),
			programOutput: []*big.Int{big.NewInt(289)},
			expected: HexToBN(
				"0xe6168c0a865aa80d724ad05627fa65fbcfe4b1d66a586e9f348f461b076072c4",
			),
		},
		{
			programHash: HexToBN(
				"0x79920d895101ad1fbdea9adf141d8f362fdea9ee35f33dfcd07f38e4a589bab",
			),
			programOutput: []*big.Int{
				StrToBig(
					"2754806153357301156380357983574496185342034785016738734224771556919270737441",
				),
			},
			expected: HexToBN("0x1d174fa1443deea9aab54bbca8d9be308dd14a0323dd827556c173bd132098db"),
		},
	}

	for _, tt := range testFacts {
		hash := ComputeFact(tt.programHash, tt.programOutput)
		assert.Equal(t, tt.expected, hash)
	}
}

// TestSplitFactStr is a test function that tests the SplitFactStr function.
//
// It verifies the behaviour of the SplitFactStr function by providing different inputs and checking the output.
// The function takes no parameters and returns no values.
//
// Parameters:
//   - t: The testing.T object for running the test
//
// Returns:
//
//	none
func TestSplitFactStr(t *testing.T) {
	type tescase struct {
		input string
		h     string
		l     string
		err   bool
	}
	data := []tescase{
		{
			input: "0x3",
			h:     "0x0",
			l:     "0x3",
		},
		{
			input: "0x300000000000000000000000000000000",
			h:     "0x3",
			l:     "0x0",
		},
		{
			input: "11111111111111111111111111111111111111111111111111111111111111010",
			err:   true,
		},
		{
			input: "X",
			err:   true,
		},
	}
	for _, d := range data {
		l, h, err := SplitFactStr(d.input)
		if d.err {
			require.Error(t, err)
		} else {
			require.NoError(t, err)
			assert.Equal(t, d.l, l)
			assert.Equal(t, d.h, h)
		}
	}
}
