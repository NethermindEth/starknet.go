// Code generated by MockGen. DO NOT EDIT.
// Source: account.go
//
// Generated by this command:
//
//	mockgen -destination=../mocks/mock_account.go -package=mocks -source=account.go AccountInterface
//

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"
	time "time"

	felt "github.com/NethermindEth/juno/core/felt"
	rpc "github.com/NethermindEth/starknet.go/rpc"
	gomock "go.uber.org/mock/gomock"
)

// MockAccountInterface is a mock of AccountInterface interface.
type MockAccountInterface struct {
	ctrl     *gomock.Controller
	recorder *MockAccountInterfaceMockRecorder
}

// MockAccountInterfaceMockRecorder is the mock recorder for MockAccountInterface.
type MockAccountInterfaceMockRecorder struct {
	mock *MockAccountInterface
}

// NewMockAccountInterface creates a new mock instance.
func NewMockAccountInterface(ctrl *gomock.Controller) *MockAccountInterface {
	mock := &MockAccountInterface{ctrl: ctrl}
	mock.recorder = &MockAccountInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAccountInterface) EXPECT() *MockAccountInterfaceMockRecorder {
	return m.recorder
}

// PrecomputeAccountAddress mocks base method.
func (m *MockAccountInterface) PrecomputeAccountAddress(salt, classHash *felt.Felt, constructorCalldata []*felt.Felt) (*felt.Felt, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PrecomputeAccountAddress", salt, classHash, constructorCalldata)
	ret0, _ := ret[0].(*felt.Felt)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PrecomputeAccountAddress indicates an expected call of PrecomputeAccountAddress.
func (mr *MockAccountInterfaceMockRecorder) PrecomputeAccountAddress(salt, classHash, constructorCalldata any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PrecomputeAccountAddress", reflect.TypeOf((*MockAccountInterface)(nil).PrecomputeAccountAddress), salt, classHash, constructorCalldata)
}

// Sign mocks base method.
func (m *MockAccountInterface) Sign(ctx context.Context, msg *felt.Felt) ([]*felt.Felt, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Sign", ctx, msg)
	ret0, _ := ret[0].([]*felt.Felt)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Sign indicates an expected call of Sign.
func (mr *MockAccountInterfaceMockRecorder) Sign(ctx, msg any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Sign", reflect.TypeOf((*MockAccountInterface)(nil).Sign), ctx, msg)
}

// SignDeclareTransaction mocks base method.
func (m *MockAccountInterface) SignDeclareTransaction(ctx context.Context, tx *rpc.DeclareTxnV2) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SignDeclareTransaction", ctx, tx)
	ret0, _ := ret[0].(error)
	return ret0
}

// SignDeclareTransaction indicates an expected call of SignDeclareTransaction.
func (mr *MockAccountInterfaceMockRecorder) SignDeclareTransaction(ctx, tx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignDeclareTransaction", reflect.TypeOf((*MockAccountInterface)(nil).SignDeclareTransaction), ctx, tx)
}

// SignDeployAccountTransaction mocks base method.
func (m *MockAccountInterface) SignDeployAccountTransaction(ctx context.Context, tx *rpc.DeployAccountTxn, precomputeAddress *felt.Felt) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SignDeployAccountTransaction", ctx, tx, precomputeAddress)
	ret0, _ := ret[0].(error)
	return ret0
}

// SignDeployAccountTransaction indicates an expected call of SignDeployAccountTransaction.
func (mr *MockAccountInterfaceMockRecorder) SignDeployAccountTransaction(ctx, tx, precomputeAddress any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignDeployAccountTransaction", reflect.TypeOf((*MockAccountInterface)(nil).SignDeployAccountTransaction), ctx, tx, precomputeAddress)
}

// SignInvokeTransaction mocks base method.
func (m *MockAccountInterface) SignInvokeTransaction(ctx context.Context, tx *rpc.BroadcastInvokev1Txn) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SignInvokeTransaction", ctx, tx)
	ret0, _ := ret[0].(error)
	return ret0
}

// SignInvokeTransaction indicates an expected call of SignInvokeTransaction.
func (mr *MockAccountInterfaceMockRecorder) SignInvokeTransaction(ctx, tx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignInvokeTransaction", reflect.TypeOf((*MockAccountInterface)(nil).SignInvokeTransaction), ctx, tx)
}

// TransactionHashDeclare mocks base method.
func (m *MockAccountInterface) TransactionHashDeclare(tx rpc.DeclareTxnType) (*felt.Felt, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TransactionHashDeclare", tx)
	ret0, _ := ret[0].(*felt.Felt)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// TransactionHashDeclare indicates an expected call of TransactionHashDeclare.
func (mr *MockAccountInterfaceMockRecorder) TransactionHashDeclare(tx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TransactionHashDeclare", reflect.TypeOf((*MockAccountInterface)(nil).TransactionHashDeclare), tx)
}

// TransactionHashDeployAccount mocks base method.
func (m *MockAccountInterface) TransactionHashDeployAccount(tx rpc.DeployAccountType, contractAddress *felt.Felt) (*felt.Felt, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TransactionHashDeployAccount", tx, contractAddress)
	ret0, _ := ret[0].(*felt.Felt)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// TransactionHashDeployAccount indicates an expected call of TransactionHashDeployAccount.
func (mr *MockAccountInterfaceMockRecorder) TransactionHashDeployAccount(tx, contractAddress any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TransactionHashDeployAccount", reflect.TypeOf((*MockAccountInterface)(nil).TransactionHashDeployAccount), tx, contractAddress)
}

// TransactionHashInvoke mocks base method.
func (m *MockAccountInterface) TransactionHashInvoke(invokeTxn rpc.InvokeTxnType) (*felt.Felt, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TransactionHashInvoke", invokeTxn)
	ret0, _ := ret[0].(*felt.Felt)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// TransactionHashInvoke indicates an expected call of TransactionHashInvoke.
func (mr *MockAccountInterfaceMockRecorder) TransactionHashInvoke(invokeTxn any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TransactionHashInvoke", reflect.TypeOf((*MockAccountInterface)(nil).TransactionHashInvoke), invokeTxn)
}

// WaitForTransactionReceipt mocks base method.
func (m *MockAccountInterface) WaitForTransactionReceipt(ctx context.Context, transactionHash *felt.Felt, pollInterval time.Duration) (*rpc.TransactionReceiptWithBlockInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WaitForTransactionReceipt", ctx, transactionHash, pollInterval)
	ret0, _ := ret[0].(*rpc.TransactionReceiptWithBlockInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// WaitForTransactionReceipt indicates an expected call of WaitForTransactionReceipt.
func (mr *MockAccountInterfaceMockRecorder) WaitForTransactionReceipt(ctx, transactionHash, pollInterval any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WaitForTransactionReceipt", reflect.TypeOf((*MockAccountInterface)(nil).WaitForTransactionReceipt), ctx, transactionHash, pollInterval)
}
