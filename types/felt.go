package types

import (
	"database/sql"
	"database/sql/driver"
	"encoding/hex"
	"fmt"
	"io"
	"math/big"
	"strconv"
	"strings"
)

type Felt struct {
	*big.Int
}

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

func (f Felt) Value() (driver.Value, error) {
	if f.Int == nil {
		return "", nil
	}
	return f.String(), nil
}

func (f Felt) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, f.String())), nil
}

func (f *Felt) UnmarshalJSON(p []byte) error {
	if string(p) == "null" {
		return nil
	}

	f.Int = new(big.Int)

	// Ints are represented as strings in JSON. We
	// remove the enclosing quotes to provide a plain
	// string number to SetString.
	s := string(p[1 : len(p)-1])

	if has0xPrefix(s) {
		s := s[2:]

		if len(s)%2 == 1 {
			s = "0" + s
		}

		b, err := hex.DecodeString(s)
		if err != nil {
			return err
		}

		f.Int.SetBytes(b)
		return nil
	}

	if i, _ := f.Int.SetString(s, 10); i == nil {
		return fmt.Errorf("unmarshalling big int: %s", string(p))
	}

	return nil
}

func (f Felt) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(f.String()))
}

func (b *Felt) UnmarshalGQL(v interface{}) error {
	if bi, ok := v.(string); ok {
		b.Int = new(big.Int)
		b.Int, ok = b.Int.SetString(bi, 10)
		if !ok {
			return fmt.Errorf("invalid big number: %s", bi)
		}

		return nil
	}

	return fmt.Errorf("invalid big number")
}

func has0xPrefix(input string) bool {
	return len(input) >= 2 && input[0] == '0' && (input[1] == 'x' || input[1] == 'X')
}
