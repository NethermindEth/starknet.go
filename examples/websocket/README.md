This example demonstrates how to subscribe to new block headers using WebSocket. It can be adapted to subscribe to other methods as well.

Steps:
1. Rename the ".env.template" file located at the root of the "examples" folder to ".env"
1. Uncomment, and assign your Sepolia WebSocket testnet endpoint to the `WS_PROVIDER_URL` variable in the ".env" file
1. Make sure you are in the "websocket" directory
1. Execute `go run main.go`

The calls outuputs will be returned at the end of the execution.