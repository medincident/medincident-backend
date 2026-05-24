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

// BufferPriority mirrors domain.buffer_priority.
type BufferPriority string

const (
	BufferPriorityNormal BufferPriority = "normal"
	BufferPriorityHigh   BufferPriority = "high"
)

// PatientIncidentBuffer is a write-side row in
// domain.patient_incident_buffer. CategoryID and TypeID are mutable
// while status is pending. After a transition out of pending the row
// is frozen.
type PatientIncidentBuffer struct {
	ID                   uuid.UUID      `gorm:"primaryKey;<-:create"`
	OrganizationID       uuid.UUID      `gorm:"<-:create"`
	PatientZitadelUserID string         `gorm:"<-:create"`
	CategoryID           uuid.NullUUID  `gorm:"<-"`
	TypeID               uuid.NullUUID  `gorm:"<-"`
	Description          string         `gorm:"<-"`
	Summary              string         `gorm:"<-"`
	Priority             BufferPriority `gorm:"<-"`
	OccurredAt           null.Time      `gorm:"<-"`
	Status               BufferStatus   `gorm:"<-"`
	PublishedIncidentID  uuid.NullUUID  `gorm:"<-"`
	CreatedAt            time.Time      `gorm:"<-:create"`
	UpdatedAt            time.Time      `gorm:"<-"`
}

// TableName binds PatientIncidentBuffer to its domain table.
func (PatientIncidentBuffer) TableName() string { return "domain.patient_incident_buffer" }
