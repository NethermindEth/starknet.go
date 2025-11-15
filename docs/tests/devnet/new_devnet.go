package main

import (
	"fmt"
	"log"
	"os"

	"github.com/NethermindEth/starknet.go/devnet"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")

	// Get DevNet URL from environment variable
	devnetURL := os.Getenv("DEVNET_URL")
	if devnetURL == "" {
		devnetURL = "http://localhost:5050"
	}

	// Create DevNet instance
	devNet := devnet.NewDevNet(devnetURL)

	fmt.Printf("DevNet instance created for: %s\n", devnetURL)

	// Verify connection
	if devNet.IsAlive() {
		fmt.Println("✓ DevNet is running")
	} else {
		log.Fatal("✗ DevNet is not running")
	}
}
