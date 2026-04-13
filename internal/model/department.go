package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null/v6"
)

// Department belongs to exactly one Clinic; that link is set at
// creation time and never changes. Departments have no address.
type Department struct {
	ID          uuid.UUID   `gorm:"primaryKey;<-:create"`
	ClinicID    uuid.UUID   `gorm:"<-:create"`
	Name        string      `gorm:"<-"`
	Description null.String `gorm:"<-"`
	CreatedAt   time.Time   `gorm:"<-:create"`
	UpdatedAt   time.Time   `gorm:"<-"`
}

func (Department) TableName() string { return "domain.departments" }
