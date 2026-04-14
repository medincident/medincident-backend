package model

import (
	"time"

	"github.com/google/uuid"
)

// DepartmentResponsible links a department to the employee acting as
// its responsible. The primary key is composite (department_id,
// employee_id); holder uniqueness across departments is enforced by
// a separate unique index on employee_id.
type DepartmentResponsible struct {
	DepartmentID     uuid.UUID  `gorm:"primaryKey;<-:create"`
	EmployeeID       uuid.UUID  `gorm:"primaryKey;<-:create"`
	DeputyEmployeeID *uuid.UUID `gorm:"<-"`
	CreatedAt        time.Time  `gorm:"<-:create"`
	UpdatedAt        time.Time  `gorm:"<-"`
}

func (DepartmentResponsible) TableName() string { return "domain.department_responsibles" }
