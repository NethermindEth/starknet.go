%lang starknet

from starkware.cairo.common.cairo_builtins import HashBuiltin, SignatureBuiltin
from starkware.cairo.common.math import assert_not_zero
from starkware.starknet.common.syscalls import get_tx_info, get_caller_address

//###################
// STORAGE VARIABLES
//###################
@storage_var
func counter() -> (count: felt) {
}

@storage_var
func rand() -> (val: felt) {
}

//###################
// CONSTRUCTOR
//###################
@constructor
func constructor{syscall_ptr: felt*, pedersen_ptr: HashBuiltin*, range_check_ptr}() {
    counter.write(0);

    return ();
}

//###################
// GETTERS FUNCTION
//###################
@view
func get_count{syscall_ptr: felt*, pedersen_ptr: HashBuiltin*, range_check_ptr}() -> (count: felt) {
    let (count) = counter.read();

    return (count,);
}

@view
func get_rand{syscall_ptr: felt*, pedersen_ptr: HashBuiltin*, range_check_ptr}() -> (val: felt) {
    let (val) = rand.read();

    return (val,);
}

//###################
// EXTERNAL FUNCTIONS
//###################
@external
func set_rand{syscall_ptr: felt*, pedersen_ptr: HashBuiltin*, range_check_ptr}(val: felt) {
    rand.write(val);

    return ();
}

@external
func set_rand_signed{syscall_ptr: felt*, pedersen_ptr: HashBuiltin*, range_check_ptr}(val: felt) {
    let (caller) = get_caller_address();
    let (tx_info) = get_tx_info();

    assert_not_zero(tx_info.signature_len);

    rand.write(val);

    return ();
}

@external
func increment{syscall_ptr: felt*, pedersen_ptr: HashBuiltin*, range_check_ptr}() -> (count: felt) {
    let (count) = counter.read();
    counter.write(count + 1);

    let (new_count) = counter.read();

    return (count=new_count);
}

@external
func decrement{syscall_ptr: felt*, pedersen_ptr: HashBuiltin*, range_check_ptr}() -> (count: felt) {
    let (count) = counter.read();
    assert_not_zero(count);

    counter.write(count - 1);

    let (new_count) = counter.read();

    return (count=new_count);
}
