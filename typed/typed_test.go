package typed

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/NethermindEth/starknet.go/utils"
	"github.com/stretchr/testify/require"
)

type Mail struct {
	From     Person
	To       Person
	Contents string
}

type Person struct {
	Name   string
	Wallet string
}

// FmtDefinitionEncoding formats the encoding for the given field in the Mail struct.
//
// Parameters:
// - field: the field to format the encoding for
// Returns:
// - fmtEnc: a slice of big integers
func (mail Mail) FmtDefinitionEncoding(field string) (fmtEnc []*big.Int) {
	if field == "from" {
		fmtEnc = append(fmtEnc, utils.UTF8StrToBig(mail.From.Name))
		fmtEnc = append(fmtEnc, utils.HexToBN(mail.From.Wallet))
	} else if field == "to" {
		fmtEnc = append(fmtEnc, utils.UTF8StrToBig(mail.To.Name))
		fmtEnc = append(fmtEnc, utils.HexToBN(mail.To.Wallet))
	} else if field == "contents" {
		fmtEnc = append(fmtEnc, utils.UTF8StrToBig(mail.Contents))
	}
	return fmtEnc
}

// MockTypedData generates a TypedData object for testing purposes.
// It creates example types and initializes a Domain object. Then it uses the example types and the domain to create a new TypedData object.
// The function returns the generated TypedData object.
//
// Parameters:
//
//	none
//
// Returns:
// - ttd: the generated TypedData object
func MockTypedData() (ttd TypedData, err error) {
	exampleTypes := make(map[string]TypeDef)
	domDefs := []Definition{{"name", "felt"}, {"version", "felt"}, {"chainId", "felt"}}
	exampleTypes["StarkNetDomain"] = TypeDef{Definitions: domDefs}
	mailDefs := []Definition{{"from", "Person"}, {"to", "Person"}, {"contents", "felt"}}
	exampleTypes["Mail"] = TypeDef{Definitions: mailDefs}
	persDefs := []Definition{{"name", "felt"}, {"wallet", "felt"}}
	exampleTypes["Person"] = TypeDef{Definitions: persDefs}

	dm := Domain{
		Name:    "StarkNet Mail",
		Version: "1",
		ChainId: "1",
	}

	ttd, err = NewTypedData(exampleTypes, "Mail", dm)
	if err != nil {
		return TypedData{}, err
	}
	return ttd, err
}

// TestGeneral_GetMessageHash tests the GetMessageHash function.
//
// It creates a mock TypedData and sets up a test case for hashing a mail message.
// The mail message contains information about the sender and recipient, as well as the contents of the message.
// The function then calls the GetMessageHash function with the necessary parameters to calculate the message hash.
// If an error occurs during the hashing process, an error is reported using the t.Errorf function.
// The expected hash value is compared with the actual hash value returned by the function.
// If the values do not match, an error is reported using the t.Errorf function.
//
// Parameters:
// - t: a testing.T object that provides methods for testing functions
// Returns:
// - None
func TestGeneral_GetMessageHash(t *testing.T) {
	ttd, err := MockTypedData()
	require.NoError(t, err)

	mail := Mail{
		From: Person{
			Name:   "Cow",
			Wallet: "0xCD2a3d9F938E13CD947Ec05AbC7FE734Df8DD826",
		},
		To: Person{
			Name:   "Bob",
			Wallet: "0xbBbBBBBbbBBBbbbBbbBbbbbBBbBbbbbBbBbbBBbB",
		},
		Contents: "Hello, Bob!",
	}

	hash := ttd.GetMessageHash(utils.HexToBN("0xCD2a3d9F938E13CD947Ec05AbC7FE734Df8DD826"), mail)

	exp := "0x6fcff244f63e38b9d88b9e3378d44757710d1b244282b435cb472053c8d78d0"
	require.Equal(t, exp, utils.BigToHex(hash))
}

// BenchmarkGetMessageHash is a benchmark function for testing the GetMessageHash function.
//
// It tests the performance of the GetMessageHash function by running it with different input sizes.
// The input size is determined by the bit length of the address parameter, which is converted from
// a hexadecimal string to a big integer using the HexToBN function from the utils package.
//
// Parameters:
// - b: a testing.B object that provides methods for benchmarking the function
// Returns:
//
//	none
func BenchmarkGetMessageHash(b *testing.B) {
	ttd, err := MockTypedData()
	require.NoError(b, err)

	mail := Mail{
		From: Person{
			Name:   "Cow",
			Wallet: "0xCD2a3d9F938E13CD947Ec05AbC7FE734Df8DD826",
		},
		To: Person{
			Name:   "Bob",
			Wallet: "0xbBbBBBBbbBBBbbbBbbBbbbbBBbBbbbbBbBbbBBbB",
		},
		Contents: "Hello, Bob!",
	}
	addr := utils.HexToBN("0xCD2a3d9F938E13CD947Ec05AbC7FE734Df8DD826")
	b.Run(fmt.Sprintf("input_size_%d", addr.BitLen()), func(b *testing.B) {
		result := ttd.GetMessageHash(addr, mail)
		require.NotEmpty(b, result)
	})
}

// TestGeneral_GetDomainHash tests the GetDomainHash function.
// It creates a mock TypedData object and generates the hash of a typed message using the Starknet domain and curve.
// If there is an error during the hashing process, it logs the error.
// It then compares the generated hash with the expected hash and logs an error if they do not match.
//
// Parameters:
// - t: a testing.T object that provides methods for testing functions
// Returns:
//
//	none
func TestGeneral_GetDomainHash(t *testing.T) {
	ttd, err := MockTypedData()
	require.NoError(t, err)

	hash := ttd.GetTypedMessageHash("StarkNetDomain", ttd.Domain)

	exp := "0x54833b121883a3e3aebff48ec08a962f5742e5f7b973469c1f8f4f55d470b07"
	require.Equal(t, exp, utils.BigToHex(hash))
}

// TestGeneral_GetTypedMessageHash is a unit test for the GetTypedMessageHash function
// equivalent of get struct hash.
//
// It tests the generation of a typed message hash for a given mail object using a specific curve.
// The function expects the mail object to have a "From" field of type Person, a "To" field of type Person,
// and a "Contents" field of type string. It returns the generated hash as a byte array and an error object.
//
// Parameters:
// - t: a testing.T object that provides methods for testing functions
// Returns:
//
//	none
func TestGeneral_GetTypedMessageHash(t *testing.T) {
	ttd, err := MockTypedData()
	require.NoError(t, err)

	mail := Mail{
		From: Person{
			Name:   "Cow",
			Wallet: "0xCD2a3d9F938E13CD947Ec05AbC7FE734Df8DD826",
		},
		To: Person{
			Name:   "Bob",
			Wallet: "0xbBbBBBBbbBBBbbbBbbBbbbbBBbBbbbbBbBbbBBbB",
		},
		Contents: "Hello, Bob!",
	}

	hash := ttd.GetTypedMessageHash("Mail", mail)

	exp := "0x4758f1ed5e7503120c228cbcaba626f61514559e9ef5ed653b0b885e0f38aec"
	require.Equal(t, exp, utils.BigToHex(hash))
}

// TestGeneral_GetTypeHash tests the GetTypeHash function.
//
// It tests the GetTypeHash function by calling it with different input values
// and comparing the result with expected values. It also checks that the
// encoding of the types matches the expected values.
//
// Parameters:
// - t: The testing.T object used for reporting test failures and logging test output
// Returns:
//
//	none
func TestGeneral_GetTypeHash(t *testing.T) {
	require := require.New(t)

	ttd, err := MockTypedData()
	require.NoError(err)

	hash, err := ttd.GetTypeHash("StarkNetDomain")
	require.NoError(err)

	exp := "0x1bfc207425a47a5dfa1a50a4f5241203f50624ca5fdf5e18755765416b8e288"
	require.Equal(exp, utils.BigToHex(hash))

	enc := ttd.Types["StarkNetDomain"]
	require.Equal(exp, utils.BigToHex(enc.Encoding))

	pHash, err := ttd.GetTypeHash("Person")
	require.NoError(err)

	exp = "0x2896dbe4b96a67110f454c01e5336edc5bbc3635537efd690f122f4809cc855"
	require.Equal(exp, utils.BigToHex(pHash))

	enc = ttd.Types["Person"]
	require.Equal(exp, utils.BigToHex(enc.Encoding))
}

// TestGeneral_GetSelectorFromName tests the GetSelectorFromName function.
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
func TestGeneral_GetSelectorFromName(t *testing.T) {
	sel1 := utils.BigToHex(utils.GetSelectorFromName("initialize"))
	sel2 := utils.BigToHex(utils.GetSelectorFromName("mint"))
	sel3 := utils.BigToHex(utils.GetSelectorFromName("test"))

	exp1 := "0x79dc0da7c54b95f10aa182ad0a46400db63156920adb65eca2654c0945a463"
	exp2 := "0x2f0b3c5710379609eb5495f1ecd348cb28167711b73609fe565a72734550354"
	exp3 := "0x22ff5f21f0b81b113e63f7db6da94fedef11b2119b4088b89664fb9a3cb658"

	if sel1 != exp1 || sel2 != exp2 || sel3 != exp3 {
		t.Errorf("invalid Keccak256 encoding: %v %v %v\n", sel1, sel2, sel3)
	}
}

// TestGeneral_EncodeType tests the EncodeType function.
//
// It creates a mock typed data and calls the EncodeType method with the
// parameter "Mail". It checks if the returned encoding matches the expected
// encoding. If there is an error during the encoding process, it fails the
// test.
//
// Parameters:
// - t: The testing.T object used for reporting test failures and logging test output
// Returns:
//
//	none
func TestGeneral_EncodeType(t *testing.T) {
	ttd, err := MockTypedData()
	require.NoError(t, err)

	enc, err := ttd.EncodeType("Mail")
	require.NoError(t, err)

	exp := "Mail(from:Person,to:Person,contents:felt)Person(name:felt,wallet:felt)"
	require.Equal(t, exp, enc)
}
