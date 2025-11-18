package main

import (
	"fmt"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/curve"
)

func main() {
	// Sample x-coordinate on the curve
	xCoord, _ := new(felt.Felt).SetString("0x5f1614a304f6a40630792a4f6562eee36805f52e94f777d6a78463ae541229b")

	// Get the y-coordinate
	yCoord := curve.GetYCoordinate(xCoord)

	fmt.Println("GetYCoordinate:")
	fmt.Printf("  X-Coordinate: %s\n", xCoord.String())
	fmt.Printf("  Y-Coordinate: %s\n", yCoord.String())
}
