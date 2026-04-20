package utils

import (
	"testing"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEventWith(t *testing.T) {
	transferSelector := GetSelectorFromNameFelt("Transfer")
	approvalSelector := GetSelectorFromNameFelt("Approval")

	fromAddr := new(felt.Felt).SetUint64(12345)
	dataFelt := new(felt.Felt).SetUint64(100)

	tests := []struct {
		name     string
		events   []rpc.Event
		eventKey string
		want     *rpc.Event
	}{
		{
			name:     "empty slice returns nil",
			events:   []rpc.Event{},
			eventKey: "Transfer",
			want:     nil,
		},
		{
			name: "event not found returns nil",
			events: []rpc.Event{
				{
					FromAddress: fromAddr,
					EventContent: rpc.EventContent{
						Keys: []*felt.Felt{approvalSelector},
						Data: []*felt.Felt{dataFelt},
					},
				},
			},
			eventKey: "Transfer",
			want:     nil,
		},
		{
			name: "event found returns correct event",
			events: []rpc.Event{
				{
					FromAddress: fromAddr,
					EventContent: rpc.EventContent{
						Keys: []*felt.Felt{transferSelector},
						Data: []*felt.Felt{dataFelt},
					},
				},
			},
			eventKey: "Transfer",
			want: &rpc.Event{
				FromAddress: fromAddr,
				EventContent: rpc.EventContent{
					Keys: []*felt.Felt{transferSelector},
					Data: []*felt.Felt{dataFelt},
				},
			},
		},
		{
			name: "multiple events returns first match",
			events: []rpc.Event{
				{
					FromAddress: fromAddr,
					EventContent: rpc.EventContent{
						Keys: []*felt.Felt{approvalSelector},
						Data: []*felt.Felt{dataFelt},
					},
				},
				{
					FromAddress: fromAddr,
					EventContent: rpc.EventContent{
						Keys: []*felt.Felt{transferSelector},
						Data: []*felt.Felt{new(felt.Felt).SetUint64(200)},
					},
				},
				{
					FromAddress: fromAddr,
					EventContent: rpc.EventContent{
						Keys: []*felt.Felt{transferSelector},
						Data: []*felt.Felt{new(felt.Felt).SetUint64(300)},
					},
				},
			},
			eventKey: "Transfer",
			want: &rpc.Event{
				FromAddress: fromAddr,
				EventContent: rpc.EventContent{
					Keys: []*felt.Felt{transferSelector},
					Data: []*felt.Felt{new(felt.Felt).SetUint64(200)},
				},
			},
		},
		{
			name: "event with empty keys is skipped",
			events: []rpc.Event{
				{
					FromAddress: fromAddr,
					EventContent: rpc.EventContent{
						Keys: []*felt.Felt{},
						Data: []*felt.Felt{dataFelt},
					},
				},
				{
					FromAddress: fromAddr,
					EventContent: rpc.EventContent{
						Keys: []*felt.Felt{transferSelector},
						Data: []*felt.Felt{dataFelt},
					},
				},
			},
			eventKey: "Transfer",
			want: &rpc.Event{
				FromAddress: fromAddr,
				EventContent: rpc.EventContent{
					Keys: []*felt.Felt{transferSelector},
					Data: []*felt.Felt{dataFelt},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := EventWith(tt.events, tt.eventKey)
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

func TestEventsWithKey(t *testing.T) {
	transferSelector := GetSelectorFromNameFelt("Transfer")
	approvalSelector := GetSelectorFromNameFelt("Approval")

	fromAddr := new(felt.Felt).SetUint64(12345)
	dataFelt := new(felt.Felt).SetUint64(100)

	tests := []struct {
		name      string
		events    []rpc.Event
		eventKey  string
		wantCount int
	}{
		{
			name:      "empty slice returns empty slice",
			events:   []rpc.Event{},
			eventKey: "Transfer",
			wantCount: 0,
		},
		{
			name: "no matches returns empty slice",
			events: []rpc.Event{
				{
					FromAddress: fromAddr,
					EventContent: rpc.EventContent{
						Keys: []*felt.Felt{approvalSelector},
						Data: []*felt.Felt{dataFelt},
					},
				},
			},
			eventKey: "Transfer",
			wantCount: 0,
		},
		{
			name: "single match returns one event",
			events: []rpc.Event{
				{
					FromAddress: fromAddr,
					EventContent: rpc.EventContent{
						Keys: []*felt.Felt{transferSelector},
						Data: []*felt.Felt{dataFelt},
					},
				},
			},
			eventKey: "Transfer",
			wantCount: 1,
		},
		{
			name: "multiple matches returns all events",
			events: []rpc.Event{
				{
					FromAddress: fromAddr,
					EventContent: rpc.EventContent{
						Keys: []*felt.Felt{transferSelector},
						Data: []*felt.Felt{new(felt.Felt).SetUint64(100)},
					},
				},
				{
					FromAddress: fromAddr,
					EventContent: rpc.EventContent{
						Keys: []*felt.Felt{approvalSelector},
						Data: []*felt.Felt{new(felt.Felt).SetUint64(200)},
					},
				},
				{
					FromAddress: fromAddr,
					EventContent: rpc.EventContent{
						Keys: []*felt.Felt{transferSelector},
						Data: []*felt.Felt{new(felt.Felt).SetUint64(300)},
					},
				},
				{
					FromAddress: fromAddr,
					EventContent: rpc.EventContent{
						Keys: []*felt.Felt{transferSelector},
						Data: []*felt.Felt{new(felt.Felt).SetUint64(400)},
					},
				},
			},
			eventKey: "Transfer",
			wantCount: 3,
		},
		{
			name: "events with empty keys are skipped",
			events: []rpc.Event{
				{
					FromAddress: fromAddr,
					EventContent: rpc.EventContent{
						Keys: []*felt.Felt{},
						Data: []*felt.Felt{dataFelt},
					},
				},
				{
					FromAddress: fromAddr,
					EventContent: rpc.EventContent{
						Keys: []*felt.Felt{transferSelector},
						Data: []*felt.Felt{dataFelt},
					},
				},
			},
			eventKey: "Transfer",
			wantCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := EventsWithKey(tt.events, tt.eventKey)
			assert.Len(t, got, tt.wantCount)

			// Verify all returned events have the correct selector
			for _, event := range got {
				assert.Equal(t, transferSelector.String(), event.Keys[0].String())
			}
		})
	}
}
