package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null/v6"
)

// Clinic belongs to exactly one Organization; that link is set at
// creation time and never changes.
type Clinic struct {
	ID              uuid.UUID   `gorm:"primaryKey;<-:create"`
	OrganizationID  uuid.UUID   `gorm:"<-:create"`
	Name            string      `gorm:"<-"`
	Description     null.String `gorm:"<-"`
	PhysicalAddress Address     `gorm:"embedded;embeddedPrefix:physical_address_;<-"`
	CreatedAt       time.Time   `gorm:"<-:create"`
	UpdatedAt       time.Time   `gorm:"<-"`
}

func (Clinic) TableName() string { return "domain.clinics" }
