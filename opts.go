package caigo

import (
	"net/http"
)

type gatewayOptions struct {
	client       *http.Client
	chainID      string
	errorHandler func(e error) error
}

type curveOptions struct {
	initConstants bool
	paramsPath    string
}

// funcGatewayOption wraps a function that modifies gatewayOptions into an
// implementation of the GatewayOption interface.
type funcGatewayOption struct {
	f func(*gatewayOptions)
}

func (fso *funcGatewayOption) apply(do *gatewayOptions) {
	fso.f(do)
}

func newFuncGatewayOption(f func(*gatewayOptions)) *funcGatewayOption {
	return &funcGatewayOption{
		f: f,
	}
}

// GatewayOption configures how we set up the connection.
type GatewayOption interface {
	apply(*gatewayOptions)
}

func WithHttpClient(client http.Client) GatewayOption {
	return newFuncGatewayOption(func(o *gatewayOptions) {
		o.client = &client
	})
}

func WithChain(chainID string) GatewayOption {
	return newFuncGatewayOption(func(o *gatewayOptions) {
		o.chainID = chainID
	})
}

// WithErrorHandler returns an Option to set the error handler to be used.
func WithErrorHandler(f func(e error) error) GatewayOption {
	return newFuncGatewayOption(func(o *gatewayOptions) {
		o.errorHandler = f
	})
}

// funcCurveOptions wraps a function that modifies curveOptions into an
// implementation of the CurveOption interface.
type funcCurveOption struct {
	f func(*curveOptions)
}

func (fso *funcCurveOption) apply(do *curveOptions) {
	fso.f(do)
}

func newFuncCurveOption(f func(*curveOptions)) *funcCurveOption {
	return &funcCurveOption{
		f: f,
	}
}

type CurveOption interface {
	apply(*curveOptions)
}

// functions that require pedersen hashes must be run on
// a curve initialized with constant points
func WithConstants(paramsPath ...string) CurveOption {
	return newFuncCurveOption(func(o *curveOptions) {
		o.initConstants = true

		if len(paramsPath) == 1 && paramsPath[0] != "" {
			o.paramsPath = paramsPath[0]
		}
	})
}
