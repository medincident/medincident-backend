package model

import (
	"time"

	"github.com/google/uuid"
)

type ServiceRequestStatus string

const (
	ServiceRequestStatusCreated       ServiceRequestStatus = "created"
	ServiceRequestStatusInWork        ServiceRequestStatus = "in_work"
	ServiceRequestStatusOnHold        ServiceRequestStatus = "on_hold"
	ServiceRequestStatusPendingReview ServiceRequestStatus = "pending_review"
	ServiceRequestStatusCompleted     ServiceRequestStatus = "completed"
	ServiceRequestStatusCancelled     ServiceRequestStatus = "cancelled"
)

type ServiceRequest struct {
	ID             uuid.UUID            `gorm:"primaryKey;<-:create"`
	OrganizationID uuid.UUID            `gorm:"<-:create"`
	ClinicID       uuid.UUID            `gorm:"<-:create"`
	DepartmentID   uuid.UUID            `gorm:"<-:create"`
	TypeID         uuid.UUID            `gorm:"<-:create"`
	IncidentID     uuid.NullUUID        `gorm:"<-:create"`
	Description    string               `gorm:"<-"`
	Status         ServiceRequestStatus `gorm:"<-"`
	AuthorID       string               `gorm:"<-:create"`
	CreatedAt      time.Time            `gorm:"<-:create"`
	UpdatedAt      time.Time            `gorm:"<-"`
}

func (ServiceRequest) TableName() string { return "domain.service_requests" }

type ServiceRequestExecutor struct {
	ID           uuid.UUID `gorm:"primaryKey;<-:create"`
	RequestID    uuid.UUID `gorm:"<-:create"`
	EmployeeID   uuid.UUID `gorm:"<-:create"`
	AssignedAt   time.Time `gorm:"<-:create"`
	AssignedByID string    `gorm:"<-:create"`
}

func (ServiceRequestExecutor) TableName() string { return "domain.service_request_executors" }
