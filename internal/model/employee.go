package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null/v6"
)

// Employee is a membership record linking a Zitadel user to a specific
// department inside a specific organisation. organization_id is
// denormalised from the department's parent lineage at hire time and is
// immutable afterwards.
type Employee struct {
	ID             uuid.UUID   `gorm:"primaryKey;<-:create"`
	ZitadelUserID  string      `gorm:"<-:create;column:zitadel_user_id"`
	OrganizationID uuid.UUID   `gorm:"<-:create"`
	DepartmentID   uuid.UUID   `gorm:"<-"`
	Position       null.String `gorm:"<-"`
	IsActive       bool        `gorm:"<-"`
	CreatedAt      time.Time   `gorm:"<-:create"`
	UpdatedAt      time.Time   `gorm:"<-"`
}

func (Employee) TableName() string { return "domain.employees" }
