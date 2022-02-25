package caigo

import (
	"fmt"
	"math/big"
	"strings"
	"testing"
)

func TestCallContract(t *testing.T) {
	gw := NewGateway()

	req := StarknetRequest{
		ContractAddress:    "0x201d2543c023f51efc7043390ff1b24f855a6878b38aae05fd8c9ec440c3b44",
		EntryPointSelector: BigToHex(GetSelectorFromName("get_signer")),
		Calldata:           []string{},
		Signature:          []string{},
	}

	resp, err := gw.Call(req)
	if err != nil || resp[0] != "0x4e52f2f40700e9cdd0f386c31a1f160d0f310504fc508a1051b747a26070d10" {
		t.Errorf("Could not get call get_signer: %v\n", err)
	}

	resp, err = gw.Call(req, "0x75f625980778cd3f3dd5e3f2337e1872a89b3e83cf935c9975e20e1645ed295")
	if !strings.Contains(err.Error(), "UNINITIALIZED_CONTRACT") {
		t.Errorf("Could not get call get_signer: %v\n", err)
	}

	resp, err = gw.Call(req, "28062")
	if !strings.Contains(err.Error(), "UNINITIALIZED_CONTRACT") {
		t.Errorf("Could not get call get_signer: %v\n", err)
	}
}

// {
// 	"type": "INVOKE_FUNCTION",
// 	"contract_address": "0x0077b19d49e6069372d53e535fc9f3230a99b85ad46cc0934491bb6fb59a5a29",
// 	"entry_point_selector": "0x3e05da5242cf902163355db35f315cab809e5fed8079e4b39e9c75f277674a3",
// 	"nonce": "22",
// 	"calldata": [
// 	  "0",
// 	  "109356099333997",
// 	  "0x217d176acd37d6d456c433dd5246af96afc03d9f4d9241e815917ad81d639a1",
// 	  "0x217d176acd37d6d456c433dd5246af96afc03d9f4d9241e815917ad81d639a1",
// 	  "1500",
// 	  "0",
// 	  "109356099333997",
// 	  "1953066617",
// 	  "2342",
// 	  "132412"
// 	]
//   }
// {
// 	"type": "INVOKE_FUNCTION",
// 	"contract_address": "0x00f4c5e82ddb6894411d6ae48b33284ed1cc6b167551e0a76f27811700b1c3c2",
// 	"entry_point_selector": "0x7a44dde9fea32737a5cf3f9683b3235138654aa2d189f6fe44af37a61dc60d"
//   }
func TestInvokeContract(t *testing.T) {
	curve, err := SCWithConstants("./pedersen_params.json")
	if err != nil {
		t.Errorf("Could not init with constant points: %v\n", err)
	}

	priv := HexToBN("0x2fcf785f63f75236df5cd99e9fc61eb85923230f085e004945aa702a7f53a82")
	x, _, _ := curve.PrivateToPoint(priv)

	gw := NewGateway()

	tx := Transaction{
		ContractAddress:    HexToBN("0x00f4c5e82ddb6894411d6ae48b33284ed1cc6b167551e0a76f27811700b1c3c2"),
		EntryPointSelector: GetSelectorFromName("increment"),
		Calldata:           []*big.Int{},
		Signature:          []*big.Int{},
		Nonce:              big.NewInt(4),
	}

	hash, err := curve.HashTx(INVOKE, gw.ChainId, tx)

	pk := HexToBN("0x019800ea6a9a73f94aee6a3d2edf018fc770443e90c7ba121e8303ec6b349279")
	x, _, _ = curve.PrivateToPoint(pk)

	hash, err = curve.HashMsg(x, tx)

	r, s, err := curve.Sign(hash, pk)
	if err != nil {
		t.Errorf("Could not convert gen signature: %v\n", err)
	}

	req := StarknetRequest{
		ContractAddress:    "0x00f4c5e82ddb6894411d6ae48b33284ed1cc6b167551e0a76f27811700b1c3c2",
		EntryPointSelector: BigToHex(GetSelectorFromName("increment")),
		Calldata:           []string{},
		Signature:          []string{r.String(), s.String()},
		Type:               "INVOKE_FUNCTION",
	}

	resp, err := gw.Invoke(req)
	fmt.Println("RESP ERR: ", resp, err)
}

func TestGetters(t *testing.T) {
	gw := NewGateway()

	ret, err := gw.GetBlockHashById("3")
	if err != nil || ret != "0x2e65d0ff5b424d5fe9c71d5a1c3263e622234dc3bc4f4595090ee2c54205670" {
		t.Errorf("Could not get block hash by id: %v\n", err)
	}

	ret, err = gw.GetBlockIdByHash("0x60113ac2e217700f13406c6b7429331105484872e4cfa0ed3ffcf08f4c14f95")
	if err != nil || ret != "64499" {
		t.Errorf("Could not get block id by hash: %v\n", err)
	}

	ret, err = gw.GetTransactionHashById("3")
	if err != nil || ret != "0x1822471b7751cbaf98a5cce0003181af95d588e38c958739213af59f389fdc5" {
		t.Errorf("Could not get transaction hash by id: %v\n", err)
	}

	ret, err = gw.GetTransactionIdByHash("0x1822471b7751cbaf98a5cce0003181af95d588e38c958739213af59f389fdc5")
	if err != nil || ret != "3" {
		t.Errorf("Could not get transaction id by hash: %v\n", err)
	}

	ret, err = gw.GetStorageAt("0x01d1f307c073bb786a66e6e042ec2a9bdc385a3373bb3738d95b966d5ce56166", "0", "36663")
	if err != nil || ret != "0x0" {
		t.Errorf("Could not get storage: %v\n", err)
	}

	ret, err = gw.GetStorageAt("0x01d1f307c073bb786a66e6e042ec2a9bdc385a3373bb3738d95b966d5ce56166", "0", "")
	if err != nil || ret != "0x0" {
		t.Errorf("Could not get storage w/o blockId: %v\n", err)
	}

	code, err := gw.GetCode("0x057f67ac7904bfa10fa7870d5d1776d694e912cfcc9eff2dfa09938e2fa8d05d", "")
	if err != nil || code.Bytecode[0] != "0x40780017fff7fff" || len(code.Abi) == 0 {
		t.Errorf("Could not get code: %v\n", err)
	}

	block, err := gw.GetBlock("0x75b944d03a204b13c6f40a6ef842a69721c1343a8b381cf9e7c12759b4ffb75")
	if err != nil || block.BlockNumber != 64454 {
		t.Errorf("Could not get block by hash: %v\n", err)
	}

	block, err = gw.GetBlock("64454")
	if err != nil || block.BlockHash != "0x75b944d03a204b13c6f40a6ef842a69721c1343a8b381cf9e7c12759b4ffb75" {
		t.Errorf("Could not get block by num: %v\n", err)
	}

	status, err := gw.GetTransactionStatus("0x28f36ac0f14f21cf5d562ddbc0d1d875a103a1b4e44f640ccf0a88299c3aa33")
	if err != nil || status.TxStatus != "ACCEPTED_ON_L1" {
		t.Errorf("Could not get tx status: %v\n", err)
	}

	tx, err := gw.GetTransaction("0x28f36ac0f14f21cf5d562ddbc0d1d875a103a1b4e44f640ccf0a88299c3aa33")
	if err != nil || tx.BlockNumber != 28062 {
		t.Errorf("Could not get tx: %v\n", err)
	}

	receipt, err := gw.GetTransactionReceipt("0x28f36ac0f14f21cf5d562ddbc0d1d875a103a1b4e44f640ccf0a88299c3aa33")
	if err != nil || receipt.BlockHash != "0x75f625980778cd3f3dd5e3f2337e1872a89b3e83cf935c9975e20e1645ed295" {
		t.Errorf("Could not get tx receipt: %v\n", err)
	}
}
