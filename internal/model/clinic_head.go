package model

import (
	"time"

	"github.com/google/uuid"
)

// ClinicHead links a clinic to the employee acting as its head.
// Primary key is composite (clinic_id, employee_id).
type ClinicHead struct {
	ClinicID         uuid.UUID  `gorm:"primaryKey;<-:create"`
	EmployeeID       uuid.UUID  `gorm:"primaryKey;<-:create"`
	DeputyEmployeeID *uuid.UUID `gorm:"<-"`
	CreatedAt        time.Time  `gorm:"<-:create"`
	UpdatedAt        time.Time  `gorm:"<-"`
}

func (ClinicHead) TableName() string { return "domain.clinic_heads" }
