package xsessions

import (
	ctypes "github.com/NethermindEth/starknet.go/types"
)

var (
	SESSION_TYPE_HASH         = ctypes.HexToBN("0x1aa0e1c56b45cf06a54534fa1707c54e520b842feb21d03b7deddb6f1e340c")
	STARKNET_MESSAGE          = ctypes.UTF8StrToBig("StarkNet Message")
	STARKNET_DOMAIN_TYPE_HASH = ctypes.HexToBN("0x13cda234a04d66db62c06b8e3ad5f91bd0c67286c2c7519a826cf49da6ba478")
	POLICY_TYPE_HASH          = ctypes.HexToBN("0x2f0026e78543f036f33e26a8f5891b88c58dc1e20cbbfaf0bb53274da6fa568")
)
