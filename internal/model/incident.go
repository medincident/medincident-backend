package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null/v6"
)

// IncidentStatus mirrors domain.incident_status.
type IncidentStatus string

const (
	IncidentStatusPending    IncidentStatus = "pending"
	IncidentStatusInProgress IncidentStatus = "in_progress"
	IncidentStatusDone       IncidentStatus = "done"
	IncidentStatusRejected   IncidentStatus = "rejected"
	IncidentStatusCancelled  IncidentStatus = "cancelled"
)

// IsTerminal reports whether the status forbids further mutations.
func (s IncidentStatus) IsTerminal() bool {
	switch s {
	case IncidentStatusDone, IncidentStatusRejected, IncidentStatusCancelled:
		return true
	default:
		return false
	}
}

// IncidentPriority mirrors domain.incident_priority.
type IncidentPriority string

const (
	IncidentPriorityLow      IncidentPriority = "low"
	IncidentPriorityNormal   IncidentPriority = "normal"
	IncidentPriorityHigh     IncidentPriority = "high"
	IncidentPriorityCritical IncidentPriority = "critical"
)

// Incident is a write-side row in domain.incidents. organization_id and
// clinic_id are denormalized from the parent department at creation
// time and are immutable afterwards. CategoryID and TypeID are also
// immutable — once set they describe what happened, not what is
// currently classified.
type Incident struct {
	ID                         uuid.UUID        `gorm:"primaryKey;<-:create"`
	OrganizationID             uuid.UUID        `gorm:"<-:create"`
	ClinicID                   uuid.UUID        `gorm:"<-:create"`
	DepartmentID               uuid.UUID        `gorm:"<-:create"`
	CategoryID                 uuid.UUID        `gorm:"<-:create"`
	TypeID                     uuid.UUID        `gorm:"<-:create"`
	Status                     IncidentStatus   `gorm:"<-"`
	Priority                   IncidentPriority `gorm:"<-"`
	Description                null.String      `gorm:"<-"`
	PatientOriginalDescription null.String      `gorm:"<-:create"`
	OccurredAt                 time.Time        `gorm:"<-:create"`
	RegistrarEmployeeID        uuid.UUID        `gorm:"<-:create"`
	SourcePatientZitadelUserID null.String      `gorm:"<-:create"`
	SourceBufferID             uuid.NullUUID    `gorm:"<-:create"`
	ReopenedFromIncidentID     uuid.NullUUID    `gorm:"<-:create"`
	CreatedAt                  time.Time        `gorm:"<-:create"`
	UpdatedAt                  time.Time        `gorm:"<-"`
}

// TableName binds Incident to domain.incidents.
func (Incident) TableName() string { return "domain.incidents" }
