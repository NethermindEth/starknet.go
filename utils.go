package caigo

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
)

func (sc StarkCurve) XToPubKey(x string) (*big.Int, *big.Int) {
	xin := HexToBN(x)

	yout := sc.GetYCoordinate(xin)

	return xin, yout
}

func StrToBig(str string) *big.Int {
	b, _ := new(big.Int).SetString(str, 10)

	return b
}

func HexToBN(hexString string) (n *big.Int) {
	numStr := strings.Replace(hexString, "0x", "", -1)

	n = new(big.Int)
	n.SetString(numStr, 16)
	return n
}

func HexToBytes(hexString string) ([]byte, error) {
	numStr := strings.Replace(hexString, "0x", "", -1)
	if (len(numStr) % 2) != 0 {
		numStr = fmt.Sprintf("%s%s", "0", numStr)
	}

	return hex.DecodeString(numStr)
}

func BigToHex(in *big.Int) string {
	return fmt.Sprintf("0x%x", in)
}
