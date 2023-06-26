package gateway_test

// func TestWithErrorHandler(t *testing.T) {

// 	type testSetType struct{}
// 	testSet := map[string][]testSetType{
// 		"devnet":  {},
// 		"testnet": {},
// 		"mainnet": {},
// 		"mock":    {},
// 	}[testEnv]

// 	for range testSet {
// 		spy := false
// 		spyPtr := &spy

// 		opt := gateway.WithErrorHandler(func(err error) error {
// 			*spyPtr = true
// 			return err
// 		})

// 		client := gateway.NewClient(opt)
// 		var result string
// 		request, _ := http.NewRequest(http.MethodGet, "https://httpstat.us/400", nil)
// 		err := client.do(request, &result)

// 		if err == nil {
// 			t.Errorf("should fail due to http-400")
// 		}
// 		if !*spyPtr {
// 			t.Errorf("spy should be true, got %t", *spyPtr)
// 		}
// 	}
// }
