package model

// Address is the persistence form of an address. Text is always
// required. Point is optional — nil means no coordinates.
type Address struct {
	Text  string `gorm:"column:text;<-"`
	Point *Point `gorm:"column:point;<-"`
}

// Equal reports whether two addresses have the same text and point.
func (a Address) Equal(other Address) bool {
	if a.Text != other.Text {
		return false
	}
	if a.Point == nil && other.Point == nil {
		return true
	}
	if a.Point == nil || other.Point == nil {
		return false
	}
	return a.Point.Equal(*other.Point)
}
