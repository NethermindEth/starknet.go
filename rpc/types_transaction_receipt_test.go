package rpc

import (
	"testing"

	"github.com/NethermindEth/juno/core/felt"
	internalUtils "github.com/NethermindEth/starknet.go/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTransactionReceiptWithBlockInfo_EventWith(t *testing.T) {
	transferSelector := internalUtils.GetSelectorFromNameFelt("Transfer")
	approvalSelector := internalUtils.GetSelectorFromNameFelt("Approval")

	fromAddr := new(felt.Felt).SetUint64(12345)
	dataFelt := new(felt.Felt).SetUint64(100)

	tests := []struct {
		name     string
		receipt  TransactionReceiptWithBlockInfo
		eventKey string
		want     *Event
	}{
		{
			name: "empty events returns nil",
			receipt: TransactionReceiptWithBlockInfo{
				TransactionReceipt: TransactionReceipt{
					Events: []Event{},
				},
			},
			eventKey: "Transfer",
			want:     nil,
		},
		{
			name: "event not found returns nil",
			receipt: TransactionReceiptWithBlockInfo{
				TransactionReceipt: TransactionReceipt{
					Events: []Event{
						{
							FromAddress: fromAddr,
							EventContent: EventContent{
								Keys: []*felt.Felt{approvalSelector},
								Data: []*felt.Felt{dataFelt},
							},
						},
					},
				},
			},
			eventKey: "Transfer",
			want:     nil,
		},
		{
			name: "event found returns correct event",
			receipt: TransactionReceiptWithBlockInfo{
				TransactionReceipt: TransactionReceipt{
					Events: []Event{
						{
							FromAddress: fromAddr,
							EventContent: EventContent{
								Keys: []*felt.Felt{transferSelector},
								Data: []*felt.Felt{dataFelt},
							},
						},
					},
				},
			},
			eventKey: "Transfer",
			want: &Event{
				FromAddress: fromAddr,
				EventContent: EventContent{
					Keys: []*felt.Felt{transferSelector},
					Data: []*felt.Felt{dataFelt},
				},
			},
		},
		{
			name: "multiple events returns first match",
			receipt: TransactionReceiptWithBlockInfo{
				TransactionReceipt: TransactionReceipt{
					Events: []Event{
						{
							FromAddress: fromAddr,
							EventContent: EventContent{
								Keys: []*felt.Felt{approvalSelector},
								Data: []*felt.Felt{dataFelt},
							},
						},
						{
							FromAddress: fromAddr,
							EventContent: EventContent{
								Keys: []*felt.Felt{transferSelector},
								Data: []*felt.Felt{new(felt.Felt).SetUint64(200)},
							},
						},
						{
							FromAddress: fromAddr,
							EventContent: EventContent{
								Keys: []*felt.Felt{transferSelector},
								Data: []*felt.Felt{new(felt.Felt).SetUint64(300)},
							},
						},
					},
				},
			},
			eventKey: "Transfer",
			want: &Event{
				FromAddress: fromAddr,
				EventContent: EventContent{
					Keys: []*felt.Felt{transferSelector},
					Data: []*felt.Felt{new(felt.Felt).SetUint64(200)},
				},
			},
		},
		{
			name: "event with empty keys is skipped",
			receipt: TransactionReceiptWithBlockInfo{
				TransactionReceipt: TransactionReceipt{
					Events: []Event{
						{
							FromAddress: fromAddr,
							EventContent: EventContent{
								Keys: []*felt.Felt{},
								Data: []*felt.Felt{dataFelt},
							},
						},
						{
							FromAddress: fromAddr,
							EventContent: EventContent{
								Keys: []*felt.Felt{transferSelector},
								Data: []*felt.Felt{dataFelt},
							},
						},
					},
				},
			},
			eventKey: "Transfer",
			want: &Event{
				FromAddress: fromAddr,
				EventContent: EventContent{
					Keys: []*felt.Felt{transferSelector},
					Data: []*felt.Felt{dataFelt},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.receipt.EventWith(tt.eventKey)
			if tt.want == nil {
				assert.Nil(t, got)
			} else {
				require.NotNil(t, got)
				assert.Equal(t, tt.want.FromAddress.String(), got.FromAddress.String())
				assert.Equal(t, tt.want.Keys[0].String(), got.Keys[0].String())
				assert.Equal(t, tt.want.Data[0].String(), got.Data[0].String())
			}
		})
	}
}
