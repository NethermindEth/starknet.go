# TypedData Test Files - Real Outputs

This directory contains test files that capture REAL outputs from the starknet.go typeddata package.

## Test Files with Real Outputs ✅

1. ✅ **get_message_hash.go** - GetMessageHash() method
   - Creates typed data for a Person message
   - Calculates message hash for signing
   - Output: `0x2eab1684598adbcfe30a1aec930ebccbbe66656b8f713739eb799316f4647ce`

2. ✅ **get_type_hash.go** - GetTypeHash() method
   - Calculates type hashes for Person and Mail types
   - Outputs:
     - Person type hash: `0x2896dbe4b96a67110f454c01e5336edc5bbc3635537efd690f122f4809cc855`
     - Mail type hash: `0x13d89452df9512bf750f539ba3001b945576243288137ddb6c788457d4b2f79`

3. ✅ **get_struct_hash.go** - GetStructHash() method
   - Calculates struct hashes for message data
   - Outputs:
     - Person struct hash: `0x7aefbe519d42308ecc89e6bb0a4cd3b86507eeff38ab3b2de4fcd1b1ca5c3e8`
     - Domain struct hash: `0x1df03453b9b5f32266af6a69382de17acf1eb7d5b43a65bd4334658891b8519`

4. ✅ **new_typed_data.go** - NewTypedData() function
   - Creates TypedData from type definitions
   - Output: Message hash `0x68adb267019e501f5e86385407bb0564d9d76cb059717946aa8abbc5dc13644`

## Total: 4 TypedData functions tested with REAL outputs ✅

All tests demonstrate SNIP-12 compliant typed data signing for Starknet.
