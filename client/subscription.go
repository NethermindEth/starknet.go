// Copyright 2016 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package client

import (
	"bytes"
	"container/list"
	"context"
	crand "crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"math/rand"
	"reflect"
	"strings"
	"sync"
	"time"
)

var (
	// ErrNotificationsUnsupported is returned by the client when the connection doesn't
	// support notifications. You can use this error value to check for subscription
	// support like this:
	//
	//	sub, err := client.EthSubscribe(ctx, channel, "newHeads", true)
	//	if errors.Is(err, rpc.ErrNotificationsUnsupported) {
	//		// Server does not support subscriptions, fall back to polling.
	//	}
	//
	ErrNotificationsUnsupported = notificationsUnsupportedError{}

	// ErrSubscriptionNotFound is returned when the notification for the given id is not found
	ErrSubscriptionNotFound = errors.New("subscription not found")
)

var globalGen = randomIDGenerator()

// ID defines a pseudo random number that is used to identify RPC subscriptions.
type ID string

// NewID returns a new, random ID.
func NewID() ID {
	return globalGen()
}

// randomIDGenerator returns a function generates a random IDs.
func randomIDGenerator() func() ID {
	buf := make([]byte, 8)
	var seed int64
	if _, err := crand.Read(buf); err == nil {
		seed = int64(binary.BigEndian.Uint64(buf))
	} else {
		seed = int64(time.Now().Nanosecond())
	}

	var (
		mu  sync.Mutex
		rng = rand.New(rand.NewSource(seed))
	)

	return func() ID {
		mu.Lock()
		defer mu.Unlock()
		id := make([]byte, 16)
		rng.Read(id)

		return encodeID(id)
	}
}

func encodeID(b []byte) ID {
	id := hex.EncodeToString(b)
	id = strings.TrimLeft(id, "0")
	if id == "" {
		id = "0" // ID's are RPC quantities, no leading zero's and 0 is 0x0.
	}

	return ID("0x" + id)
}

type notifierKey struct{}

// NotifierFromContext returns the Notifier value stored in ctx, if any.
func NotifierFromContext(ctx context.Context) (*Notifier, bool) {
	n, ok := ctx.Value(notifierKey{}).(*Notifier)

	return n, ok
}

// Notifier is tied to an RPC connection that supports subscriptions.
// Server callbacks use the notifier to send notifications.
type Notifier struct {
	h         *handler
	namespace string

	mu           sync.Mutex
	sub          *Subscription
	buffer       []any
	callReturned bool
	activated    bool
}

// CreateSubscription returns a new subscription that is coupled to the
// RPC connection. By default subscriptions are inactive and notifications
// are dropped until the subscription is marked as active. This is done
// by the RPC server after the subscription ID is send to the client.
func (n *Notifier) CreateSubscription() *Subscription {
	n.mu.Lock()
	defer n.mu.Unlock()

	if n.sub != nil {
		panic("can't create multiple subscriptions with Notifier")
	} else if n.callReturned {
		panic("can't create subscription after subscribe call has returned")
	}
	n.sub = &Subscription{ID: n.h.idgen(), namespace: n.namespace, err: make(chan error, 1)}

	return n.sub
}

// Notify sends a notification to the client with the given data as payload.
// If an error occurs the RPC connection is closed and the error is returned.
func (n *Notifier) Notify(id ID, data any) error {
	n.mu.Lock()
	defer n.mu.Unlock()

	if n.sub == nil {
		panic("can't Notify before subscription is created")
	} else if n.sub.ID != id {
		panic("Notify with wrong ID")
	}
	if n.activated {
		return n.send(n.sub, data)
	}
	n.buffer = append(n.buffer, data)

	return nil
}

// takeSubscription returns the subscription (if one has been created). No subscription can
// be created after this call.
func (n *Notifier) takeSubscription() *Subscription {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.callReturned = true

	return n.sub
}

// activate is called after the subscription ID was sent to client. Notifications are
// buffered before activation. This prevents notifications being sent to the client before
// the subscription ID is sent to the client.
func (n *Notifier) activate() error {
	n.mu.Lock()
	defer n.mu.Unlock()

	for _, data := range n.buffer {
		if err := n.send(n.sub, data); err != nil {
			return err
		}
	}
	n.activated = true

	return nil
}

func (n *Notifier) send(sub *Subscription, data any) error {
	msg := jsonrpcSubscriptionNotification{
		Version: vsn,
		Method:  n.namespace + notificationMethodSuffix,
		Params: subscriptionResultEnc{
			ID:     string(sub.ID),
			Result: data,
		},
	}

	return n.h.conn.writeJSON(context.Background(), &msg, false)
}

// A Subscription is created by a notifier and tied to that notifier. The client can use
// this subscription to wait for an unsubscribe request for the client, see Err().
type Subscription struct {
	ID        ID
	namespace string
	err       chan error // closed on unsubscribe
}

// Err returns a channel that is closed when the client send an unsubscribe request.
func (s *Subscription) Err() <-chan error {
	return s.err
}

// MarshalJSON marshals a subscription as its ID.
func (s *Subscription) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.ID)
}

// ClientSubscription is a subscription established through the Client's Subscribe
type ClientSubscription struct {
	client       *Client
	etype        reflect.Type
	channel      reflect.Value
	reorgEtype   reflect.Type
	reorgChannel chan *ReorgEvent
	namespace    string
	subid        string

	// The in channel receives notification values from client dispatcher.
	in chan json.RawMessage

	// The error channel receives the error from the forwarding loop.
	// It is closed by Unsubscribe.
	err     chan error
	errOnce sync.Once

	// Closing of the subscription is requested by sending on 'quit'. This is handled by
	// the forwarding loop, which closes 'forwardDone' when it has stopped sending to
	// sub.channel. Finally, 'unsubDone' is closed after unsubscribing on the server side.
	quit        chan error
	forwardDone chan struct{}
	unsubDone   chan struct{}
}

// This is the sentinel value sent on sub.quit when Unsubscribe is called.
var errUnsubscribed = errors.New("unsubscribed")

func NewClientSubscription(c *Client, namespace string, channel reflect.Value) *ClientSubscription {
	sub := &ClientSubscription{
		client:       c,
		namespace:    namespace,
		etype:        channel.Type().Elem(),
		channel:      channel,
		reorgEtype:   reflect.TypeOf(&ReorgEvent{}),
		reorgChannel: make(chan *ReorgEvent),
		in:           make(chan json.RawMessage),
		quit:         make(chan error),
		forwardDone:  make(chan struct{}),
		unsubDone:    make(chan struct{}),
		err:          make(chan error, 1),
	}

	return sub
}

// Err returns the subscription error channel. The intended use of Err is to schedule
// resubscription when the client connection is closed unexpectedly.
//
// The error channel receives a value when the subscription has ended due to an error. The
// received error is nil if Close has been called on the underlying client and no other
// error has occurred.
//
// The error channel is closed when Unsubscribe is called on the subscription.
func (sub *ClientSubscription) Err() <-chan error {
	return sub.err
}

// Reorg returns a channel that notifies the subscriber of a reorganisation of the chain.
// A reorg event can be received from subscribing to any Starknet subscription.
func (sub *ClientSubscription) Reorg() <-chan *ReorgEvent {
	return sub.reorgChannel
}

// Unsubscribe unsubscribes the notification by calling the 'starknet_unsubscribe' method and closes the error channel.
// It can safely be called more than once.
func (sub *ClientSubscription) Unsubscribe() {
	sub.errOnce.Do(func() {
		select {
		case sub.quit <- errUnsubscribed:
			<-sub.unsubDone
		case <-sub.unsubDone:
		}
		close(sub.err)
	})
}

// ID returns the subscription ID.
func (sub *ClientSubscription) ID() string {
	return sub.subid
}

// deliver is called by the client's message dispatcher to send a notification value.
func (sub *ClientSubscription) deliver(result json.RawMessage) (ok bool) {
	select {
	case sub.in <- result:
		return true
	case <-sub.forwardDone:
		return false
	}
}

// close is called by the client's message dispatcher when the connection is closed.
func (sub *ClientSubscription) close(err error) {
	select {
	case sub.quit <- err:
	case <-sub.forwardDone:
	}
}

// run is the forwarding loop of the subscription. It runs in its own goroutine and
// is launched by the client's handler after the subscription has been created.
func (sub *ClientSubscription) run() {
	defer close(sub.unsubDone)
	defer close(sub.reorgChannel)

	unsubscribe, err := sub.forward()

	// The client's dispatch loop won't be able to execute the unsubscribe call if it is
	// blocked in sub.deliver() or sub.close(). Closing forwardDone unblocks them.
	close(sub.forwardDone)

	// Call the unsubscribe method on the server.
	if unsubscribe {
		_ = sub.requestUnsubscribe()
	}

	// Send the error.
	if err != nil {
		if err == ErrClientQuit {
			// ErrClientQuit gets here when Client.Close is called. This is reported as a
			// nil error because it's not an error, but we can't close sub.err here.
			err = nil
		}
		sub.err <- err
	}
}

// forward is the forwarding loop. It takes in RPC notifications and sends them
// on the subscription channel.
func (sub *ClientSubscription) forward() (unsubscribeServer bool, err error) {
	cases := []reflect.SelectCase{
		{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(sub.quit)},
		{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(sub.in)},
		{Dir: reflect.SelectSend, Chan: sub.channel},
	}

	// a workaround to handle reorg events as it'll come in the same subscription
	casesWithReorg := []reflect.SelectCase{
		cases[0],
		cases[1],
		{Dir: reflect.SelectSend, Chan: reflect.ValueOf(sub.reorgChannel)},
	}

	buffer := list.New()
	reorgBuffer := list.New()

	for {
		var chosen int
		var recv reflect.Value
		var isReorg bool

		if buffer.Len() == 0 && reorgBuffer.Len() == 0 {
			// Idle, omit send cases.
			chosen, recv, _ = reflect.Select(cases[:2])
		} else {
			// Non-empty buffer, send the first queued item.

			if reorgBuffer.Len() > 0 {
				casesWithReorg[2].Send = reflect.ValueOf(reorgBuffer.Front().Value)
				chosen, recv, _ = reflect.Select(casesWithReorg)
				isReorg = true
			} else {
				cases[2].Send = reflect.ValueOf(buffer.Front().Value)
				chosen, recv, _ = reflect.Select(cases)
			}
		}

		switch chosen {
		case 0: // <-sub.quit
			if !recv.IsNil() {
				err = recv.Interface().(error)
			}
			if err == errUnsubscribed {
				// Exiting because Unsubscribe was called, unsubscribe on server.
				return true, nil
			}

			return false, err

		case 1: // <-sub.in
			resp, isReorgVal, err := sub.unmarshal(recv.Interface().(json.RawMessage))
			if err != nil {
				return true, err
			}
			if buffer.Len()+reorgBuffer.Len() == maxClientSubscriptionBuffer {
				return true, ErrSubscriptionQueueOverflow
			}

			if isReorgVal {
				reorgBuffer.PushBack(resp)
			} else {
				buffer.PushBack(resp)
			}

		case 2: // sub.channel<- || sub.reorgChannel<-
			if isReorg {
				casesWithReorg[2].Send = reflect.Value{} // Cleaning up memory
				reorgBuffer.Remove(reorgBuffer.Front())
			} else {
				cases[2].Send = reflect.Value{} // Cleaning up memory
				buffer.Remove(buffer.Front())
			}
		}
	}
}

func (sub *ClientSubscription) unmarshal(value json.RawMessage) (resp interface{}, isReorg bool, err error) {
	val := reflect.New(sub.etype)
	dec := json.NewDecoder(bytes.NewReader(value))
	dec.DisallowUnknownFields()
	err = dec.Decode(val.Interface())

	// If there's an error when unmarshalling to the main channel type, maybe it's a reorg event
	if err != nil && sub.reorgEtype != nil {
		reorgVal := reflect.New(sub.reorgEtype)
		dec := json.NewDecoder(bytes.NewReader(value))
		dec.DisallowUnknownFields()
		err2 := dec.Decode(reorgVal.Interface())
		if err2 != nil {
			return nil, false, errors.Join(err, err2)
		}

		return reorgVal.Elem().Interface(), true, nil
	}

	return val.Elem().Interface(), false, err
}

func (sub *ClientSubscription) requestUnsubscribe() error {
	var result interface{}
	ctx, cancel := context.WithTimeout(context.Background(), unsubscribeTimeout)
	defer cancel()

	err := sub.client.CallContextWithSliceArgs(ctx, &result, sub.namespace+unsubscribeMethodSuffix, sub.subid)

	return err
}
