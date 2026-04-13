package model

import (
	"time"

	"gorm.io/datatypes"
)

// OutboxEvent is one row in outbox.events. Payload holds the serialized
// medincident.event.v1.Envelope; Subject holds the NATS subject used by
// the drainer when shipping the row.
//
// The command service NEVER mutates an existing row: Append creates,
// the separate publisher service updates published_at. Subject/Payload/
// Headers are create-only; published_at is not modelled at all here
// because this service has no business reading or writing it.
type OutboxEvent struct {
	ID        int64          `gorm:"primaryKey"`
	Subject   string         `gorm:"<-:create"`
	Payload   []byte         `gorm:"<-:create"`
	Headers   datatypes.JSON `gorm:"<-:create"`
	CreatedAt time.Time      `gorm:"<-:create"`
}

func (OutboxEvent) TableName() string { return "outbox.events" }
