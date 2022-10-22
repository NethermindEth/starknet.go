package main

import (
	_ "embed"
	"log"
	"os"
)

func main() {
	c, err := parse(os.Args[1:])
	if err != nil {
		log.Fatalf("could not run the command...")
	}
	switch c.provider {
	// case PROVIDER_RPCV01:
	// 	switch c.command {
	// 	case "install":
	// 		c.installAccountWithRPCv01()
	// 	case "execute":
	// 		if c.withPlugin {
	// 			c.incrementWithSessionKey()
	// 			return
	// 		}
	// 		c.incrementWithRPCv01()
	// 	case "sum":
	// 		log.Fatalf("rpcv01 not yet implemented")
	// 	default:
	// 		log.Fatalf("unknown command: %s\n", c.command)
	// 	}
	// case PROVIDER_GATEWAY:
	// 	switch c.command {
	// 	case "install":
	// 		c.installAccountWithGateway()
	// 	case "execute":
	// 		if c.withPlugin {
	// 			c.incrementWithSessionKey()
	// 			return
	// 		}
	// 		c.incrementWithGateway()
	// 	case "sum":
	// 		c.sumWithGateway()
	// 	default:
	// 		log.Fatalf("unknown command: %s\n", c.command)
	// 	}
	default:
		log.Fatal("provider provider only supports rpcv01 and gateway")
	}
}
