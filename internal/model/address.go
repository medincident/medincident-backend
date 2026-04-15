package model

// Address is the persistence form of an address. Text is always
// required (non-empty). Point is optional and embedded flat.
//
// Address itself is embedded into Organization/Clinic with an
// embeddedPrefix, so the owning table ends up with columns:
//
//	<prefix>_text
//	<prefix>_longitude
//	<prefix>_latitude
type Address struct {
	Text  string `gorm:"<-"`
	Point Point  `gorm:"embedded"`
}
