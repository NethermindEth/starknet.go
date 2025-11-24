package main
 
import (
	"fmt"
	"math/big"
 
	"github.com/NethermindEth/starknet.go/merkle"
)
 
func main() {
	// Create two values to hash
	value1 := big.NewInt(100)
	value2 := big.NewInt(200)
 
	// Calculate the Merkle hash
	hash := merkle.MerkleHash(value1, value2)
 
	fmt.Printf("Hash of %d and %d: 0x%s\n",
		value1.Int64(),
		value2.Int64(),
		hash.Text(16))
}