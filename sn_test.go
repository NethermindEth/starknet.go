package caigo

import (
	"fmt"
	"math/big"
	// "strings"
	"testing"
)

// TEST SIGNER:
// PRIVATE:  3525226865423913736626831370712856874147697093596274695646823855285486342436
// PRIVATE:  0x7cb352bb7ce23d083dc150c10e89c2592dae372adda12de807ef131e275d124
// PUBLIC:  1782664861422806383676456243932506795376436825307043764029510638703215652083
// PUBLIC:  0x03f0f3daafa031ea153899d7911282c84906fd47f443675e3ec8d4f4d0902cf3

func TestEvoke(t *testing.T) {
	curve, err := SCWithConstants("./pedersen_params.json")
	if err != nil {
		t.Errorf("Could not init with constant points: %v\n", err)
	}

	priv, _ := new(big.Int).SetString("3525226865423913736626831370712856874147697093596274695646823855285486342436", 10)
	x, y, _ := curve.PrivateToPoint(priv)

	hashy, err := curve.PedersenHash([]*big.Int{big.NewInt(1), big.NewInt(2)})

	r, s, err := curve.Sign(hashy, priv, big.NewInt(1231))
	if err != nil {
		t.Errorf("Could not convert gen signature: %v\n", err)
	}

	fmt.Println("R S: ", r, s, BigToHex(x), BigToHex(hashy))
	fmt.Println("VALID: ", curve.Verify(hashy, r, s, x, y))

	gw := NewGateway()

	n, stat, err := gw.PollTx("0x51c0676db2996b683f793ed4e3a6bbac7a725a3ee73313620183befd37a34c1", ACCEPTED_ON_L2, 1, 100)
	fmt.Println("DIS: ", n, stat, err)
}

// func TestDeploy(t *testing.T) {
// 	curve := SC()

// 	pk, _ := curve.GetRandomPrivateKey()
// 	x, _, err := curve.PrivateToPoint(pk)
// 	if err != nil {
// 		t.Errorf("Could not convert random private key to point: %v\n", err)
// 	}

// 	dr := DeployRequest{
// 		Type:                "DEPLOY",
// 		ContractAddressSalt: BigToHex(x),
// 		ConstructorCalldata: []string{x.String(), "0"},
// 	}

// 	gw := NewGateway()
// 	resp, err := gw.Deploy("ArgentAccount.json", dr)
// 	if err != nil {
// 		t.Errorf("Could not deploy contract: %v\n", err)
// 	}

// 	_, stat, err = gw.PollTx(resp.TransactionHash, "ACCEPTED_ON_L2", 1, 100)
// 	if err != nil {
// 		t.Errorf("Could not get success: %v\n", err)
// 	}
// }

// func TestInvokeContract(t *testing.T) {
// 	curve, err := SCWithConstants("./pedersen_params.json")
// 	if err != nil {
// 		t.Errorf("Could not init with constant points: %v\n", err)
// 	}

// 	priv := HexToBN("0x2fcf785f63f75236df5cd99e9fc61eb85923230f085e004945aa702a7f53a82")
// 	x, _, _ := curve.PrivateToPoint(priv)

// 	gw := NewGateway()

// 	hash, err := curve.HashTx(INVOKE, gw.ChainId, tx)

// 	pk := HexToBN("0x019800ea6a9a73f94aee6a3d2edf018fc770443e90c7ba121e8303ec6b349279")
// 	x, _, _ = curve.PrivateToPoint(pk)

// 	hash, err = curve.HashMsg(x, tx)

// 	r, s, err := curve.Sign(hash, pk)
// 	if err != nil {
// 		t.Errorf("Could not convert gen signature: %v\n", err)
// 	}

// 	req := StarknetRequest{
// 		ContractAddress:    "0x00f4c5e82ddb6894411d6ae48b33284ed1cc6b167551e0a76f27811700b1c3c2",
// 		EntryPointSelector: BigToHex(GetSelectorFromName("increment")),
// 		Calldata:           []string{},
// 		Signature:          []string{r.String(), s.String()},
// 		Type:               "INVOKE_FUNCTION",
// 	}

// 	resp, err := gw.Invoke(req)
// 	fmt.Println("RESP ERR: ", resp, err)
// }

// func TestInvokeContract(t *testing.T) {
// 	gw := NewGateway()

// 	req := StarknetRequest{
// 		Type:               "INVOKE_FUNCTION",
// 		ContractAddress:    "0x0077b19d49e6069372d53e535fc9f3230a99b85ad46cc0934491bb6fb59a5a29",
// 		EntryPointSelector: BigToHex(GetSelectorFromName("update_l1_address")),
// 		Calldata:           []string{HexToBN("0xDEADBEEF").String()},
// 		Signature:          []string{},
// 	}

// 	resp, err := gw.Invoke(req)
// 	if err != nil && strings.Contains(code, "NOT") {
// 		t.Errorf("Could not add tx: %v\n", err)
// 	}
// }

// func TestInvokeWithSig(t *testing.T) {
// 	gw := NewGateway()

// 	tx := Transaction{
// 		ContractAddress:    HexToBN("0x0077b19d49e6069372d53e535fc9f3230a99b85ad46cc0934491bb6fb59a5a29"),
// 		EntryPointSelector: GetSelectorFromName("update_l1_address"),
// 		Calldata:           []*big.Int{HexToBN("0xDEADBEEF")},
// 		Nonce:              big.NewInt(25),
// 		Signature:          []*big.Int{},
// 	}

// 	req := StarknetRequest{
// 		Type:               "INVOKE_FUNCTION",
// 		ContractAddress:    "0x0077b19d49e6069372d53e535fc9f3230a99b85ad46cc0934491bb6fb59a5a29",
// 		EntryPointSelector: BigToHex(GetSelectorFromName("update_l1_address")),
// 		Calldata:           []string{HexToBN("0xDEADBEEF").String()},
// 		Signature:          []string{},
// 	}

// 	resp, err := gw.Invoke(req)
// 	if err != nil && strings.Contains(code, "NOT") {
// 		t.Errorf("Could not add tx: %v\n", err)
// 	}
// }

// func TestGetters(t *testing.T) {
// 	gw := NewGateway()

// 	ret, err := gw.GetBlockHashById("3")
// 	if err != nil || ret != "0x2e65d0ff5b424d5fe9c71d5a1c3263e622234dc3bc4f4595090ee2c54205670" {
// 		t.Errorf("Could not get block hash by id: %v\n", err)
// 	}

// 	ret, err = gw.GetBlockIdByHash("0x60113ac2e217700f13406c6b7429331105484872e4cfa0ed3ffcf08f4c14f95")
// 	if err != nil || ret != "64499" {
// 		t.Errorf("Could not get block id by hash: %v\n", err)
// 	}

// 	ret, err = gw.GetTransactionHashById("3")
// 	if err != nil || ret != "0x1822471b7751cbaf98a5cce0003181af95d588e38c958739213af59f389fdc5" {
// 		t.Errorf("Could not get transaction hash by id: %v\n", err)
// 	}

// 	ret, err = gw.GetTransactionIdByHash("0x1822471b7751cbaf98a5cce0003181af95d588e38c958739213af59f389fdc5")
// 	if err != nil || ret != "3" {
// 		t.Errorf("Could not get transaction id by hash: %v\n", err)
// 	}

// 	ret, err = gw.GetStorageAt("0x01d1f307c073bb786a66e6e042ec2a9bdc385a3373bb3738d95b966d5ce56166", "0", "36663")
// 	if err != nil || ret != "0x0" {
// 		t.Errorf("Could not get storage: %v\n", err)
// 	}

// 	ret, err = gw.GetStorageAt("0x01d1f307c073bb786a66e6e042ec2a9bdc385a3373bb3738d95b966d5ce56166", "0", "")
// 	if err != nil || ret != "0x0" {
// 		t.Errorf("Could not get storage w/o blockId: %v\n", err)
// 	}

// 	code, err := gw.GetCode("0x057f67ac7904bfa10fa7870d5d1776d694e912cfcc9eff2dfa09938e2fa8d05d", "")
// 	if err != nil || code.Bytecode[0] != "0x40780017fff7fff" || len(code.Abi) == 0 {
// 		t.Errorf("Could not get code: %v\n", err)
// 	}

// 	block, err := gw.GetBlock("0x75b944d03a204b13c6f40a6ef842a69721c1343a8b381cf9e7c12759b4ffb75")
// 	if err != nil || block.BlockNumber != 64454 {
// 		t.Errorf("Could not get block by hash: %v\n", err)
// 	}

// 	block, err = gw.GetBlock("64454")
// 	if err != nil || block.BlockHash != "0x75b944d03a204b13c6f40a6ef842a69721c1343a8b381cf9e7c12759b4ffb75" {
// 		t.Errorf("Could not get block by num: %v\n", err)
// 	}

// 	status, err := gw.GetTransactionStatus("0x28f36ac0f14f21cf5d562ddbc0d1d875a103a1b4e44f640ccf0a88299c3aa33")
// 	if err != nil || status.TxStatus != "ACCEPTED_ON_L1" {
// 		t.Errorf("Could not get tx status: %v\n", err)
// 	}

// 	tx, err := gw.GetTransaction("0x28f36ac0f14f21cf5d562ddbc0d1d875a103a1b4e44f640ccf0a88299c3aa33")
// 	if err != nil || tx.BlockNumber != 28062 {
// 		t.Errorf("Could not get tx: %v\n", err)
// 	}

// 	receipt, err := gw.GetTransactionReceipt("0x28f36ac0f14f21cf5d562ddbc0d1d875a103a1b4e44f640ccf0a88299c3aa33")
// 	if err != nil || receipt.BlockHash != "0x75f625980778cd3f3dd5e3f2337e1872a89b3e83cf935c9975e20e1645ed295" {
// 		t.Errorf("Could not get tx receipt: %v\n", err)
// 	}
// }
