package model

import (
	"database/sql/driver"
	"fmt"
	"strings"

	"github.com/guregu/null/v6"
	"github.com/samber/oops"
)

// Address is the persistence form of an address. Text is always
// required. Point is optional — Valid=false means no coordinates.
type Address struct {
	Text  string            `gorm:"column:text;<-"`
	Point null.Value[Point] `gorm:"column:point;<-"`
}

// Equal reports whether two addresses have the same text and point.
func (a Address) Equal(other Address) bool {
	if a.Text != other.Text {
		return false
	}
	if !a.Point.Valid && !other.Point.Valid {
		return true
	}
	if !a.Point.Valid || !other.Point.Valid {
		return false
	}
	return a.Point.V.Equal(other.Point.V)
}

// Value serialises Address into the Postgres composite text format
// accepted by the domain.address type.
func (a Address) Value() (driver.Value, error) {
	textStr := quoteCompositeField(a.Text)
	if !a.Point.Valid {
		return fmt.Sprintf("(%s,)", textStr), nil
	}
	v, err := a.Point.V.Value()
	if err != nil {
		return nil, err
	}
	pointStr := quoteCompositeField(v.(string))
	return fmt.Sprintf("(%s,%s)", textStr, pointStr), nil
}

// Scan parses the Postgres composite text format for domain.address.
func (a *Address) Scan(src any) error {
	if src == nil {
		return oops.Errorf("cannot scan nil into Address")
	}
	var s string
	switch v := src.(type) {
	case string:
		s = v
	case []byte:
		s = string(v)
	default:
		return oops.Errorf("unsupported scan type %T for Address", src)
	}

	// Strip outer parens.
	if len(s) < 2 || s[0] != '(' || s[len(s)-1] != ')' {
		return oops.Errorf("invalid composite format: %q", s)
	}
	inner := s[1 : len(s)-1]

	fields := splitCompositeFields(inner)
	if len(fields) != 2 {
		return oops.Errorf("expected 2 fields in address composite, got %d", len(fields))
	}

	// Field 0: text (required).
	a.Text = unquoteCompositeField(fields[0])

	// Field 1: point (optional — empty string means NULL).
	if fields[1] == "" {
		a.Point = null.Value[Point]{}
	} else {
		raw := unquoteCompositeField(fields[1])
		var p Point
		if err := p.Scan(raw); err != nil {
			return oops.Wrap(err)
		}
		a.Point = null.ValueFrom(p)
	}
	return nil
}

// quoteCompositeField wraps s in double-quotes, escaping inner
// double-quotes and backslashes per Postgres composite text rules.
func quoteCompositeField(s string) string {
	var b strings.Builder
	b.WriteByte('"')
	for _, c := range s {
		if c == '"' || c == '\\' {
			b.WriteByte('\\')
		}
		b.WriteRune(c)
	}
	b.WriteByte('"')
	return b.String()
}

// splitCompositeFields splits a composite inner string (without outer parens)
// into fields, respecting double-quoted values and nested composites (parens).
func splitCompositeFields(s string) []string {
	var fields []string
	var current strings.Builder
	inQuote := false
	depth := 0

	for i := 0; i < len(s); i++ {
		c := s[i]
		if inQuote {
			if c == '\\' && i+1 < len(s) {
				current.WriteByte(s[i+1])
				i++
				continue
			}
			if c == '"' {
				// Check for escaped quote ("") — Postgres uses "" not \".
				if i+1 < len(s) && s[i+1] == '"' {
					current.WriteByte('"')
					i++
					continue
				}
				inQuote = false
				continue
			}
			current.WriteByte(c)
			continue
		}
		switch {
		case c == '"':
			inQuote = true
		case c == '(':
			depth++
			current.WriteByte(c)
		case c == ')':
			depth--
			current.WriteByte(c)
		case c == ',' && depth == 0:
			fields = append(fields, current.String())
			current.Reset()
		default:
			current.WriteByte(c)
		}
	}
	fields = append(fields, current.String())
	return fields
}

// unquoteCompositeField removes surrounding double-quotes and unescapes
// inner quotes. If the field is not quoted, returns as-is.
func unquoteCompositeField(s string) string {
	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		inner := s[1 : len(s)-1]
		inner = strings.ReplaceAll(inner, `""`, `"`)
		inner = strings.ReplaceAll(inner, `\"`, `"`)
		inner = strings.ReplaceAll(inner, `\\`, `\`)
		return inner
	}
	return s
}
