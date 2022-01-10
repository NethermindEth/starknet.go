package caigo

/*
	Although the library adheres to the 'elliptic/curve' interface.
	All testing has been done against library function explicity.
	It is recommended to use in the same way(i.e. `curve.Sign` and not `ecdsa.Sign`).
*/
import (
	"crypto/elliptic"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"os"
)

var sc StarkCurve

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
	Alpha            *big.Int
	ConstantPoints   [][]*big.Int
}

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

func SC() StarkCurve {
	InitCurve()
	return sc
}

func SCWithConstants(path string) (StarkCurve, error) {
	err := InitWithConstants(path)
	return sc, err
}

/*
	Not all operations require a stark curve initialization
	including the provided constant points. Here you can
	initialize the curve without the constant points
*/
func InitCurve() {
	sc.CurveParams = &elliptic.CurveParams{Name: "stark-curve"}
	sc.P, _ = new(big.Int).SetString("3618502788666131213697322783095070105623107215331596699973092056135872020481", 10)  // Field Prime ./pedersen_json
	sc.N, _ = new(big.Int).SetString("3618502788666131213697322783095070105526743751716087489154079457884512865583", 10)  // Order of base point ./pedersen_json
	sc.B, _ = new(big.Int).SetString("3141592653589793238462643383279502884197169399375105820974944592307816406665", 10)  // Constant of curve equation ./pedersen_json
	sc.Gx, _ = new(big.Int).SetString("2089986280348253421170679821480865132823066470938446095505822317253594081284", 10) // (x, _) of basepoint ./pedersen_json
	sc.Gy, _ = new(big.Int).SetString("1713931329540660377023406109199410414810705867260802078187082345529207694986", 10) // (_, y) of basepoint ./pedersen_json
	sc.EcGenX, _ = new(big.Int).SetString("874739451078007766457464989774322083649278607533249481151382481072868806602", 10)
	sc.EcGenY, _ = new(big.Int).SetString("152666792071518830868575557812948353041420400780739481342941381225525861407", 10)
	sc.MinusShiftPointX, _ = new(big.Int).SetString("2089986280348253421170679821480865132823066470938446095505822317253594081284", 10) // MINUS_SHIFT_POINT = (SHIFT_POINT[0], FIELD_PRIME - SHIFT_POINT[1])
	sc.MinusShiftPointY, _ = new(big.Int).SetString("1904571459125470836673916673895659690812401348070794621786009710606664325495", 10) // MINUS_SHIFT_POINT = (SHIFT_POINT[0], FIELD_PRIME - SHIFT_POINT[1])
	sc.Alpha = big.NewInt(1)
	sc.BitSize = 251
}

/*
	Various starknet functions require constant points be initialized.
	In this case use 'InitWithConstants'. Given an empty string this will
	init the curve by pulling the 'pedersen_params.json' file from Starkware
	official github repository. For production deployments it is recommended
	to have the file stored locally.
*/
func InitWithConstants(path string) (err error) {
	sc.CurveParams = &elliptic.CurveParams{Name: "stark-curve-with-constants"}
	scPayload := &StarkCurvePayload{}

	if path != "" {
		scFile, err := os.Open(path)
		if err != nil {
			return err
		}
		defer scFile.Close()

		scBytes, err := ioutil.ReadAll(scFile)
		if err != nil {
			return err
		}

		json.Unmarshal(scBytes, &scPayload)
	} else {
		url := "https://raw.githubusercontent.com/starkware-libs/cairo-lang/master/src/starkware/crypto/starkware/crypto/signature/pedersen_params.json"
		method := "GET"

		client := &http.Client{}

		req, err := http.NewRequest(method, url, nil)
		if err != nil {
			return err
		}
		req.Header.Add("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		err = json.NewDecoder(resp.Body).Decode(scPayload)
		if err != nil {
			return err
		}
	}

	if len(scPayload.ConstantPoints) == 0 {
		return fmt.Errorf("could not decode stark curve json")
	}

	sc.P = scPayload.FieldPrime
	sc.N = scPayload.EcOrder
	sc.B = scPayload.Beta
	sc.Gx = scPayload.ConstantPoints[0][0]
	sc.Gy = scPayload.ConstantPoints[0][1]
	sc.EcGenX = scPayload.ConstantPoints[1][0]
	sc.EcGenY = scPayload.ConstantPoints[1][1]
	sc.MinusShiftPointX, _ = new(big.Int).SetString("2089986280348253421170679821480865132823066470938446095505822317253594081284", 10) // MINUS_SHIFT_POINT = (SHIFT_POINT[0], FIELD_PRIME - SHIFT_POINT[1])
	sc.MinusShiftPointY, _ = new(big.Int).SetString("1904571459125470836673916673895659690812401348070794621786009710606664325495", 10)
	sc.Alpha = big.NewInt(scPayload.Alpha)
	sc.BitSize = 251
	sc.ConstantPoints = scPayload.ConstantPoints

	return nil
}

func (sc StarkCurve) Params() *elliptic.CurveParams {
	return sc.CurveParams
}

// Gets two points on an elliptic curve mod p and returns their sum.
// Assumes affine form (x, y) is spread (x1 *big.Int, y1 *big.Int)
//
// (ref: https://github.com/starkware-libs/cairo-lang/blob/master/src/starkware/crypto/starkware/crypto/signature/math_utils.py)
func (sc StarkCurve) Add(x1, y1, x2, y2 *big.Int) (x, y *big.Int) {
	yDelta := new(big.Int)
	xDelta := new(big.Int)
	yDelta.Sub(y1, y2)
	xDelta.Sub(x1, x2)

	m := DivMod(yDelta, xDelta, sc.P)

	xm := new(big.Int)
	xm = xm.Mul(m, m)

	x = new(big.Int)
	x = x.Sub(xm, x1)
	x = x.Sub(x, x2)
	x = x.Mod(x, sc.P)

	y = new(big.Int)
	y = y.Sub(x1, x)
	y = y.Mul(m, y)
	y = y.Sub(y, y1)
	y = y.Mod(y, sc.P)

	return x, y
}

// Doubles a point on an elliptic curve with the equation y^2 = x^3 + alpha*x + beta mod p.
// Assumes affine form (x, y) is spread (x1 *big.Int, y1 *big.Int)
//
// (ref: https://github.com/starkware-libs/cairo-lang/blob/master/src/starkware/crypto/starkware/crypto/signature/math_utils.py)
func (sc StarkCurve) Double(x1, y1 *big.Int) (x, y *big.Int) {
	xin := new(big.Int)
	xin = xin.Mul(big.NewInt(3), x1)
	xin = xin.Mul(xin, x1)
	xin = xin.Add(xin, sc.Alpha)

	yin := new(big.Int)
	yin = yin.Mul(y1, big.NewInt(2))

	m := DivMod(xin, yin, sc.P)

	xout := new(big.Int)
	xout = xout.Mul(m, m)
	xmed := new(big.Int)
	xmed = xmed.Mul(big.NewInt(2), x1)
	xout = xout.Sub(xout, xmed)
	xout = xout.Mod(xout, sc.P)

	yout := new(big.Int)
	yout = yout.Sub(x1, xout)
	yout = yout.Mul(m, yout)
	yout = yout.Sub(yout, y1)
	yout = yout.Mod(yout, sc.P)

	return xout, yout
}

func (sc StarkCurve) ScalarMult(x1, y1 *big.Int, k []byte) (x, y *big.Int) {
	fmt.Println("K HERE: ", k)
	m := new(big.Int)
	m = m.SetBytes(k)
	fmt.Println("M HERE: ", m)
	x, y = sc.EcMult(m, x1, y1)
	return x, y
}

func (sc StarkCurve) ScalarBaseMult(k []byte) (x, y *big.Int) {
	return sc.ScalarMult(sc.Gx, sc.Gy, k)
}

func (sc StarkCurve) IsOnCurve(x, y *big.Int) bool {
	left := new(big.Int)
	left = left.Mul(y, y)
	left = left.Mod(left, sc.P)

	right := new(big.Int)
	right = right.Mul(x, x)
	right = right.Mul(right, x)
	right = right.Mod(right, sc.P)

	ri := new(big.Int)
	// ALPHA = big.NewInt(1)
	ri = ri.Mul(big.NewInt(1), x)

	right = right.Add(right, ri)
	right = right.Add(right, sc.B)
	right = right.Mod(right, sc.P)

	if left.Cmp(right) == 0 {
		return true
	} else {
		return false
	}
}

// (ref: https://github.com/starkware-libs/cairo-lang/blob/master/src/starkware/crypto/starkware/crypto/signature/math_utils.py)
func (sc StarkCurve) InvModCurveSize(x *big.Int) *big.Int {
	return DivMod(big.NewInt(1), x, sc.N)
}

// Given the x coordinate of a stark_key, returns a possible y coordinate such that together the
// point (x,y) is on the curve.
// Note: the real y coordinate is either y or -y.
//
// (ref: https://github.com/starkware-libs/cairo-lang/blob/master/src/starkware/crypto/starkware/crypto/signature/signature.py)
func (sc StarkCurve) GetYCoordinate(starkX *big.Int) *big.Int {
	y := new(big.Int)
	y = y.Mul(starkX, starkX)
	y = y.Mul(y, starkX)
	yin := new(big.Int)
	yin = yin.Mul(sc.Alpha, starkX)

	y = y.Add(y, yin)
	y = y.Add(y, sc.B)
	y = y.Mod(y, sc.P)

	y = y.ModSqrt(y, sc.P)
	return y
}

// Computes m * point + shift_point using the same steps like the AIR and throws an exception if
// and only if the AIR errors.
//
// (ref: https://github.com/starkware-libs/cairo-lang/blob/master/src/starkware/crypto/starkware/crypto/signature/signature.py)
func (sc StarkCurve) MimicEcMultAir(mout, x1, y1, x2, y2 *big.Int) (x *big.Int, y *big.Int, err error) {
	m := new(big.Int)
	m = m.Set(mout)
	if m.Cmp(big.NewInt(0)) != 1 || m.BitLen() > 502 {
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

// Multiplies by m a point on the elliptic curve with equation y^2 = x^3 + alpha*x + beta mod p.
// Assumes affine form (x, y) is spread (x1 *big.Int, y1 *big.Int) and that 0 < m < order(point).
//
// (ref: https://github.com/starkware-libs/cairo-lang/blob/master/src/starkware/crypto/starkware/crypto/signature/math_utils.py)
func (sc StarkCurve) EcMult(m, x1, y1 *big.Int) (x, y *big.Int) {
	var _ecMult func(m, x1, y1 *big.Int) (x, y *big.Int)
	var _add func(x1, y1, x2, y2 *big.Int) (x, y *big.Int)

	_add = func(x1, y1, x2, y2 *big.Int) (x, y *big.Int) {
		yDelta := new(big.Int)
		xDelta := new(big.Int)
		yDelta.Sub(y1, y2)
		xDelta.Sub(x1, x2)

		m := DivMod(yDelta, xDelta, sc.P)

		xm := new(big.Int)
		xm = xm.Mul(m, m)

		x = new(big.Int)
		x = x.Sub(xm, x1)
		x = x.Sub(x, x2)
		x = x.Mod(x, sc.P)

		y = new(big.Int)
		y = y.Sub(x1, x)
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
		mk := new(big.Int)
		mk = mk.Mod(m, big.NewInt(2))
		if mk.Cmp(big.NewInt(0)) == 0 {
			h := new(big.Int)
			h = h.Div(m, big.NewInt(2))
			c, d := sc.Double(x1, y1)
			return _ecMult(h, c, d)
		}
		n := new(big.Int)
		n = n.Sub(m, big.NewInt(1))
		e, f := _ecMult(n, x1, y1)
		return _add(e, f, x1, y1)
	}

	x, y = _ecMult(m, x1, y1)
	return x, y
}

// Finds a nonnegative integer 0 <= x < p such that (m * x) % p == n
//
// (ref: https://github.com/starkware-libs/cairo-lang/blob/master/src/starkware/crypto/starkware/crypto/signature/math_utils.py)
func DivMod(n, m, p *big.Int) *big.Int {
	q := new(big.Int)
	gx := new(big.Int)
	gy := new(big.Int)
	q = q.GCD(gx, gy, m, p)

	r := new(big.Int)
	r = r.Mul(n, gx)
	r = r.Mod(r, p)
	return r
}
