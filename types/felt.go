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
	"regexp"
)

const (
	FIELD_GEN   int    = 3
	FIELD_PRIME string = "3618502788666131213697322783095070105623107215331596699973092056135872020481"
)

var (
	MaxFelt = StrToFelt(FIELD_PRIME)
	utfRegexp = regexp.MustCompile(`\w+`)
)

type Felt struct {
	*big.Int
}

func (f *Felt) Hex() string {
	return fmt.Sprintf("0x%x", f)
}

func (f *Felt) Big() *big.Int {
	return new(big.Int).SetBytes(f.Int.Bytes())
}

func StrToFelt(str string) *Felt {
	f := new(Felt)
	if ok := f.strToFelt(str); ok {
		return f
	}
	return nil
}

func BigToFelt(b *big.Int) *Felt {
	return &Felt{Int: b}
}

func BytesToFelt(b []byte) *Felt {
	return &Felt{Int: new(big.Int).SetBytes(b)}
}

func (f *Felt) strToFelt(str string) bool {
	if b, ok := new(big.Int).SetString(str, 0); ok {
		f.Int = b
		return ok
	}
	if IsUTF(str) {
		hexStr := hex.EncodeToString([]byte(str))
		if b, ok := new(big.Int).SetString(hexStr, 16); ok {
			f.Int = b
			return ok
		}
	}
	return false
}

func IsUTF(str string) bool {
	return utfRegexp.MatchString(str)
}

func (f Felt) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, f.String())), nil
}

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

func (f Felt) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(f.Hex()))
}

func (b *Felt) UnmarshalGQL(v interface{}) error {
	switch bi := v.(type) {
	case string:
		if ok := b.strToFelt(bi); ok {
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

func (f Felt) Value() (driver.Value, error) {
	if f.Int == nil {
		return "", nil
	}
	return f.String(), nil
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
