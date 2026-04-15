package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null/v6"
)

// Organization is the top-level aggregate of the organisational
// structure. Every clinic and every department ultimately belongs to
// exactly one Organization, and that parenthood is immutable.
type Organization struct {
	ID           uuid.UUID   `gorm:"primaryKey;<-:create"`
	Name         string      `gorm:"<-"`
	Description  null.String `gorm:"<-"`
	LegalAddress Address     `gorm:"embedded;embeddedPrefix:legal_address_;<-"`
	CreatedAt    time.Time   `gorm:"<-:create"`
	UpdatedAt    time.Time   `gorm:"<-"`
}

// TableName pins the gorm mapping to the domain schema.
func (Organization) TableName() string { return "domain.organizations" }
