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
		"example_array",
		"example_baseTypes",
		"example_enum",
		"example_presetTypes",
		"mail_StructArray",
		"session_MerkleTree",
		"v1Nested",
		"allInOne",
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
	type testSetType struct {
		TypedData           TypedData
		Address             string
		ExpectedMessageHash string
	}
	testSet := []testSetType{
		{
			TypedData:           typedDataExamples["baseExample"],
			Address:             "0xCD2a3d9F938E13CD947Ec05AbC7FE734Df8DD826",
			ExpectedMessageHash: "0x6fcff244f63e38b9d88b9e3378d44757710d1b244282b435cb472053c8d78d0",
		},
		{
			TypedData:           typedDataExamples["example_array"],
			Address:             "0xCD2a3d9F938E13CD947Ec05AbC7FE734Df8DD826",
			ExpectedMessageHash: "0x88edea26d6177a8bc545b2e73c960ab7ddd67b46237b386b514e50315ce0f4",
		},
		{
			TypedData:           typedDataExamples["example_baseTypes"],
			Address:             "0xCD2a3d9F938E13CD947Ec05AbC7FE734Df8DD826",
			ExpectedMessageHash: "0xdb7829db8909c0c5496f5952bcfc4fc894341ce01842537fc4f448743480b6",
		},
		{
			TypedData:           typedDataExamples["example_presetTypes"],
			Address:             "0xCD2a3d9F938E13CD947Ec05AbC7FE734Df8DD826",
			ExpectedMessageHash: "0x185b339d5c566a883561a88fb36da301051e2c0225deb325c91bb7aa2f3473a",
		},
		{
			TypedData:           typedDataExamples["session_MerkleTree"],
			Address:             "0xCD2a3d9F938E13CD947Ec05AbC7FE734Df8DD826",
			ExpectedMessageHash: "0x751fb7d98545f7649d0d0eadc80d770fcd88d8cfaa55590b284f4e1b701ef0a",
		},
		{
			TypedData:           typedDataExamples["mail_StructArray"],
			Address:             "0xCD2a3d9F938E13CD947Ec05AbC7FE734Df8DD826",
			ExpectedMessageHash: "0x5914ed2764eca2e6a41eb037feefd3d2e33d9af6225a9e7fe31ac943ff712c",
		},
		{
			TypedData:           typedDataExamples["v1Nested"],
			Address:             "0xCD2a3d9F938E13CD947Ec05AbC7FE734Df8DD826",
			ExpectedMessageHash: "0x69b57cf0cd7c151c51f9616cc58a1f0a877fec28c8c15ff7537cf777c54a30d",
		},
		{
			TypedData:           typedDataExamples["example_enum"],
			Address:             "0xCD2a3d9F938E13CD947Ec05AbC7FE734Df8DD826",
			ExpectedMessageHash: "0x6e61abaf480b1370bbf231f54e298c5f4872f40a6d2dd409ff30accee5bbd1e",
		},
		{
			TypedData:           typedDataExamples["allInOne"],
			Address:             "0xCD2a3d9F938E13CD947Ec05AbC7FE734Df8DD826",
			ExpectedMessageHash: "0x8fa4e453de78c2762493760efd449a38eb46f85b2e02b116b77b3daa9075c8",
		},
	}

	for _, test := range testSet {
		hash, err := test.TypedData.GetMessageHash(test.Address)
		require.NoError(t, err)

		require.Equal(t, test.ExpectedMessageHash, hash.String())
	}
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
	type testSetType struct {
		TypedData    TypedData
		TypeName     string
		ExpectedHash string
	}
	testSet := []testSetType{
		{
			TypedData:    typedDataExamples["baseExample"],
			TypeName:     "StarkNetDomain",
			ExpectedHash: "0x1bfc207425a47a5dfa1a50a4f5241203f50624ca5fdf5e18755765416b8e288",
		},
		{
			TypedData:    typedDataExamples["baseExample"],
			TypeName:     "Mail",
			ExpectedHash: "0x13d89452df9512bf750f539ba3001b945576243288137ddb6c788457d4b2f79",
		},
		{
			TypedData:    typedDataExamples["example_baseTypes"],
			TypeName:     "Example",
			ExpectedHash: "0x1f94cd0be8b4097a41486170fdf09a4cd23aefbc74bb2344718562994c2c111",
		},
		{
			TypedData:    typedDataExamples["example_presetTypes"],
			TypeName:     "Example",
			ExpectedHash: "0x1a25a8bb84b761090b1fadaebe762c4b679b0d8883d2bedda695ea340839a55",
		},
		{
			TypedData:    typedDataExamples["session_MerkleTree"],
			TypeName:     "Session",
			ExpectedHash: "0x1aa0e1c56b45cf06a54534fa1707c54e520b842feb21d03b7deddb6f1e340c",
		},
	}
	for _, test := range testSet {
		hash, err := test.TypedData.GetTypeHash(test.TypeName)
		require.NoError(t, err)

		require.Equal(t, test.ExpectedHash, hash.String())
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
	type testSetType struct {
		TypedData      TypedData
		TypeName       string
		ExpectedEncode string
	}
	testSet := []testSetType{
		{
			TypedData:      typedDataExamples["baseExample"],
			TypeName:       "StarkNetDomain",
			ExpectedEncode: "StarkNetDomain(name:felt,version:felt,chainId:felt)",
		},
		{
			TypedData:      typedDataExamples["baseExample"],
			TypeName:       "Mail",
			ExpectedEncode: "Mail(from:Person,to:Person,contents:felt)Person(name:felt,wallet:felt)",
		},
		{
			TypedData:      typedDataExamples["example_array"],
			TypeName:       "StarknetDomain",
			ExpectedEncode: `"StarknetDomain"("name":"shortstring","version":"shortstring","chainId":"shortstring","revision":"shortstring")`,
		},
		{
			TypedData:      typedDataExamples["example_baseTypes"],
			TypeName:       "Example",
			ExpectedEncode: `"Example"("n0":"felt","n1":"bool","n2":"string","n3":"selector","n4":"u128","n5":"i128","n6":"ContractAddress","n7":"ClassHash","n8":"timestamp","n9":"shortstring")`,
		},
		{
			TypedData:      typedDataExamples["example_presetTypes"],
			TypeName:       "Example",
			ExpectedEncode: `"Example"("n0":"TokenAmount","n1":"NftId")"NftId"("collection_address":"ContractAddress","token_id":"u256")"TokenAmount"("token_address":"ContractAddress","amount":"u256")"u256"("low":"u128","high":"u128")`,
		},
		{
			TypedData:      typedDataExamples["session_MerkleTree"],
			TypeName:       "Session",
			ExpectedEncode: `Session(key:felt,expires:felt,root:merkletree)`,
		},
		{
			TypedData:      typedDataExamples["mail_StructArray"],
			TypeName:       "Mail",
			ExpectedEncode: `Mail(from:Person,to:Person,posts_len:felt,posts:Post*)Person(name:felt,wallet:felt)Post(title:felt,content:felt)`,
		},
		{
			TypedData:      typedDataExamples["v1Nested"],
			TypeName:       "TransferERC721",
			ExpectedEncode: `"TransferERC721"("MessageId":"felt","From":"Account1","To":"Account1","Nft_to_transfer":"Nft","Comment1":"string","Comment2":"string","Comment3":"string")"Account1"("Name":"string","Address":"felt")"Nft"("Collection":"string","Address":"felt","Nft_id":"felt","Negotiated_for":"Transaction")"Transaction"("Qty":"string","Unit":"string","Token_address":"felt","Amount":"felt")`,
		},
		{
			TypedData:      typedDataExamples["example_enum"],
			TypeName:       "Example",
			ExpectedEncode: `"Example"("someEnum1":"EnumA","someEnum2":"EnumB")"EnumA"("Variant 1":(),"Variant 2":("u128","u128*"),"Variant 3":("u128"))"EnumB"("Variant 1":(),"Variant 2":("u128"))`,
		},
	}
	for _, test := range testSet {
		require.Equal(t, test.ExpectedEncode, test.TypedData.Types[test.TypeName].EncoddingString)
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
	type testSetType struct {
		TypedData    TypedData
		TypeName     string
		Context      []string
		ExpectedHash string
	}
	testSet := []testSetType{
		{
			TypedData:    typedDataExamples["baseExample"],
			TypeName:     "StarkNetDomain",
			ExpectedHash: "0x54833b121883a3e3aebff48ec08a962f5742e5f7b973469c1f8f4f55d470b07",
		},
		{
			TypedData:    typedDataExamples["example_baseTypes"],
			TypeName:     "Example",
			ExpectedHash: "0x75db031c1f5bf980cc48f46943b236cb85a95c8f3b3c8203572453075d3d39",
		},
		{
			TypedData:    typedDataExamples["example_presetTypes"],
			TypeName:     "Example",
			ExpectedHash: "0x74fba3f77f8a6111a9315bac313bf75ecfa46d1234e0fda60312fb6a6517667",
		},
		{
			TypedData:    typedDataExamples["session_MerkleTree"],
			TypeName:     "Session",
			ExpectedHash: "0x73602062421caf6ad2e942253debfad4584bff58930981364dcd378021defe8",
		},
		{
			TypedData:    typedDataExamples["v1Nested"],
			TypeName:     "TransferERC721",
			ExpectedHash: "0x11b5fb80dd88c3d8b6239b065def4ac9a79e6995b117ed5940a3a0734324b79",
		},
		{
			TypedData:    typedDataExamples["example_enum"],
			TypeName:     "Example",
			ExpectedHash: "0x1e1bb5d477e92cbf562b3b766c5c1e5f8590f2df868d4c8249c0db8416f8c37",
		},
	}
	for _, test := range testSet {
		hash, err := test.TypedData.GetStructHash(test.TypeName, test.Context...)
		require.NoError(t, err)

		require.Equal(t, test.ExpectedHash, hash.String())
	}
}
