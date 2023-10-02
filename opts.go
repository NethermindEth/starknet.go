package starknetgo

type curveOptions struct {
	initConstants bool
	paramsPath    string
}

// funcCurveOptions wraps a function that modifies curveOptions into an
// implementation of the CurveOption interface.
type funcCurveOption struct {
	f func(*curveOptions)
}

// apply helper function that applies the funcCurveOption to the curveOptions.
//
// It takes a pointer to a curveOptions object as its parameter.
// There is no return value.
func (fso *funcCurveOption) apply(do *curveOptions) {
	fso.f(do)
}

// newFuncCurveOption creates a new funcCurveOption with the given function and returns a pointer to it.
//
// The function parameter f is a function that takes a pointer to a curveOptions struct as its argument.
// The return type of the function is a pointer to a funcCurveOption struct.
func newFuncCurveOption(f func(*curveOptions)) *funcCurveOption {
	return &funcCurveOption{
		f: f,
	}
}

type CurveOption interface {
	apply(*curveOptions)
}

// WithConstants creates a curve option with the specified parameter paths.
// functions that require Pedersen hashes must be run on
// a curve initialized with constant points
//
// paramsPath: A variadic list of strings representing the paths to the parameters.
// Returns: A CurveOption function.
func WithConstants(paramsPath ...string) CurveOption {
	return newFuncCurveOption(func(o *curveOptions) {
		o.initConstants = true

		if len(paramsPath) == 1 && paramsPath[0] != "" {
			o.paramsPath = paramsPath[0]
		}
	})
}
