package caigo

import (
	"context"
	"fmt"
	"math/big"
)

func (sg *StarknetGateway) AccountNonce(ctx context.Context, address string) (nonce *big.Int, err error) {
	resp, err := sg.Call(ctx, Transaction{
		ContractAddress:    address,
		EntryPointSelector: "get_nonce",
	})
	if err != nil {
		return nonce, err
	}
	if len(resp) == 0 {
		return nonce, fmt.Errorf("no resp in contract call 'get_nonce' %v", address)
	}

	return HexToBN(resp[0]), nil
}

func fmtBlockId(blockId string) string {
	if len(blockId) < 2 {
		return ""
	}

	if blockId[:2] == "0x" {
		return fmt.Sprintf("&blockHash=%s", blockId)
	}
	return fmt.Sprintf("&blockNumber=%s", blockId)
}
