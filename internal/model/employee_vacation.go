package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null/v6"
)

// EmployeeVacation is a child row of Employee. Non-overlap is
// enforced at the schema level via an EXCLUDE USING GIST constraint.
// Deleting the parent Employee cascades here.
type EmployeeVacation struct {
	ID         uuid.UUID `gorm:"primaryKey;<-:create"`
	EmployeeID uuid.UUID `gorm:"<-:create"`
	StartsAt   time.Time `gorm:"<-:create"`
	EndsAt     null.Time `gorm:"<-"`
	CreatedAt  time.Time `gorm:"<-:create"`
	UpdatedAt  time.Time `gorm:"<-"`
}

func (EmployeeVacation) TableName() string { return "domain.employee_vacations" }
