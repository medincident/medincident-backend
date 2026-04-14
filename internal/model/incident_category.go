package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null/v6"
)

// IncidentCategory is a node in an organisation's incident classifier
// tree. Unlike clinics/departments, IncidentCategory's parent FK is
// mutable: the classifier explicitly supports MoveCategory, and the
// service layer enforces cycle and max-depth invariants on each move.
// OrganizationID is denormalised from the parent chain specifically so
// that the partial unique index
// (organization_id, name) WHERE is_active can enforce per-organisation
// name uniqueness for active rows at the database level.
type IncidentCategory struct {
	ID               uuid.UUID     `gorm:"primaryKey;<-:create"`
	OrganizationID   uuid.UUID     `gorm:"<-:create"`
	ParentCategoryID uuid.NullUUID `gorm:"<-"`
	Name             string        `gorm:"<-"`
	Description      null.String   `gorm:"<-"`
	IsActive         bool          `gorm:"<-"`
	CreatedAt        time.Time     `gorm:"<-:create"`
	UpdatedAt        time.Time     `gorm:"<-"`
}

func (IncidentCategory) TableName() string { return "domain.incident_categories" }
