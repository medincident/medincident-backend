package model

import "github.com/guregu/null/v6"

// Point is an optional geographic coordinate pair. Both fields must
// be Valid together or both Invalid — enforced by the service-layer
// validator, never by a DB check constraint.
//
// Point is embedded flat into Address (no extra prefix) so the owning
// table's columns are named <prefix>_longitude and <prefix>_latitude.
type Point struct {
	Longitude null.Float `gorm:"<-"`
	Latitude  null.Float `gorm:"<-"`
}
