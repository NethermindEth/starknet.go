package types

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"io"
	"math/big"
	"strconv"
	"strings"
)

const (
	FIELD_GEN   int    = 3
	FIELD_PRIME string = "3618502788666131213697322783095070105623107215331596699973092056135872020481"
)

var (
	MaxFelt = ToFelt(FIELD_PRIME)
)

// Felt represents Field Element or Felt from cairo.
type Felt struct {
	*big.Int
}

// ToFelt converts a string containing a decimal or hexadecimal into a Felt.
func ToFelt(str string) *Felt {
	if b, ok := new(big.Int).SetString(str, 0); ok {
		return &Felt{Int: b}
	}
	return nil
}

func (f *Felt) setString(str string) bool {
	if b, ok := new(big.Int).SetString(str, 0); ok {
		f.Int = b
		return ok
	}
	return false
}

// String converts a Felt into its hexadecimal string representation and implement fmt.Stringer.
func (f *Felt) String() string {
	return fmt.Sprintf("0x%x", f)
}

// MarshalJSON implements the json Marshaller interface to marshal types to []byte.
func (f Felt) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, f.String())), nil
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

	if ok := f.setString(s); !ok {
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
		if ok := b.setString(bi); ok {
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
