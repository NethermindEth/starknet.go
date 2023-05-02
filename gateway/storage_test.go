package gateway_test

import (
	"context"
	"testing"

	"github.com/dontpanicdao/caigo/gateway"
	//	"github.com/google/go-cmp/cmp"
)

/*
func TestProvider_StorageAt(t *testing.T) {
	// TODO: Find a better test case
	want := "0x0"

	for _, tc := range []struct {
		address string
		key     uint64
		opts    *StorageAtOptions
	}{{
		address: "0x06eeefce63bc81620e375c7501cb7b5aecdf9fb99aa5ec25886b8b854c4293cb",
		key:     1,
		opts:    nil,
	}, {
		address: "0x06eeefce63bc81620e375c7501cb7b5aecdf9fb99aa5ec25886b8b854c4293cb",
		key:     1,
		opts:    &StorageAtOptions{BlockNumber: 582},
	}, {
		address: "0x06eeefce63bc81620e375c7501cb7b5aecdf9fb99aa5ec25886b8b854c4293cb",
		key:     1,
		opts:    &StorageAtOptions{BlockHash: "0x182d83f0ed972e97fa529be0088e20b5a7cb32e3bba0d164d668147713e79f9"},
	}} {
		ctx := context.Background()
		sg := NewClient(WithChain("main"))
		got, err := sg.StorageAt(ctx, tc.address, tc.key, tc.opts)
		if err != nil {
			t.Fatalf("getting storage at: %v", err)
		}

		if tc.opts != nil {
			if diff := cmp.Diff(want, got, nil); diff != "" {
				t.Errorf("Storage value diff mismatch (-want +got):\n%s", diff)
			}
		}
	}
}
*/

func TestStorageAt(t *testing.T) {
	testConfig := beforeEach(t)

	type testSetType struct {
		address string
		key     string
		value   string
		opts    *gateway.StorageAtOptions
	}

	testSet := map[string][]testSetType{
		"devnet": {},
		"testnet": {{address: "0x035401b96dc690eda2716068d3b03732d7c18af7c0327787660179108789d84f",
			key:   "475322019845212235330707245667153666023074534120350221048512561271566416926",
			value: "0x14a7a59e3e2d058d4c7c868e05907b2b49e324cc5b6af71182f008feb939e91",
			opts:  &gateway.StorageAtOptions{BlockNumber: 281263}}},
	}[testEnv]

	for _, test := range testSet {
		val, err := testConfig.client.StorageAt(context.Background(), test.address, test.key, test.opts)

		if err != nil {
			t.Fatal(err)
		}
		if val != test.value {
			t.Fatalf("expecting %s, instead: %s", "", val)
		}
	}
}
