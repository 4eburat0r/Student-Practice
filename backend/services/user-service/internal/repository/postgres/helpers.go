package postgres

import (
	"time"
)

// pgNullDate helps mapping DATE -> string (YYYY-MM-DD)
type pgNullDate struct {
	Valid bool
	Value string
}

func (n *pgNullDate) Scan(src interface{}) error {
	if src == nil {
		n.Valid = false
		n.Value = ""
		return nil
	}
	switch v := src.(type) {
	case time.Time:
		n.Valid = true
		n.Value = v.Format("2006-01-02")
		return nil
	case []byte:
		n.Valid = true
		n.Value = string(v)
		return nil
	case string:
		n.Valid = true
		n.Value = v
		return nil
	default:
		n.Valid = false
		n.Value = ""
		return nil
	}
}

func nullIfEmptyDate(s string) interface{} {
	if s == "" {
		return nil
	}
	if _, err := time.Parse("2006-01-02", s); err != nil {
		return nil
	}
	return s
}
