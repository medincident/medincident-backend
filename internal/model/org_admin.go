package model

import (
	"time"

	"github.com/google/uuid"
)

// OrgAdmin links an organisation to an employee holding the
// organisation-admin role. Primary key is composite.
type OrgAdmin struct {
	OrganizationID   uuid.UUID  `gorm:"primaryKey;<-:create"`
	EmployeeID       uuid.UUID  `gorm:"primaryKey;<-:create"`
	DeputyEmployeeID *uuid.UUID `gorm:"<-"`
	CreatedAt        time.Time  `gorm:"<-:create"`
	UpdatedAt        time.Time  `gorm:"<-"`
}

func (OrgAdmin) TableName() string { return "domain.org_admins" }
