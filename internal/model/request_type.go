package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null/v6"
)

type RequestType struct {
	ID             uuid.UUID   `gorm:"primaryKey;<-:create"`
	OrganizationID uuid.UUID   `gorm:"<-:create"`
	Name           string      `gorm:"<-"`
	Description    null.String `gorm:"<-"`
	IsActive       bool        `gorm:"<-"`
	CreatedAt      time.Time   `gorm:"<-:create"`
	UpdatedAt      time.Time   `gorm:"<-"`
}

func (RequestType) TableName() string { return "domain.request_types" }
