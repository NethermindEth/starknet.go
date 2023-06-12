package types

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

const FeltLength = 32

// Felt represents Field Element or Felt from cairo.
type Felt [FeltLength]byte

// Big converts a Felt to its big.Int representation.
func (f Felt) Big() *big.Int { return new(big.Int).SetBytes(f[:]) }

// Bytes gets the byte representation of the Felt.
func (f Felt) Bytes() []byte { return f[:] }

// StrToFelt converts a string containing a decimal, hexadecimal or UTF8 charset into a Felt.
func StrToFelt(str string) Felt {
	var f Felt
	f.strToFelt(str)
	return f
}

func (f *Felt) strToFelt(str string) bool {
	if b, ok := new(big.Int).SetString(str, 0); ok {
		b.FillBytes(f[:])
		return ok
	}

	// TODO: revisit conversation on seperate 'ShortString' conversion
	if asciiRegexp.MatchString(str) {
		hexStr := hex.EncodeToString([]byte(str))
		if b, ok := new(big.Int).SetString(hexStr, 16); ok {
			b.FillBytes(f[:])
			return ok
		}
	}

	return false
}

// BigToFelt converts a big.Int to its Felt representation.
func BigToFelt(b *big.Int) Felt {
	var f Felt
	b.FillBytes(f[:])
	return f
}

// BytesToFelt converts a []byte to its Felt representation.
func BytesToFelt(b []byte) Felt {
	var f Felt
	copy(f[:], b)
	return f
}

// String converts a Felt into its 'short string' representation.
func (f Felt) ShortString() string {
	str := string(f.Big().Bytes())
	if asciiRegexp.MatchString(str) {
		return str
	}
	return ""
}

// String converts a Felt into its hexadecimal string representation and implement fmt.Stringer.
func (f Felt) String() string {
	return fmt.Sprintf("0x%x", f.Big())
}

// MarshalJSON implements the json Marshaller interface for a Signature array to marshal types to []byte.
func (s Signature) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`["%s","%s"]`, s[0].String(), s[1].String())), nil
}

// MarshalJSON implements the json Marshaller interface for Felt to marshal types to []byte.
func (f Felt) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, f.String())), nil
}

// MarshalText implements encoding.TextMarshaler
func (f Felt) MarshalText() ([]byte, error) {
	return []byte(f.String()), nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (f *Felt) UnmarshalText(input []byte) error {
	if f.strToFelt(string(input)) {
		return nil
	} else {
		return fmt.Errorf("Could not unmarshal to felt")
	}
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

	if ok := f.strToFelt(s); !ok {
		return fmt.Errorf("unmarshalling big int: %s", string(p))
	}

	return nil
}

// MarshalGQL implements the gqlgen Marshaller interface to marshal Felt into an io.Writer.
func (f Felt) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(f.String()))
}

// UnmarshalGQL implements the gqlgen Unmarshaller interface to unmarshal an interface into a Felt.
func (f *Felt) UnmarshalGQL(v interface{}) error {
	switch bi := v.(type) {
	case string:
		if ok := f.strToFelt(bi); ok {
			return nil
		}
	case int:
		b := big.NewInt(int64(bi))
		b.FillBytes(f[:])
		return nil
	}

	return fmt.Errorf("invalid big number")
}

// Value is used by database/sql drivers to store data in databases
func (f Felt) Value() (driver.Value, error) {
	return f.String(), nil
}

// Scan implements the database/sql Scanner interface to read Felt from a databases.
func (f Felt) Scan(src interface{}) error {
	var i sql.NullString
	if err := i.Scan(src); err != nil {
		return err
	}
	if !i.Valid {
		return nil
	}
	// Value came in a floating point format.
	if strings.ContainsAny(i.String, ".+e") {
		flt := big.NewFloat(0)
		if _, err := fmt.Sscan(i.String, flt); err != nil {
			return err
		}
		i, _ := flt.Int(nil)
		i.FillBytes(f[:])
	} else {
		intValue := big.NewInt(0)
		if _, err := fmt.Sscan(i.String, intValue); err != nil {
			return err
		}
		intValue.FillBytes(f[:])
	}
	return nil
}
