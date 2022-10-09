%lang starknet

from starkware.cairo.common.math import assert_not_zero
from starkware.cairo.common.cairo_builtins import HashBuiltin, SignatureBuiltin, BitwiseBuiltin
from starkware.starknet.common.syscalls import get_caller_address

from openzeppelin.account.library import Account, AccountCallArray

const USE_PLUGIN = 1121675007639292412441492001821602921366030142137563176027248191276862353634

@contract_interface
namespace IPlugin:
    # Method to call during validation
    func validate(
        plugin_data_len: felt,
        plugin_data: felt*,
        call_array_len: felt,
        call_array: AccountCallArray*,
        calldata_len: felt,
        calldata: felt*
    ):
    end
end

@event
func PluginRegistered(accountAddress: felt, classHash: felt):
end

@storage_var
func PluginEnabled(plugin: felt) -> (res: felt):
end

namespace PluginUtils:
    func initializer{
            syscall_ptr: felt*,
            pedersen_ptr: HashBuiltin*,
            range_check_ptr
        } (
            plugin: felt
        ):
        with_attr error_message("plugin class hash cannnot be 0"):
            assert_not_zero(plugin)
        end
        PluginEnabled.write(plugin, 1)
        return()
    end

    func enable{
            syscall_ptr: felt*,
            pedersen_ptr: HashBuiltin*,
            range_check_ptr
        } (
            plugin: felt
        ):
        Account.assert_only_self()

        with_attr error_message("plugin class hash cannnot be 0"):
            assert_not_zero(plugin)
        end
        PluginEnabled.write(plugin, 1)
        return()
    end

    func is_enabled{
            syscall_ptr: felt*, 
            pedersen_ptr: HashBuiltin*,
            range_check_ptr
        } (plugin: felt) -> (success: felt):
        let (res) = PluginEnabled.read(plugin)
        return (success=res)
    end

    func disable{
            syscall_ptr: felt*,
            pedersen_ptr: HashBuiltin*,
            range_check_ptr
        } (
            plugin: felt
        ):
        Account.assert_only_self()

        PluginEnabled.write(plugin, 0)
        return()
    end

    func validate{
            syscall_ptr: felt*,
            pedersen_ptr: HashBuiltin*,
            ecdsa_ptr: SignatureBuiltin*,
            range_check_ptr
        } (
            call_array_len: felt,
            call_array: AccountCallArray*,
            calldata_len: felt,
            calldata: felt*
        ):
        alloc_locals

        let plugin_class = calldata[call_array[0].data_offset]
        let (enabled) = PluginEnabled.read(plugin_class)
        with_attr error_message("plugin class not enabled"):
            assert_not_zero(enabled)
        end

        IPlugin.library_call_validate(
            class_hash=plugin_class,
            plugin_data_len=call_array[0].data_len - 1,
            plugin_data=calldata + call_array[0].data_offset + 1,
            call_array_len=call_array_len - 1,
            call_array=call_array + AccountCallArray.SIZE,
            calldata_len=calldata_len - call_array[0].data_len,
            calldata=calldata + call_array[0].data_offset + call_array[0].data_len)
        return()
    end
end
