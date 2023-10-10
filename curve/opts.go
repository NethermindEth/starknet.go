package curve

type curveOptions struct {
	initConstants bool
	paramsPath    string
}

// funcCurveOptions wraps a function that modifies curveOptions into an
// implementation of the CurveOption interface.
type funcCurveOption struct {
	f func(*curveOptions)
}

// apply applies the given curve options to the funcCurveOption.
//
// It takes a curveOptions pointer as a parameter and does not return anything.
func (fso *funcCurveOption) apply(do *curveOptions) {
	fso.f(do)
}

// newFuncCurveOption returns a new instance of funcCurveOption.
//
// It takes a function f as input, which is of type func(*curveOptions),
// and returns a pointer to funcCurveOption.
func newFuncCurveOption(f func(*curveOptions)) *funcCurveOption {
	return &funcCurveOption{
		f: f,
	}
}

type CurveOption interface {
	apply(*curveOptions)
}

// WithConstants creates a CurveOption (a curve initialized with constant points) that initializes the constants of the curve.
//
// paramsPath: A variadic parameter of type string, representing the path(s) to the parameters.
// Returns: A CurveOption.
func WithConstants(paramsPath ...string) CurveOption {
	return newFuncCurveOption(func(o *curveOptions) {
		o.initConstants = true

		if len(paramsPath) == 1 && paramsPath[0] != "" {
			o.paramsPath = paramsPath[0]
		}
	})
}
