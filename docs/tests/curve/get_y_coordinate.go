package main

import (
	"fmt"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/curve"
)

func main() {
	// Known X-coordinate
	x, _ := new(felt.Felt).SetString("0x5f1614a304f6a40630792a4f6562eee36805f52e94f777d6a78463ae541229b")

	// Derive Y-coordinate
	y := curve.GetYCoordinate(x)

	fmt.Println("GetYCoordinate:")
	fmt.Printf("  X-Coordinate: %s\n", x.String())
	fmt.Printf("  Y-Coordinate: %s\n", y.String())
}
