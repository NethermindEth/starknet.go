package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/NethermindEth/starknet.go/client"
)

// The purpose of the Spy type is to spy on the subscriptions made by the client.
// It's used in the tests to observe and store the responses from the subscriptions.
type WSSpy struct {
	wsConn
	spyCh chan json.RawMessage // must have a buffer of 1
	debug bool
}

// The wsConn interface used in `rpc.websocket` tests.
// It's implemented by the `client.Client` type.
type wsConn interface {
	Subscribe(
		ctx context.Context,
		namespace string,
		methodSuffix string,
		channel interface{},
		args interface{},
	) (*client.ClientSubscription, error)
	SubscribeWithSliceArgs(
		ctx context.Context,
		namespace string,
		methodSuffix string,
		channel interface{},
		args ...interface{},
	) (*client.ClientSubscription, error)
	Close()
}

// The Spyer interface implemented by the Spy type.
type WSSpyer interface {
	wsConn
	SpyChannel() <-chan json.RawMessage
	ToggleDebug()
}

// Assert that the Spy type implements the WSSpyer interface.
var (
	_ WSSpyer = (*WSSpy)(nil)
)

// NewWSSpy creates a new WSSpy object. A new spy should be created for each subscription,
// since it uses a single channel to capture the notifications sent by the node.
//
// It takes a client wsConn as the first parameter and an optional debug parameter.
// The client wsConn is the interface that the spy will be based on.
// The debug parameter is a variadic parameter that specifies whether debug mode is enabled.
//
// Parameters:
//   - client: the interface that the spy will be based on
//   - debug: a boolean flag indicating whether to print debug information
//
// Returns:
//   - WsSpyer: a new WSSpy object that implements the WsSpyer interface
func NewWSSpy(client wsConn, debug ...bool) WSSpyer {
	d := false
	if len(debug) > 0 {
		d = debug[0]
	}

	return &WSSpy{
		wsConn: client,
		spyCh:  make(chan json.RawMessage, 1),
		debug:  d,
	}
}

// Subscribe calls the original Subscribe function with the given parameters
// and captures the notifications sent by the node.
func (s *WSSpy) Subscribe(
	ctx context.Context,
	namespace string,
	methodSuffix string,
	channel interface{},
	args interface{},
) (*client.ClientSubscription, error) {
	if s.debug {
		fmt.Printf("### Spy Debug mode: in parameters\n")
		fmt.Printf("   arg.(%T): %+v\n", args, args)
		PrettyPrint(args)
		fmt.Println("--------------------------------------------")
	}

	userCh := reflect.ValueOf(channel)
	eventType := reflect.TypeOf(channel).Elem()
	mainCh := make(chan json.RawMessage)

	go listenAndForward(ctx, mainCh, userCh, eventType, s.spyCh, &s.debug)

	sub, err := s.wsConn.Subscribe(ctx, namespace, methodSuffix, mainCh, args)
	if err != nil {
		return nil, err
	}

	return sub, nil
}

// SubscribeWithSliceArgs calls the original SubscribeWithSliceArgs function with the given parameters
// and captures the notifications sent by the node.
func (s *WSSpy) SubscribeWithSliceArgs(
	ctx context.Context,
	namespace string,
	methodSuffix string,
	channel interface{},
	args ...interface{},
) (*client.ClientSubscription, error) {
	if s.debug {
		fmt.Printf("### Spy Debug mode: in parameters\n")
		for i, v := range args {
			fmt.Printf("   Arg[%d].(%T): %+v\n", i, v, v)
			PrettyPrint(v)
			fmt.Println("--------------------------------------------")
		}
	}

	userCh := reflect.ValueOf(channel)
	eventType := reflect.TypeOf(channel).Elem()
	mainCh := make(chan json.RawMessage)

	go listenAndForward(ctx, mainCh, userCh, eventType, s.spyCh, &s.debug)

	sub, err := s.wsConn.SubscribeWithSliceArgs(ctx, namespace, methodSuffix, mainCh, args...)
	if err != nil {
		return nil, err
	}

	return sub, nil
}

// SpyChannel returns the channel that captures raw JSON responses from the node.
// In other words, it returns the raw JSON response received from the node when
// sending notifications to the subscription.
// It is filled right after the main channel is filled with a new message, so
// it should be read right after the main channel is read, not before.
func (s *WSSpy) SpyChannel() <-chan json.RawMessage {
	return s.spyCh
}

// Toggles the debug mode of the spy to the opposite of the current value.
func (s *WSSpy) ToggleDebug() {
	s.debug = !s.debug
}

// listenAndForward listens for messages on the main channel
// and forwards them to the user channel and the spy channel.
func listenAndForward(
	ctx context.Context,
	mainCh <-chan json.RawMessage,
	userCh reflect.Value,
	eventType reflect.Type,
	spyCh chan json.RawMessage,
	debug *bool,
) {
	for {
		select {
		case <-ctx.Done():
			return
		case rawMsg := <-mainCh:
			if *debug {
				fmt.Printf("### Spy Debug mode: msg received\n")
				PrettyPrint(rawMsg)
			}

			msg := reflect.New(eventType)
			dec := json.NewDecoder(bytes.NewReader(rawMsg))
			err := dec.Decode(msg.Interface())
			if err != nil {
				panic(fmt.Errorf(
					"failed to unmarshal message to variable of type %T: %w",
					msg.Interface(),
					err,
				))
			}
			userCh.Send(msg.Elem())

			// Non-blocking receive - discards value if present
			select {
			case <-spyCh:
				spyCh <- rawMsg
			default:
				spyCh <- rawMsg
			}
		}
	}
}
