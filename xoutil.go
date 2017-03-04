// Package xoutil provides additional types for the xo cli tool.
package xoutil

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/mattn/go-sqlite3"
)

// SqTime provides a type that will correctly scan the various timestamps
// values stored by the github.com/mattn/go-sqlite3 driver for time.Time
// values, as well as correctly satisfying the sql/driver/Valuer interface.
type SqTime time.Time

// Value satisfies the Valuer interface.
func (t SqTime) Value() (driver.Value, error) {
	return time.Time(t), nil
}

// Scan satisfies the Scanner interface.
func (t *SqTime) Scan(v interface{}) error {
	switch x := v.(type) {
	case []byte:
		return t.parse(string(x))

	case string:
		return t.parse(x)
	}

	return fmt.Errorf("cannot convert type %s to time.Time", reflect.TypeOf(v))
}

// parse attempts to parse string s to t.
func (t *SqTime) parse(s string) error {
	if s == "" {
		return nil
	}

	for _, f := range sqlite3.SQLiteTimestampFormats {
		z, err := time.Parse(f, s)
		if err == nil {
			*t = SqTime(z)
			return nil
		}
	}

	return errors.New("could not parse time")
}
