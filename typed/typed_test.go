package typed

import (
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

type Mail struct {
	From     Person `json:"from"`
	To       Person `json:"to"`
	Contents string `json:"contents"`
}

type Person struct {
	Name   string `json:"name"`
	Wallet string `json:"wallet"`
}

var types = []TypeDefinition{
	{
		Name: "StarkNetDomain",
		Parameters: []TypeParameter{
			{Name: "name", Type: "felt"},
			{Name: "version", Type: "felt"},
			{Name: "chainId", Type: "felt"},
		},
	},
	{
		Name: "Mail",
		Parameters: []TypeParameter{
			{Name: "from", Type: "Person"},
			{Name: "to", Type: "Person"},
			{Name: "contents", Type: "felt"},
		},
	},
	{
		Name: "Person",
		Parameters: []TypeParameter{
			{Name: "name", Type: "felt"},
			{Name: "wallet", Type: "felt"},
		},
	},
}

var dm = Domain{
	Name:    "StarkNet Mail",
	Version: "1",
	ChainId: "1",
}

var typedDataExamples = make(map[string]TypedData)

func TestMain(m *testing.M) {
	//TODO: implement v1 so we can use other examples
	fileNames := []string{
		"baseExample",
		// "example_baseTypes",
		// "example_enum",
		// "example_presetTypes",
		// "mail_StructArray",
		// "session_MerkleTree",
		// "v1Nested",
	}

	for _, fileName := range fileNames {
		var ttd TypedData
		content, err := os.ReadFile(fmt.Sprintf("./tests/%s.json", fileName))
		if err != nil {
			panic(fmt.Errorf("fail to read file: %w", err))
		}
		err = json.Unmarshal(content, &ttd)
		if err != nil {
			panic(fmt.Errorf("fail to unmarshal TypedData: %w", err))
		}

		typedDataExamples[fileName] = ttd
	}

	os.Exit(m.Run())
}

func BMockTypedData(b *testing.B) (ttd TypedData) {
	b.Helper()
	content, err := os.ReadFile("./tests/baseExample.json")
	require.NoError(b, err)

	err = json.Unmarshal(content, &ttd)
	require.NoError(b, err)

	return
}

// The TestUnmarshal function tests the ability to correctly unmarshal (deserialize) JSON content from
// a file into a Go TypedData struct. It starts by reading a json file. The JSON content is then unmarshaled
// into a TypedData struct using the json.Unmarshal function. After unmarshaling, the test checks if there were
// any errors during the unmarshaling process, and if an error is found, the test will fail.
//
// Parameters:
// - t: a testing.T object that provides methods for testing functions
// Returns:
// - None
func TestUnmarshal(t *testing.T) {
	content, err := os.ReadFile("./tests/baseExample.json")
	require.NoError(t, err)

	var typedData TypedData
	err = json.Unmarshal(content, &typedData)
	require.NoError(t, err)
}

func TestGeneral_CreateMessageWithTypes(t *testing.T) {
	t.Skip("TODO: need to implement encodeData method")
	// for testSetType 2
	type Example1 struct {
		N0 Felt            `json:"n0"`
		N1 Bool            `json:"n1"`
		N2 String          `json:"n2"`
		N3 Selector        `json:"n3"`
		N4 U128            `json:"n4"`
		N5 I128            `json:"n5"`
		N6 ContractAddress `json:"n6"`
		N7 ClassHash       `json:"n7"`
		N8 Timestamp       `json:"n8"`
		N9 Shortstring     `json:"n9"`
	}

	// for testSetType 3
	type Example2 struct {
		N0 TokenAmount `json:"n0"`
		N1 NftId       `json:"n1"`
	}

	hex1, ok := new(big.Int).SetString("0x3e8", 0)
	require.True(t, ok)
	hex2, ok := new(big.Int).SetString("0x0", 0)
	require.True(t, ok)

	type testSetType struct {
		MessageWithString string
		MessageWithTypes  any
	}
	testSet := []testSetType{
		{
			MessageWithString: `
			{
				"from": {
					"name": "Cow",
					"wallet": "0xCD2a3d9F938E13CD947Ec05AbC7FE734Df8DD826"
				},
				"to": {
					"name": "Bob",
					"wallet": "0xbBbBBBBbbBBBbbbBbbBbbbbBBbBbbbbBbBbbBBbB"
				},
				"contents": "Hello, Bob!"
			}`,
			MessageWithTypes: Mail{
				From: Person{
					Name:   "Cow",
					Wallet: "0xCD2a3d9F938E13CD947Ec05AbC7FE734Df8DD826",
				},
				To: Person{
					Name:   "Bob",
					Wallet: "0xbBbBBBBbbBBBbbbBbbBbbbbBBbBbbbbBbBbbBBbB",
				},
				Contents: "Hello, Bob!",
			},
		},
		{
			MessageWithString: `
			{
				"n0": "0x3e8",
				"n1": true,
				"n2": "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.",
				"n3": "transfer",
				"n4": 10,
				"n5": -10,
				"n6": "0x3e8",
				"n7": "0x3e8",
				"n8": 1000,
				"n9": "transfer"
			}`,
			MessageWithTypes: Example1{
				N0: "0x3e8",
				N1: true,
				N2: "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.",
				N3: "transfer",
				N4: big.NewInt(10),
				N5: big.NewInt(-10),
				N6: "0x3e8",
				N7: "0x3e8",
				N8: big.NewInt(1000),
				N9: "transfer",
			},
		},
		{
			MessageWithString: `
			{
				"n0": {
					"token_address": "0x049d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7",
					"amount": {
						"low": "0x3e8",
						"high": "0x0"
					}
				},
				"n1": {
					"collection_address": "0x049d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7",
					"token_id": {
						"low": "0x3e8",
						"high": "0x0"
					}
				}
			}`,
			MessageWithTypes: Example2{
				N0: TokenAmount{
					TokenAddress: "0x049d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7",
					Amount: U256{
						Low:  hex1,
						High: hex2,
					},
				},
				N1: NftId{
					CollectionAddress: "0x049d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7",
					TokenID: U256{
						Low:  hex1,
						High: hex2,
					},
				},
			},
		},
	}
	for _, test := range testSet {
		ttd1, err := NewTypedData(types, "Mail", dm, []byte(test.MessageWithString))
		require.NoError(t, err)

		bytes, err := json.Marshal(test.MessageWithTypes)
		require.NoError(t, err)

		ttd2, err := NewTypedData(types, "Mail", dm, bytes)
		require.NoError(t, err)

		require.EqualValues(t, ttd1, ttd2)

	}
}

// TestMessageHash tests the GetMessageHash function.
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
func TestGetMessageHash(t *testing.T) {
	ttd := typedDataExamples["baseExample"]

	hash, err := ttd.GetMessageHash("0xCD2a3d9F938E13CD947Ec05AbC7FE734Df8DD826")
	require.NoError(t, err)

	exp := "0x6fcff244f63e38b9d88b9e3378d44757710d1b244282b435cb472053c8d78d0"
	require.Equal(t, exp, hash.String())
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
	ttd := BMockTypedData(b)

	addr := "0xCD2a3d9F938E13CD947Ec05AbC7FE734Df8DD826"
	b.Run(fmt.Sprintf("input_size_%d", len(addr)), func(b *testing.B) {
		result, err := ttd.GetMessageHash(addr)
		require.NoError(b, err)
		require.NotEmpty(b, result)
	})
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
func TestGetTypeHash(t *testing.T) {
	require := require.New(t)
	ttd := typedDataExamples["baseExample"]

	type testSetType struct {
		TypeName     string
		ExpectedHash string
	}
	testSet := []testSetType{
		// revision 0
		{
			TypeName:     "StarkNetDomain",
			ExpectedHash: "0x1bfc207425a47a5dfa1a50a4f5241203f50624ca5fdf5e18755765416b8e288",
		},
		{
			TypeName:     "Person",
			ExpectedHash: "0x2896dbe4b96a67110f454c01e5336edc5bbc3635537efd690f122f4809cc855",
		},
		{
			TypeName:     "Mail",
			ExpectedHash: "0x13d89452df9512bf750f539ba3001b945576243288137ddb6c788457d4b2f79",
		},
	}
	for _, test := range testSet {
		hash, err := ttd.GetTypeHash(test.TypeName)
		require.NoError(err)

		require.Equal(test.ExpectedHash, hash.String())
	}
}

// TestEncodeType tests the EncodeType function.
//
// It creates a mock typed data and calls the EncodeType method with the
// type name. It checks if the returned encoding matches the expected
// encoding. If there is an error during the encoding process, it fails the
// test.
//
// Parameters:
// - t: The testing.T object used for reporting test failures and logging test output
// Returns:
//
//	none
func TestEncodeType(t *testing.T) {
	require := require.New(t)
	ttd := typedDataExamples["baseExample"]

	type testSetType struct {
		TypeName       string
		ExpectedEncode string
		Revision       revision
	}
	testSet := []testSetType{
		// revision 0
		{
			TypeName:       "StarkNetDomain",
			ExpectedEncode: "StarkNetDomain(name:felt,version:felt,chainId:felt)",
			Revision:       RevisionV0,
		},
		{
			TypeName:       "Mail",
			ExpectedEncode: "Mail(from:Person,to:Person,contents:felt)Person(name:felt,wallet:felt)",
			Revision:       RevisionV0,
		},
		// revision 1
		{
			TypeName:       "StarkNetDomain",
			ExpectedEncode: `"StarkNetDomain"("name":"felt","version":"felt","chainId":"felt")`,
			Revision:       RevisionV1,
		},
		{
			TypeName:       "Mail",
			ExpectedEncode: `"Mail"("from":"Person","to":"Person","contents":"felt")"Person"("name":"felt","wallet":"felt")`,
			Revision:       RevisionV1,
		},
	}
	for _, test := range testSet {
		encode, err := encodeType(test.TypeName, ttd.Types, test.Revision.Version())
		require.NoError(err)

		require.Equal(test.ExpectedEncode, encode)
	}
}

// TestGetStructHash tests the GetStructHash function.
//
// It creates a mock typed data and calls the GetStructHash method with the
// type name. It checks if the returned encoding matches the expected
// encoding. If there is an error during the encoding process, it fails the
// test.
//
// Parameters:
// - t: The testing.T object used for reporting test failures and logging test output
// Returns:
//
//	none
func TestGetStructHash(t *testing.T) {
	ttd := typedDataExamples["baseExample"]

	type testSetType struct {
		TypeName     string
		ExpectedHash string
	}
	testSet := []testSetType{
		{
			TypeName:     "StarkNetDomain",
			ExpectedHash: "0x54833b121883a3e3aebff48ec08a962f5742e5f7b973469c1f8f4f55d470b07",
		},
		{
			TypeName:     "Mail",
			ExpectedHash: "0x4758f1ed5e7503120c228cbcaba626f61514559e9ef5ed653b0b885e0f38aec",
		},
	}
	for _, test := range testSet {
		hash, err := ttd.GetStructHash(test.TypeName)
		require.NoError(t, err)

		require.Equal(t, test.ExpectedHash, hash.String())
	}
}
