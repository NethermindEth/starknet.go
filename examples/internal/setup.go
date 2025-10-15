// Internal package to set up basic configurations for the examples contained in the 'examples' folder.
package setup

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/account"
	"github.com/joho/godotenv"
)

// Loads environment variables contained in the ".env" file in the root of "examples" folder.
func init() {
	err := godotenv.Load("../.env")
	if err != nil {
		panic(errors.Join(errors.New("error loading '.env' file"), err))
	}
}

// Validates whether the RPC_PROVIDER_URL variable has been set in the '.env' file and returns it; panics otherwise.
func GetRPCProviderURL() string {
	return getEnv("RPC_PROVIDER_URL")
}

// Validates whether the WS_PROVIDER_URL variable has been set in the '.env' file and returns it; panics otherwise.
func GetWsProviderURL() string {
	return getEnv("WS_PROVIDER_URL")
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
func GetAccountCairoVersion() account.CairoVersion {
	num, err := strconv.Atoi(getEnv("ACCOUNT_CAIRO_VERSION"))
	if err != nil {
		panic("Invalid ACCOUNT_CAIRO_VERSION number set in the '.env' file")
	}

	switch num {
	case 0:
		return account.CairoV0
	case 2:
		return account.CairoV2
	default:
		panic("Invalid ACCOUNT_CAIRO_VERSION number set in the '.env' file")
	}
}

// Loads an env variable by name and returns it; panics otherwise.
func getEnv(envName string) string {
	env := os.Getenv(envName)
	if env == "" {
		panic(envName + " variable not set in the '.env' file")
	}

	return env
}

// PadZerosInFelt it's a helper function that pads zeros to the left of a hex felt value to make sure it is 64 characters long.
func PadZerosInFelt(hexFelt *felt.Felt) string {
	length := 66
	hexStr := hexFelt.String()

	// Check if the hex value is already of the desired length
	if len(hexStr) >= length {
		return hexStr
	}

	// Extract the hex value without the "0x" prefix
	hexValue := hexStr[2:]
	// Pad zeros after the "0x" prefix
	paddedHexValue := fmt.Sprintf("%0*s", length-2, hexValue)
	// Add back the "0x" prefix to the padded hex value
	paddedHexStr := "0x" + paddedHexValue

	return paddedHexStr
}
