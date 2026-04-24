package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null/v6"
)

// IncidentType is a leaf of the classifier that lives under exactly
// one IncidentCategory. CategoryID is mutable (MoveType). OrganizationID
// is denormalised for the same reason as in IncidentCategory — to power
// the active-only unique index.
type IncidentType struct {
	ID                   uuid.UUID   `gorm:"primaryKey;<-:create"`
	OrganizationID       uuid.UUID   `gorm:"<-:create"`
	CategoryID           uuid.UUID   `gorm:"<-"`
	Name                 string      `gorm:"<-"`
	Description          null.String `gorm:"<-"`
	IsActive             bool        `gorm:"<-"`
	IsAllowedForPatients bool        `gorm:"<-"`
	CreatedAt            time.Time   `gorm:"<-:create"`
	UpdatedAt            time.Time   `gorm:"<-"`
}

func (IncidentType) TableName() string { return "domain.incident_types" }
