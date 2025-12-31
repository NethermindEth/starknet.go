package main
 
import (
	"fmt"
	"log"
 
	"github.com/NethermindEth/starknet.go/typeddata"
)
 
func main() {
	// Get Revision 1 (Poseidon hashing)
	_, err := typeddata.GetRevision(1)
	if err != nil {
		log.Fatalf("Invalid revision: %v", err)
	}
 
	fmt.Printf("Revision 1 loaded\n")
	fmt.Printf("Uses Poseidon hashing\n")
	fmt.Printf("Supports: u128, i128, ContractAddress, enums, u256 preset\n")
 
	// Get Revision 0 (Pedersen hashing)
	_, err = typeddata.GetRevision(0)
	if err != nil {
		log.Fatalf("Invalid revision: %v", err)
	}
 
	fmt.Printf("\nRevision 0 loaded\n")
	fmt.Printf("Uses Pedersen hashing\n")
	fmt.Printf("Supports: felt, bool, string, selector\n")
}