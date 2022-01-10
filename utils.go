package caigo

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
	"math/rand"
	"time"
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

func (sc StarkCurve) GetRandomPrivateKey() *big.Int {
	max := new(big.Int)
	max = max.Sub(sc.N, big.NewInt(1))

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	priv := new(big.Int)
	priv = priv.Rand(r, max)
	return priv
}

func (sc StarkCurve) PrivateToPoint(privKey *big.Int) (x, y *big.Int, err error) {
	if privKey.Cmp(big.NewInt(0)) != 1 || privKey.Cmp(sc.N) != -1 {
		return x, y, fmt.Errorf("private key not in curve range")
	}
	x, y = sc.EcMult(privKey, sc.EcGenX, sc.EcGenY)
	return x, y, nil
}
