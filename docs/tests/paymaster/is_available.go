package main

import (
	"context"
	"fmt"
)

func main() {
	fmt.Println("IsAvailable Demo:")
	fmt.Println("  Method: pm.IsAvailable(ctx)")
	fmt.Println()
	fmt.Println("Purpose:")
	fmt.Println("  Check if the paymaster service is running and accepting requests")
	fmt.Println()
	fmt.Println("Returns:")
	fmt.Println("  - bool: true if service is available")
	fmt.Println("  - error: error if check fails")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  available, err := pm.IsAvailable(ctx)")
	fmt.Println("  if available {")
	fmt.Println("      fmt.Println(\"Paymaster service is ready\")")
	fmt.Println("  }")
	
	_ = context.Background()
}
