package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null/v6"
)

// OrgDispatcher links an organisation to an employee holding the
// dispatcher role. Primary key is composite.
type OrgDispatcher struct {
	OrganizationID   uuid.UUID             `gorm:"primaryKey;<-:create"`
	EmployeeID       uuid.UUID             `gorm:"primaryKey;<-:create"`
	DeputyEmployeeID null.Value[uuid.UUID] `gorm:"<-"`
	CreatedAt        time.Time             `gorm:"<-:create"`
	UpdatedAt        time.Time             `gorm:"<-"`
}

func (OrgDispatcher) TableName() string { return "domain.org_dispatchers" }
