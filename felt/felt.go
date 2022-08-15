package felt

import (
	"database/sql"
	"database/sql/driver"
	"encoding/hex"
	"fmt"
	"io"
	"math/big"
	"regexp"
	"strconv"
	"strings"
)

const (
	FIELD_GEN   int    = 3
	FIELD_PRIME string = "3618502788666131213697322783095070105623107215331596699973092056135872020481"
)

var (
	MaxFelt     = StrToFelt(FIELD_PRIME)
	asciiRegexp = regexp.MustCompile(`^([[:graph:]]|[[:space:]]){1,31}$`)
)

// Felt represents Field Element or Felt from cairo.
type Felt struct {
	*big.Int
}

// Signature is a list of Felt.
type Signature []Felt

// IsNil tests if a Felt is nil
func (f *Felt) IsNil() bool {
	return f == nil || f.Int == nil
}

// Big converts a Felt to its big.Int representation.
func (f *Felt) BigInt() *big.Int {
	if f == nil || f.Int == nil {
		return &big.Int{}
	}
	return f.Int
}

// StrToFelt converts a string containing a decimal, hexadecimal or UTF8
// charset into a Felt.
func (f *Felt) StrToFelt(str string) bool {
	if b, ok := new(big.Int).SetString(str, 0); ok {
		f.Int = b
		return ok
	}

	// TODO: revisit conversation on separate 'ShortString' conversion
	if asciiRegexp.MatchString(str) {
		hexStr := hex.EncodeToString([]byte(str))
		if b, ok := new(big.Int).SetString(hexStr, 16); ok {
			f.Int = b
			return ok
		}
	}
	return false
}

// Add a Felt to an existing Felt
func NewFelt() *Felt {
	return &Felt{}
}

// Add a Felt to an existing Felt
func (f *Felt) Add(x, y Felt) *Felt {
	if f == nil {
		return nil
	}
	if x.IsNil() || y.IsNil() {
		f.Int = nil
		return f
	}
	if f.Int == nil {
		f.Int = new(big.Int)
	}
	f.Int.Add(x.Int, y.Int)
	return f
}

// Hex converts a Felt into its hexadecimal string representation.
func (f Felt) Hex() string {
	return fmt.Sprintf("0x%x", f)
}

// ShortString converts a Felt into its 'short string' representation.
func (f Felt) ShortString() string {
	if f.IsNil() {
		return ""
	}
	str := string(f.Bytes())
	if asciiRegexp.MatchString(str) {
		return str
	}
	return ""
}

// Equals compares 2 Felts and returns the true if they are not nil and equals
func (f Felt) Equals(g Felt) bool {
	if f.Int == nil || g.Int == nil || f.Cmp(g.Int) != 0 {
		return false
	}
	return true
}

// MarshalJSON implements the json Marshaller interface for a Signature array to marshal types to []byte.
func (s Signature) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`["%s","%s"]`, s[0].String(), s[1].String())), nil
}

// MarshalJSON implements the json Marshaller interface for Felt to marshal types to []byte.
func (f *Felt) MarshalJSON() ([]byte, error) {
	return []byte(f.String()), nil
}

// UnmarshalJSON implements the json Unmarshaller interface to unmarshal []byte into types.
func (f *Felt) UnmarshalJSON(p []byte) error {
	if string(p) == "null" || len(p) == 0 {
		return nil
	}

	var s string
	// parse double quotes
	if p[0] == 0x22 {
		s = string(p[1 : len(p)-1])
	} else {
		s = string(p)
	}

	if ok := f.StrToFelt(s); !ok {
		return fmt.Errorf("unmarshalling big int: %s", string(p))
	}

	return nil
}

// MarshalGQL implements the gqlgen Marshaller interface to marshal Felt into an io.Writer.
func (f Felt) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(f.String()))
}

// UnmarshalGQL implements the gqlgen Unmarshaller interface to unmarshal an interface into a Felt.
func (b *Felt) UnmarshalGQL(v interface{}) error {
	switch bi := v.(type) {
	case string:
		if ok := b.StrToFelt(bi); ok {
			return nil
		}
	case int:
		b.Int = big.NewInt(int64(bi))
		if b.Int != nil {
			return nil
		}
	}

	return fmt.Errorf("invalid big number")
}

// Value is used by database/sql drivers to store data in databases
func (f Felt) Value() (driver.Value, error) {
	if f.Int == nil {
		return "", nil
	}
	return f.String(), nil
}

// Scan implements the database/sql Scanner interface to read Felt from a databases.
func (f *Felt) Scan(src interface{}) error {
	var i sql.NullString
	if err := i.Scan(src); err != nil {
		return err
	}
	if !i.Valid {
		return nil
	}
	if f.Int == nil {
		f.Int = big.NewInt(0)
	}
	// Value came in a floating point format.
	if strings.ContainsAny(i.String, ".+e") {
		flt := big.NewFloat(0)
		if _, err := fmt.Sscan(i.String, f); err != nil {
			return err
		}
		f.Int, _ = flt.Int(f.Int)
	} else if _, err := fmt.Sscan(i.String, f.Int); err != nil {
		return err
	}
	return nil
}

// StrToFelt converts a string containing a decimal, hexadecimal or UTF8 charset into a Felt.
// mind there guesses and StrToFelt("1") will be 1 and not "1"
func StrToFelt(str string) Felt {
	f := Felt{}
	(&f).StrToFelt(str)
	return f
}

// BigToFelt converts a big.Int to its Felt representation.
func BigToFelt(b *big.Int) Felt {
	return Felt{Int: b}
}

// BytesToFelt converts a []byte to its Felt representation.
func BytesToFelt(b []byte) *Felt {
	return &Felt{Int: new(big.Int).SetBytes(b)}
}
