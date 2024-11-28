package model

import (
	"database/sql/driver"
	"errors"
	"time"
)

type Date struct {
	time.Time
}

const dateFormat = "2006-01-02"

// UnmarshalJSON handles custom date parsing for JSON.
func (d *Date) UnmarshalJSON(b []byte) error {
	str := string(b)
	str = str[1 : len(str)-1] // Remove quotes around the date string
	t, err := time.Parse(dateFormat, str)
	if err != nil {
		return errors.New("invalid date format, use YYYY-MM-DD")
	}
	d.Time = t
	return nil
}

// MarshalJSON formats the date back to YYYY-MM-DD for JSON responses.
func (d Date) MarshalJSON() ([]byte, error) {
	return []byte(`"` + d.Format(dateFormat) + `"`), nil
}

// Value converts the Date to a driver-compatible value (gorm.Valuer).
func (d Date) Value() (driver.Value, error) {
	return d.Format(dateFormat), nil
}

// Scan assigns a value from the database to the Date (sql.Scanner).
func (d *Date) Scan(value interface{}) error {
	if value == nil {
		d.Time = time.Time{}
		return nil
	}

	switch v := value.(type) {
	case time.Time:
		d.Time = v
		return nil
	case string:
		t, err := time.Parse(dateFormat, v)
		if err != nil {
			return err
		}
		d.Time = t
		return nil
	default:
		return errors.New("invalid type for Date")
	}
}
