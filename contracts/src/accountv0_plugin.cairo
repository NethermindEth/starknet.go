# SPDX-License-Identifier: MIT
# OpenZeppelin Contracts for Cairo v0.3.2 (account/presets/Account.cairo)

%lang starknet

from starkware.cairo.common.bool import TRUE, FALSE
from starkware.cairo.common.cairo_builtins import HashBuiltin, SignatureBuiltin, BitwiseBuiltin
from starkware.cairo.common.registers import get_fp_and_pc
from starkware.starknet.common.syscalls import get_tx_info

from openzeppelin.account.library import Account, AccountCallArray
from openzeppelin.upgrades.library import Proxy

from openzeppelin.introspection.erc165.library import ERC165

from plugin import PluginUtils, USE_PLUGIN

#
# Constructor
#

@constructor
func constructor{syscall_ptr : felt*, pedersen_ptr : HashBuiltin*, range_check_ptr}():
    return ()
end

#
# Getters
#

@view
func get_public_key{syscall_ptr : felt*, pedersen_ptr : HashBuiltin*, range_check_ptr}() -> (
    res : felt
):
    let (res) = Account.get_public_key()
    return (res=res)
end

@view
func get_nonce{syscall_ptr : felt*, pedersen_ptr : HashBuiltin*, range_check_ptr}() -> (res : felt):
    let (res) = Account.get_nonce()
    return (res=res)
end

@view
func supportsInterface{syscall_ptr : felt*, pedersen_ptr : HashBuiltin*, range_check_ptr}(
    interfaceId : felt
) -> (success : felt):
    let (success) = ERC165.supports_interface(interfaceId)
    return (success)
end

#
# Setters
#

@external
func set_public_key{syscall_ptr : felt*, pedersen_ptr : HashBuiltin*, range_check_ptr}(
    new_public_key : felt
):
    Account.set_public_key(new_public_key)
    return ()
end

@external
func initialize{syscall_ptr: felt*, pedersen_ptr: HashBuiltin*, range_check_ptr}(
    public_key : felt, plugin : felt
):
    let (current) = Account.get_public_key()
    with_attr error_message("account already initialized"):
        assert current = FALSE
    end
    Account.initializer(public_key)
    PluginUtils.initializer(plugin)
    return ()
end

@external
func upgrade{syscall_ptr: felt*, pedersen_ptr: HashBuiltin*, range_check_ptr}(
    implementation: felt
):
    Account.assert_only_self()
    # TODO: make sure the new implementation is compliant with account
    
    # Upgrade
    Proxy._set_implementation_hash(implementation)
    return ()
end

#
# Business logic
#

@view
func is_valid_signature{
    syscall_ptr : felt*, pedersen_ptr : HashBuiltin*, range_check_ptr, ecdsa_ptr : SignatureBuiltin*
}(hash : felt, signature_len : felt, signature : felt*) -> (is_valid : felt):
    let (is_valid) = Account.is_valid_signature(hash, signature_len, signature)
    return (is_valid=is_valid)
end

@external
func __execute__{
    syscall_ptr : felt*,
    pedersen_ptr : HashBuiltin*,
    range_check_ptr,
    ecdsa_ptr : SignatureBuiltin*,
    bitwise_ptr : BitwiseBuiltin*,
}(
    call_array_len : felt,
    call_array : AccountCallArray*,
    calldata_len : felt,
    calldata : felt*,
    nonce : felt,
) -> (response_len : felt, response : felt*):
    let (response_len, response) = execute(
        call_array_len, call_array, calldata_len, calldata, nonce
    )
    return (response_len=response_len, response=response)
end

# Helpers

func execute{
    syscall_ptr : felt*,
    pedersen_ptr : HashBuiltin*,
    range_check_ptr,
    ecdsa_ptr : SignatureBuiltin*,
    bitwise_ptr : BitwiseBuiltin*,
}(
    call_array_len : felt,
    call_array : AccountCallArray*,
    calldata_len : felt,
    calldata : felt*,
    nonce : felt,
) -> (response_len : felt, response : felt*):
    alloc_locals

    let (__fp__, _) = get_fp_and_pc()
    let (tx_info) = get_tx_info()

    if (call_array[0].to - tx_info.account_contract_address) + (call_array[0].selector - USE_PLUGIN) == 0:
        with_attr error_message("Account: invalid plugin verification"):
            PluginUtils.validate(call_array_len, call_array, calldata_len, calldata)
        end
        return Account._unsafe_execute(
            call_array_len - 1,
            call_array + AccountCallArray.SIZE,
            calldata_len - call_array[0].data_len,
            calldata + call_array[0].data_offset,
            nonce,
        )
    end

    # validate transaction
    with_attr error_message("Account: invalid signature"):
        let (is_valid) = is_valid_signature(
            tx_info.transaction_hash, tx_info.signature_len, tx_info.signature
        )
        assert is_valid = TRUE
    end

    return Account._unsafe_execute(call_array_len, call_array, calldata_len, calldata, nonce)
end
