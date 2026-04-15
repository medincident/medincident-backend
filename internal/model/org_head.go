package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null/v6"
)

// OrgHead links an organisation to an employee holding the
// organisation-head role. Primary key is composite.
type OrgHead struct {
	OrganizationID   uuid.UUID             `gorm:"primaryKey;<-:create"`
	EmployeeID       uuid.UUID             `gorm:"primaryKey;<-:create"`
	DeputyEmployeeID null.Value[uuid.UUID] `gorm:"<-"`
	CreatedAt        time.Time             `gorm:"<-:create"`
	UpdatedAt        time.Time             `gorm:"<-"`
}

func (OrgHead) TableName() string { return "domain.org_heads" }
