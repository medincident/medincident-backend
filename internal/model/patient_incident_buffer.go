package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null/v6"
)

// BufferStatus mirrors domain.buffer_status.
type BufferStatus string

const (
	BufferStatusPending   BufferStatus = "pending"
	BufferStatusPublished BufferStatus = "published"
	BufferStatusRejected  BufferStatus = "rejected"
	BufferStatusCancelled BufferStatus = "cancelled"
)

// PatientIncidentBuffer is a write-side row in
// domain.patient_incident_buffer. CategoryID and TypeID are mutable
// while status is pending — both the patient and the dispatcher may
// edit them. After a transition out of pending the row is frozen.
type PatientIncidentBuffer struct {
	ID                   uuid.UUID     `gorm:"primaryKey;<-:create"`
	OrganizationID       uuid.UUID     `gorm:"<-:create"`
	PatientZitadelUserID string        `gorm:"<-:create"`
	CategoryID           uuid.NullUUID `gorm:"<-"`
	TypeID               uuid.NullUUID `gorm:"<-"`
	Description          null.String   `gorm:"<-"`
	OccurredAt           null.Time     `gorm:"<-"`
	Status               BufferStatus  `gorm:"<-"`
	PublishedIncidentID  uuid.NullUUID `gorm:"<-"`
	CreatedAt            time.Time     `gorm:"<-:create"`
	UpdatedAt            time.Time     `gorm:"<-"`
}

// TableName binds PatientIncidentBuffer to its domain table.
func (PatientIncidentBuffer) TableName() string { return "domain.patient_incident_buffer" }
