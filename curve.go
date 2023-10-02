package starknetgo

/*
	Although the library adheres to the 'elliptic/curve' interface.
	All testing has been done against library function explicity.
	It is recommended to use in the same way(i.e. `curve.Sign` and not `ecdsa.Sign`).
*/
import (
	"crypto/elliptic"
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
)

var Curve StarkCurve

/*
Returned stark curve includes several values above and beyond
what the 'elliptic' interface calls for to facilitate common starkware functions
*/
type StarkCurve struct {
	*elliptic.CurveParams
	EcGenX           *big.Int
	EcGenY           *big.Int
	MinusShiftPointX *big.Int
	MinusShiftPointY *big.Int
	Max              *big.Int
	Alpha            *big.Int
	ConstantPoints   [][]*big.Int
}

//go:embed pedersen_params.json
var PedersenParamsRaw []byte
var PedersenParams StarkCurvePayload

// struct definition for parsing 'pedersen_params.json'
type StarkCurvePayload struct {
	License        []string     `json:"_license"`
	Comment        string       `json:"_comment"`
	FieldPrime     *big.Int     `json:"FIELD_PRIME"`
	FieldGen       int          `json:"FIELD_GEN"`
	EcOrder        *big.Int     `json:"EC_ORDER"`
	Alpha          int64        `json:"ALPHA"`
	Beta           *big.Int     `json:"BETA"`
	ConstantPoints [][]*big.Int `json:"CONSTANT_POINTS"`
}

// init initializes the Curve parameters using the PedersenParams json.
//
// It unmarshals the PedersenParamsRaw json into the PedersenParams variable.
// It checks if the length of PedersenParams.ConstantPoints is zero and panics if true.
// It sets the Curve.CurveParams to a new elliptic.CurveParams with the name "stark-curve-with-constants".
// It sets the Curve.P, Curve.N, Curve.B, Curve.Gx, Curve.Gy, Curve.EcGenX, Curve.EcGenY, Curve.MinusShiftPointX, Curve.MinusShiftPointY, Curve.Max, Curve.Alpha, and Curve.BitSize fields.
//
// No parameter.
// No return value.
func init() {
	if err := json.Unmarshal(PedersenParamsRaw, &PedersenParams); err != nil {
		log.Fatalf("unmarshalling pedersen params: %v", err)
	}

	if len(PedersenParams.ConstantPoints) == 0 {
		panic("decoding pedersen params json")
	}

	Curve.CurveParams = &elliptic.CurveParams{Name: "stark-curve-with-constants"}

	Curve.P = PedersenParams.FieldPrime
	Curve.N = PedersenParams.EcOrder
	Curve.B = PedersenParams.Beta
	Curve.Gx = PedersenParams.ConstantPoints[0][0]
	Curve.Gy = PedersenParams.ConstantPoints[0][1]
	Curve.EcGenX = PedersenParams.ConstantPoints[1][0]
	Curve.EcGenY = PedersenParams.ConstantPoints[1][1]
	Curve.MinusShiftPointX, _ = new(big.Int).SetString("2089986280348253421170679821480865132823066470938446095505822317253594081284", 10) // MINUS_SHIFT_POINT = (SHIFT_POINT[0], FIELD_PRIME - SHIFT_POINT[1])
	Curve.MinusShiftPointY, _ = new(big.Int).SetString("1904571459125470836673916673895659690812401348070794621786009710606664325495", 10)
	Curve.Max, _ = new(big.Int).SetString("3618502788666131106986593281521497120414687020801267626233049500247285301248", 10) // 2 ** 251
	Curve.Alpha = big.NewInt(PedersenParams.Alpha)
	Curve.BitSize = 252
	Curve.ConstantPoints = PedersenParams.ConstantPoints

	/*
		Not all operations require a stark curve initialization
		including the provided constant points. Here you can
		initialize the curve without the constant points
	*/
	Curve.CurveParams = &elliptic.CurveParams{Name: "stark-curve"}
	Curve.P, _ = new(big.Int).SetString("3618502788666131213697322783095070105623107215331596699973092056135872020481", 10)  // Field Prime ./pedersen_json
	Curve.N, _ = new(big.Int).SetString("3618502788666131213697322783095070105526743751716087489154079457884512865583", 10)  // Order of base point ./pedersen_json
	Curve.B, _ = new(big.Int).SetString("3141592653589793238462643383279502884197169399375105820974944592307816406665", 10)  // Constant of curve equation ./pedersen_json
	Curve.Gx, _ = new(big.Int).SetString("2089986280348253421170679821480865132823066470938446095505822317253594081284", 10) // (x, _) of basepoint ./pedersen_json
	Curve.Gy, _ = new(big.Int).SetString("1713931329540660377023406109199410414810705867260802078187082345529207694986", 10) // (_, y) of basepoint ./pedersen_json
	Curve.EcGenX, _ = new(big.Int).SetString("874739451078007766457464989774322083649278607533249481151382481072868806602", 10)
	Curve.EcGenY, _ = new(big.Int).SetString("152666792071518830868575557812948353041420400780739481342941381225525861407", 10)
	Curve.MinusShiftPointX, _ = new(big.Int).SetString("2089986280348253421170679821480865132823066470938446095505822317253594081284", 10) // MINUS_SHIFT_POINT = (SHIFT_POINT[0], FIELD_PRIME - SHIFT_POINT[1])
	Curve.MinusShiftPointY, _ = new(big.Int).SetString("1904571459125470836673916673895659690812401348070794621786009710606664325495", 10) // MINUS_SHIFT_POINT = (SHIFT_POINT[0], FIELD_PRIME - SHIFT_POINT[1])
	Curve.Max, _ = new(big.Int).SetString("3618502788666131106986593281521497120414687020801267626233049500247285301248", 10)              // 2 ** 251
	Curve.Alpha = big.NewInt(1)
	Curve.BitSize = 252
}

// Add adds two points on a elliptic curve (mod p) and returns the resulting point.
// Assumes affine form (x, y) is spread (x1 *big.Int, y1 *big.Int)
//
// (ref: https://github.com/starkware-libs/cairo-lang/blob/master/src/starkware/crypto/starkware/crypto/signature/math_utils.py)
// The parameters x1, y1, x2, y2 are the x and y coordinates of the two points to be added.
// The function returns the x and y coordinates of the resulting point.
func (sc StarkCurve) Add(x1, y1, x2, y2 *big.Int) (x, y *big.Int) {
	yDelta := new(big.Int).Sub(y1, y2)
	xDelta := new(big.Int).Sub(x1, x2)

	m := DivMod(yDelta, xDelta, sc.P)

	xm := new(big.Int).Mul(m, m)

	x = new(big.Int).Sub(xm, x1)
	x = x.Sub(x, x2)
	x = x.Mod(x, sc.P)

	y = new(big.Int).Sub(x1, x)
	y = y.Mul(m, y)
	y = y.Sub(y, y1)
	y = y.Mod(y, sc.P)

	return x, y
}

// Doubles a point on an elliptic curve with the equation y^2 = x^3 + alpha*x + beta mod p.
// Assumes affine form (x, y) is spread (x1 *big.Int, y1 *big.Int)
// (ref: https://github.com/starkware-libs/cairo-lang/blob/master/src/starkware/crypto/starkware/crypto/signature/math_utils.py)
//
// Takes in the x and y coordinates of a point on the StarkCurve (elliptic curve) as
// big.Int pointers, and returns the x and y coordinates of the double
// of the input point as big.Int pointers.
func (sc StarkCurve) Double(x1, y1 *big.Int) (x, y *big.Int) {
	xin := new(big.Int).Mul(big.NewInt(3), x1)
	xin = xin.Mul(xin, x1)
	xin = xin.Add(xin, sc.Alpha)

	yin := new(big.Int).Mul(y1, big.NewInt(2))

	m := DivMod(xin, yin, sc.P)

	xout := new(big.Int).Mul(m, m)
	xmed := new(big.Int).Mul(big.NewInt(2), x1)
	xout = xout.Sub(xout, xmed)
	xout = xout.Mod(xout, sc.P)

	yout := new(big.Int).Sub(x1, xout)
	yout = yout.Mul(m, yout)
	yout = yout.Sub(yout, y1)
	yout = yout.Mod(yout, sc.P)

	return xout, yout
}

// ScalarMult multiplies the point (x1, y1) on the Stark elliptic curve by the scalar value k.
//
// Parameters:
// - x1: The x-coordinate of the point.
// - y1: The y-coordinate of the point.
// - k: The scalar value to multiply the point by.
//
// Returns:
// - x: The x-coordinate of the resulting point.
// - y: The y-coordinate of the resulting point.
func (sc StarkCurve) ScalarMult(x1, y1 *big.Int, k []byte) (x, y *big.Int) {
	m := new(big.Int).SetBytes(k)
	x, y = sc.EcMult(m, x1, y1)
	return x, y
}

// ScalarBaseMult returns the result of multiplying the base point of the Stark elliptic curve by a scalar.
//
// The parameter k is the scalar to multiply the base point by.
// The function returns the x and y coordinates of the resulting point as *big.Int.
func (sc StarkCurve) ScalarBaseMult(k []byte) (x, y *big.Int) {
	return sc.ScalarMult(sc.Gx, sc.Gy, k)
}

// IsOnCurve checks if the given point (x, y) is on the Stark elliptic curve.
//
// Parameters:
// - x: the x-coordinate of the point.
// - y: the y-coordinate of the point.
//
// Returns:
// - bool: true if the point is on the curve, false otherwise.
func (sc StarkCurve) IsOnCurve(x, y *big.Int) bool {
	left := new(big.Int).Mul(y, y)
	left = left.Mod(left, sc.P)

	right := new(big.Int).Mul(x, x)
	right = right.Mul(right, x)
	right = right.Mod(right, sc.P)

	ri := new(big.Int).Mul(big.NewInt(1), x)

	right = right.Add(right, ri)
	right = right.Add(right, sc.B)
	right = right.Mod(right, sc.P)

	if left.Cmp(right) == 0 {
		return true
	} else {
		return false
	}
}


// InvModCurveSize calculates the inverse of x modulo the size of the Stark elliptic curve.
// (ref: https://github.com/starkware-libs/cairo-lang/blob/master/src/starkware/crypto/starkware/crypto/signature/math_utils.py)
//
// It takes a pointer to a big.Int, x, as its parameter.
// It returns a pointer to a big.Int.
func (sc StarkCurve) InvModCurveSize(x *big.Int) *big.Int {
	return DivMod(big.NewInt(1), x, sc.N)
}

// GetYCoordinate calculates the Y coordinate for a given X coordinate on the Stark (elliptic) curve.
// (ref: https://github.com/starkware-libs/cairo-lang/blob/master/src/starkware/crypto/starkware/crypto/signature/signature.py)
//
// Parameters:
// - starkX: the X coordinate on the Stark curve.
//
// Returns:
// - y: the calculated Y returns a possible y coordinate such that together the point (x,y) is on the curve.
// Note: the real y coordinate is either y or -y.
func (sc StarkCurve) GetYCoordinate(starkX *big.Int) *big.Int {
	y := new(big.Int).Mul(starkX, starkX)
	y = y.Mul(y, starkX)
	yin := new(big.Int).Mul(sc.Alpha, starkX)

	y = y.Add(y, yin)
	y = y.Add(y, sc.B)
	y = y.Mod(y, sc.P)

	y = y.ModSqrt(y, sc.P)
	return y
}

// MimicEcMultAir performs a scalar multiplication (m * point) on a Stark elliptic curve using the given coordinates (+ shift_point) and returns the resulting x and y values.
// (ref: https://github.com/starkware-libs/cairo-lang/blob/master/src/starkware/crypto/starkware/crypto/signature/signature.py)
//
// Parameters:
// - mout: The scalar by which to multiply the curve.
// - x1, y1: The coordinates of the first point on the curve.
// - x2, y2: The coordinates of the second point on the curve.
//
// Returns:
// - x: The x-coordinate of the resulting point.
// - y: The y-coordinate of the resulting point.
// - err: An error, if any (if and only if the AIR errors).
func (sc StarkCurve) MimicEcMultAir(mout, x1, y1, x2, y2 *big.Int) (x *big.Int, y *big.Int, err error) {
	m := new(big.Int).Set(mout)
	if m.Cmp(big.NewInt(0)) != 1 || m.Cmp(sc.Max) != -1 {
		return x, y, fmt.Errorf("too many bits %v", m.BitLen())
	}

	psx := x2
	psy := y2
	for i := 0; i < 251; i++ {
		if psx == x1 {
			return x, y, fmt.Errorf("xs are the same")
		}
		if m.Bit(0) == 1 {
			psx, psy = sc.Add(psx, psy, x1, y1)
		}
		x1, y1 = sc.Double(x1, y1)
		m = m.Rsh(m, 1)
	}
	if m.Cmp(big.NewInt(0)) != 0 {
		return psx, psy, fmt.Errorf("m doesn't equal zero")
	}
	return psx, psy, nil
}

// EcMult performs elliptic curve multiplication using the Stark elliptic curve.
// Multiplies by m a point on the elliptic curve with equation y^2 = x^3 + alpha*x + beta mod p.
// Assumes affine form (x, y) is spread (x1 *big.Int, y1 *big.Int) and that 0 < m < order(point).
// (ref: https://github.com/starkware-libs/cairo-lang/blob/master/src/starkware/crypto/starkware/crypto/signature/math_utils.py)
//
// It takes three big.Int parameters: m, x1, and y1.
// It returns two big.Int values: x and y.
func (sc StarkCurve) EcMult(m, x1, y1 *big.Int) (x, y *big.Int) {
	var _ecMult func(m, x1, y1 *big.Int) (x, y *big.Int)

	_add := func(x1, y1, x2, y2 *big.Int) (x, y *big.Int) {
		yDelta := new(big.Int).Sub(y1, y2)
		xDelta := new(big.Int).Sub(x1, x2)

		m := DivMod(yDelta, xDelta, sc.P)

		xm := new(big.Int).Mul(m, m)

		x = new(big.Int).Sub(xm, x1)
		x = x.Sub(x, x2)
		x = x.Mod(x, sc.P)

		y = new(big.Int).Sub(x1, x)
		y = y.Mul(m, y)
		y = y.Sub(y, y1)
		y = y.Mod(y, sc.P)

		return x, y
	}

	// alpha is our Y
	_ecMult = func(m, x1, y1 *big.Int) (x, y *big.Int) {
		if m.BitLen() == 1 {
			return x1, y1
		}
		mk := new(big.Int).Mod(m, big.NewInt(2))
		if mk.Cmp(big.NewInt(0)) == 0 {
			h := new(big.Int).Div(m, big.NewInt(2))
			c, d := sc.Double(x1, y1)
			return _ecMult(h, c, d)
		}
		n := new(big.Int).Sub(m, big.NewInt(1))
		e, f := _ecMult(n, x1, y1)
		return _add(e, f, x1, y1)
	}

	x, y = _ecMult(m, x1, y1)
	return x, y
}

// DivMod finds a nonnegative integer 0 <= x < p such that (m * x) % p == n
// (ref: https://github.com/starkware-libs/cairo-lang/blob/master/src/starkware/crypto/starkware/crypto/signature/math_utils.py)
//
// Parameters:
// - n: The numerator of the division as a *big.Int.
// - m: The denominator of the division as a *big.Int.
// - p: The modulus value as a *big.Int.
//
// Returns:
// - r: The remainder of the division as a *big.Int.
func DivMod(n, m, p *big.Int) *big.Int {
	q := new(big.Int)
	gx := new(big.Int)
	gy := new(big.Int)
	q.GCD(gx, gy, m, p)

	r := new(big.Int).Mul(n, gx)
	r = r.Mod(r, p)
	return r
}
