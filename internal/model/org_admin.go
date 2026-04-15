package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null/v6"
)

// OrgAdmin links an organisation to an employee holding the
// organisation-admin role. Primary key is composite.
type OrgAdmin struct {
	OrganizationID   uuid.UUID             `gorm:"primaryKey;<-:create"`
	EmployeeID       uuid.UUID             `gorm:"primaryKey;<-:create"`
	DeputyEmployeeID null.Value[uuid.UUID] `gorm:"<-"`
	CreatedAt        time.Time             `gorm:"<-:create"`
	UpdatedAt        time.Time             `gorm:"<-"`
}

func (OrgAdmin) TableName() string { return "domain.org_admins" }
