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

	devNet := devnet.NewDevNet(devnetURL)

	// Check if DevNet is alive
	if devNet.IsAlive() {
		fmt.Println("DevNet Status: Running ✓")
	} else {
		fmt.Println("DevNet Status: Not Running ✗")
		log.Fatal("Please start DevNet with: starknet-devnet")
	}
}
