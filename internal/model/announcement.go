package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null/v6"
)

// AnnouncementPriority mirrors announcement_priority enum.
type AnnouncementPriority string

const (
	AnnouncementPriorityNormal AnnouncementPriority = "normal"
	AnnouncementPriorityHigh   AnnouncementPriority = "high"
)

// Announcement is a write-side row in domain.announcements.
// OrganizationID, ClinicID, DepartmentID, AuthorID are immutable after creation.
// The nullable scope fields: ClinicID=NULL means org-level, DepartmentID=NULL means at most clinic-level.
type Announcement struct {
	ID             uuid.UUID            `gorm:"primaryKey;<-:create"`
	OrganizationID uuid.UUID            `gorm:"<-:create"`
	ClinicID       uuid.NullUUID        `gorm:"<-:create"`
	DepartmentID   uuid.NullUUID        `gorm:"<-:create"`
	AuthorID       string               `gorm:"<-:create"`
	Title          string               `gorm:"<-"`
	Content        string               `gorm:"<-"`
	Priority       AnnouncementPriority `gorm:"<-"`
	IsArchived     bool                 `gorm:"<-"`
	StartsAt       null.Time            `gorm:"<-"`
	EndsAt         null.Time            `gorm:"<-"`
	CreatedAt      time.Time            `gorm:"<-:create"`
	UpdatedAt      time.Time            `gorm:"<-"`
}

// TableName binds Announcement to domain.announcements.
func (Announcement) TableName() string { return "domain.announcements" }
