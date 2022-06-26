package rpc

// func TestIntegrationNodeGetBlockTransactionCountByHash(t *testing.T) {
// 	godotenv.Load()
// 	if os.Getenv("INTEGRATION") != "true" {
// 		t.Skip("Skipping integration test")
// 	}
// 	node, err := NewNode("node")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	result, err := node.GetBlockTransactionCountByHash(context.Background(), "0x60113ac2e217700f13406c6b7429331105484872e4cfa0ed3ffcf08f4c14f95")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	if result != 17 {
// 		t.Fatal("Transactions should be 17, current:", result)
// 	}
// }

// func TestIntegrationNodeGetBlockTransactionCountByNumber(t *testing.T) {
// 	godotenv.Load()
// 	if os.Getenv("INTEGRATION") != "true" {
// 		t.Skip("Skipping integration test")
// 	}
// 	node, err := NewNode("node")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	result, err := node.GetBlockTransactionCountByNumber(context.Background(), 230000)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	if result != 30 {
// 		t.Fatal("Transactions should be 30, current:", result)
// 	}
// }

// func TestIntegrationNodeGetTransactionByBlockHashAndIndex(t *testing.T) {
// 	godotenv.Load()
// 	if os.Getenv("INTEGRATION") != "true" {
// 		t.Skip("Skipping integration test")
// 	}
// 	node, err := NewNode("node")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	result, err := node.GetTransactionByBlockHashAndIndex(context.Background(), "0x115aa451e374dbfdeb6f8d4c70133a39c6bb7b2948a4a3f0c9d5dda30f94044", 0)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	if result.TxnHash != "0x705547f8f2f8fdfb10ed533d909f76482bb293c5a32648d476774516a0bebd0" {
// 		t.Fatal("Transactions error:", result.TxnHash)
// 	}
// }

// func TestIntegrationNodeGetTransactionByBlockNumberAndIndex(t *testing.T) {
// 	godotenv.Load()
// 	if os.Getenv("INTEGRATION") != "true" {
// 		t.Skip("Skipping integration test")
// 	}
// 	node, err := NewNode("node")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	result, err := node.GetTransactionByBlockNumberAndIndex(context.Background(), 220000, 1)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	if result.TxnHash != "0x25a909f5b88db224b3ec7e307d18e7739fd222f7e8467c57a996aff787cc5b3" {
// 		t.Fatal("Transactions error:", result.TxnHash)
// 	}
// }

// func TestIntegrationNodePendingTransactions(t *testing.T) {
// 	godotenv.Load()
// 	if os.Getenv("INTEGRATION") != "true" {
// 		t.Skip("Skipping integration test")
// 	}
// 	node, err := NewNode("node")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	_, err = node.PendingTransactions(context.Background())
// 	// TODO: this code is not yet enable in pathfinder v0.2.2-alpha.
// 	if err == nil || !strings.Contains(err.Error(), "Method not found") {
// 		t.Fatal(err)
// 	}
// }

// func TestIntegrationNodeGetBlockByHash(t *testing.T) {
// 	godotenv.Load()
// 	if os.Getenv("INTEGRATION") != "true" {
// 		t.Skip("Skipping integration test")
// 	}
// 	node, err := NewNode("node")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	result, err := node.GetBlockByHash(context.Background(), "0x115aa451e374dbfdeb6f8d4c70133a39c6bb7b2948a4a3f0c9d5dda30f94044")
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	if result.BlockHash != "0x115aa451e374dbfdeb6f8d4c70133a39c6bb7b2948a4a3f0c9d5dda30f94044" {
// 		t.Fatal("Transactions error:", result.BlockHash)
// 	}
// }

// func TestIntegrationNodeGetStorageAt(t *testing.T) {
// 	godotenv.Load()
// 	if os.Getenv("INTEGRATION") != "true" {
// 		t.Skip("Skipping integration test")
// 	}
// 	node, err := NewNode("node")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	result, err := node.GetStorageAt(context.Background(), "0x6fbd460228d843b7fbef670ff15607bf72e19fa94de21e29811ada167b4ca39", "0x0206F38F7E4F15E87567361213C28F235CCCDAA1D7FD34C9DB1DFE9489C6A091", "pending")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	if result != "0x1e240" {
// 		t.Fatal("GetStorageAt error:", result)
// 	}
// }

// func TestIntegrationNodeGetTransactionByHash(t *testing.T) {
// 	godotenv.Load()
// 	if os.Getenv("INTEGRATION") != "true" {
// 		t.Skip("Skipping integration test")
// 	}
// 	node, err := NewNode("node")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	result, err := node.GetTransactionByHash(context.Background(), "0x25a909f5b88db224b3ec7e307d18e7739fd222f7e8467c57a996aff787cc5b3")
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	if result.TxnHash != "0x25a909f5b88db224b3ec7e307d18e7739fd222f7e8467c57a996aff787cc5b3" {
// 		t.Fatal("Transactions error:", result.TxnHash)
// 	}
// }

// func TestIntegrationNodeGetTransactionReceipt(t *testing.T) {
// 	godotenv.Load()
// 	if os.Getenv("INTEGRATION") != "true" {
// 		t.Skip("Skipping integration test")
// 	}
// 	node, err := NewNode("node")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	result, err := node.GetTransactionReceipt(context.Background(), "0x25a909f5b88db224b3ec7e307d18e7739fd222f7e8467c57a996aff787cc5b3")
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	if result.TxnHash != "0x25a909f5b88db224b3ec7e307d18e7739fd222f7e8467c57a996aff787cc5b3" {
// 		t.Fatal("Transactions error:", result.TxnHash)
// 	}
// }

// // TODO: implement that test that should be failing for now
// func TestIntegrationNodeGetClass(t *testing.T) {
// 	godotenv.Load()
// 	if os.Getenv("INTEGRATION") != "true" {
// 		t.Skip("Skipping integration test")
// 	}
// 	node, err := NewNode("node")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	err = node.GetClass(context.Background(), "0x6fbd460228d843b7fbef670ff15607bf72e19fa94de21e29811ada167b4ca39")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// }

// func TestIntegrationNodeGetEvents(t *testing.T) {
// 	godotenv.Load()
// 	if os.Getenv("INTEGRATION") != "true" {
// 		t.Skip("Skipping integration test")
// 	}
// 	node, err := NewNode("node")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	from := int64(237756)
// 	to := int64(237758)
// 	address := "0x002b1e566ac40ec30f8491186c1cfcaebc0285422974fdd50336d21fd7a99b81"
// 	selector := caigo.GetSelectorFromName("Transfer")
// 	result, err := node.GetEvents(context.Background(), RpcApiEventFilter{
// 		FromBlock:   &from,
// 		ToBlock:     &to,
// 		FromAddress: &address,
// 		Keys:        []string{selector.Text(16)},
// 		PageSize:    20,
// 		PageNumber:  0})
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	if len(result.Events) != 1 {
// 		t.Fatal("It should return some data, instead", len(result.Events))
// 	}
// }

// func TestIntegrationNodeCall(t *testing.T) {
// 	godotenv.Load()
// 	if os.Getenv("INTEGRATION") != "true" {
// 		t.Skip("Skipping integration test")
// 	}
// 	node, err := NewNode("node")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	address := "0x2b1e566ac40ec30f8491186c1cfcaebc0285422974fdd50336d21fd7a99b81"
// 	tokenId := "0x41a52c5f46eb2e4f1b791623b513ce653ed8addc9060ca01c4868d389e5e31b5"
// 	selector := caigo.GetSelectorFromName("tokenURI")
// 	result, err := node.Call(context.Background(), RpcApiCallParameters{
// 		ContractAddress:    address,
// 		EntryPointSelector: "0x" + selector.Text(16),
// 		CallData: []string{
// 			"0x" + tokenId[34:66],
// 			"0x" + tokenId[2:34],
// 		},
// 	}, "latest")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	if len(result) != 6 || result[0] != "0x5" {
// 		t.Fatal("It should return some data, instead of", result[0])
// 	}
// }

// //go:embed test/openzeppelin.json
// var openzeppelin []byte

// func TestIntegrationNodeAddDeclareTransaction(t *testing.T) {
// 	godotenv.Load()
// 	if os.Getenv("INTEGRATION") != "true" {
// 		t.Skip("Skipping integration test")
// 	}
// 	node, err := NewNode("node")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	classFile := ClassFile{}
// 	err = json.Unmarshal(openzeppelin, &classFile)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	declare, err := node.AddDeclareTransaction(context.Background(), classFile.Program, classFile.EntryPointsByType)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	fmt.Printf("Class: %s\n", declare.ClassHash)
// 	fmt.Printf("Tx: %s\n", declare.TxHash)
// }

// func TestIntegrationNodeaddDeployOpenZeppelin(t *testing.T) {
// 	godotenv.Load()
// 	if os.Getenv("INTEGRATION") != "true" {
// 		t.Skip("Skipping integration test")
// 	}
// 	node, err := NewNode("node")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	classFile := ClassFile{}
// 	err = json.Unmarshal(openzeppelin, &classFile)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	callData := []string{"0x5ff5eff3bed10c5109c25ab3618323d74a436e7e0b66a512ca6dbab27f08a6"}
// 	deploy, err := node.AddDeployTransaction(context.Background(), classFile.Program, classFile.EntryPointsByType, callData)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	fmt.Printf("%s\n", deploy)
// }

// //go:embed test/balance.json
// var balance []byte

// func TestIntegrationNodeaddDeployBalance1(t *testing.T) {
// 	godotenv.Load()
// 	if os.Getenv("INTEGRATION") != "true" {
// 		t.Skip("Skipping integration test")
// 	}
// 	node, err := NewNode("node")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	classFile := ClassFile{}
// 	err = json.Unmarshal(balance, &classFile)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	callData := []string{}
// 	deploy, err := node.AddDeployTransaction(context.Background(), classFile.Program, classFile.EntryPointsByType, callData)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	fmt.Printf("%s\n", deploy)
// }

// //go:embed test/balance2.json
// var balance2 []byte

// func TestIntegrationNodeaddDeployBalance2(t *testing.T) {
// 	godotenv.Load()
// 	if os.Getenv("INTEGRATION") != "true" {
// 		t.Skip("Skipping integration test")
// 	}
// 	node, err := NewNode("node")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	classFile := ClassFile{}
// 	err = json.Unmarshal(balance2, &classFile)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	callData := []string{"0x10"}
// 	deploy, err := node.AddDeployTransaction(context.Background(), classFile.Program, classFile.EntryPointsByType, callData)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	fmt.Printf("%s\n", deploy)
// }

// func TestIntegrationNodeaddInvokeTransaction(t *testing.T) {
// 	godotenv.Load()
// 	if os.Getenv("INTEGRATION") != "true" {
// 		t.Skip("Skipping integration test")
// 	}
// 	node, err := NewNode("node")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	contractAddress := "0x6c49d194c895308b3b4267dd41326afb1a360cfc5d7670c499c13dc3b0fd8ed"
// 	methodName := "constructor"
// 	callData := []string{"0x5ff5eff3bed10c5109c25ab3618323d74a436e7e0b66a512ca6dbab27f08a6"}
// 	invoke, err := node.AddInvokeTransaction(context.Background(), contractAddress, methodName, callData)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	fmt.Printf("%v\n", invoke)
// }
