// Internal package to set up basic configurations for the examples contained in the 'examples' folder.
package setup

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/joho/godotenv"
)

// Loads environment variables contained in the ".env" file in the root of "examples" folder.
func init() {
	err := godotenv.Load("../.env")
	if err != nil {
		panic(errors.Join(errors.New("error loading '.env' file"), err))
	}
}

// Default "panic" but printing all RPCError fields (code, message, and data)
func PanicRPC(err error) {

	RPCErr, ok := err.(*rpc.RPCError)
	if !ok {
		panic("failed to cast to RPCError. This error is not a RPCError")
	}
	err = errors.Join(
		errors.New(fmt.Sprint(RPCErr.Code)),
		errors.New(RPCErr.Message),
		errors.New(fmt.Sprint(RPCErr.Data)),
	)
	panic(err)
}

// Validates whether the RPC_PROVIDER_URL variable has been set in the '.env' file and returns it; panics otherwise.
func GetRpcProviderUrl() string {
	return getEnv("RPC_PROVIDER_URL")
}

// Validates whether the PRIVATE_KEY variable has been set in the '.env' file and returns it; panics otherwise.
func GetPrivateKey() string {
	return getEnv("PRIVATE_KEY")
}

// Validates whether the PUBLIC_KEY variable has been set in the '.env' file and returns it; panics otherwise.
func GetPublicKey() string {
	return getEnv("PUBLIC_KEY")
}

// Validates whether the ACCOUNT_ADDRESS variable has been set in the '.env' file and returns it; panics otherwise.
func GetAccountAddress() string {
	return getEnv("ACCOUNT_ADDRESS")
}

// Validates whether the ACCOUNT_CAIRO_VERSION variable has been set in the '.env' file and returns it; panics otherwise.
func GetAccountCairoVersion() int {
	num, err := strconv.Atoi(getEnv("ACCOUNT_CAIRO_VERSION"))
	if err != nil {
		panic("Invalid ACCOUNT_CAIRO_VERSION number set in the '.env' file")
	}

	return num
}

// Validates whether the CONTRACT_ADDRESS variable has been set in the '.env' file and returns it; panics otherwise.
func GetContractAddress() string {
	return getEnv("CONTRACT_ADDRESS")
}

// Validates whether the FROM_BLOCK and TO_BLOCK variables have been set in the '.env' file and returns them; panics otherwise.
func GetFromAndToBlocks() (uint64, uint64) {
	fromBlock, err := strconv.ParseUint(getEnv("FROM_BLOCK"), 10, 64)
	if err != nil {
	    msg := fmt.Errorf("error parsing FROM_BLOCK from the .env file: %w", err)
		panic(msg)
	}
	toBlock, err := strconv.ParseUint(getEnv("TO_BLOCK"), 10, 64)
	if err != nil {
	    msg := fmt.Errorf("error parsing TO_BLOCK from the .env file: %w", err)
		panic(msg)
	}
	return fromBlock, toBlock
}

// Loads an env variable by name and returns it; panics otherwise.
func getEnv(envName string) string {
	env := os.Getenv(envName)
	if env == "" {
		panic(fmt.Sprintf("%s variable not set in the '.env' file", envName))
	}
	return env
}
