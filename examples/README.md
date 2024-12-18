### Welcome!
To successfully execute these examples you'll need to configure some environment variables. To do so:

1. Rename the ".env.template" file located at the root of this folder to ".env"
1. Uncomment, and assign your Sepolia testnet endpoint to the `RPC_PROVIDER_URL` variable in the ".env" file
1. Uncomment, and assign your account address to the `ACCOUNT_ADDRESS` variable in the ".env" file (make sure to have a few ETH in it)
1. Uncomment, and assign your starknet public key to the `PUBLIC_KEY` variable in the ".env" file
1. Uncomment, and assign your private key to the `PRIVATE_KEY` variable in the ".env" file

The variables used vary from example to example, you'll see the required ones on each `main.go` file, usually after a "// Load variables from '.env' file" comment.
To run an example:

1. Make sure you are in the chosen example directory
1. Execute `go run main.go` to run it

#### Some FAQ answered by these examples
1. How to deploy an account via `DEPLOY_TRANSACTION`?  
  R: See [deployAccount](./deployAccount/main.go)
1. How to estimate fees?  
  R: See [deployAccount](./deployAccount/main.go), line 89.
1. How to generate random public and private keys?  
  R: See [deployAccount](./deployAccount/main.go), line 46.
1. How to use my existing account importing my account address, and public and private keys?  
  R: See [deployContractUDC](./deployContractUDC/main.go), lines 54 and 64.
1. How to get my nonce?  
  R: See [deployContractUDC](./deployContractUDC/main.go), line 70.
1. How to deploy a smart contract using UDC?  
  R: See [deployContractUDC](./deployContractUDC/main.go).
1. How to send an invoke transaction?  
  R: See [simpleInvoke](./simpleInvoke/main.go).
1. How to make multiple function calls in the same transaction?
  R: See [simpleInvoke](./simpleInvoke/main.go), line 92.
1. How to get the transaction status?  
  R: See [simpleInvoke](./simpleInvoke/main.go), line 131.
1. How to deploy an ERC20 token?  
  R: See [deployContractUDC](./deployContractUDC/main.go).
1. How to get the balance of a ERC20 token?  
  R: See [simpleCall](./simpleCall/main.go).
1. How to make a function call?  
  R: See [simpleCall](./simpleCall/main.go).
1. How to sign and verify a typed data?  
  R: See [typedData](./typedData/main.go).
