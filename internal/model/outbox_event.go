package model

import "time"

// OutboxEvent is one row in outbox.events. Payload holds the serialized
// event.v1.Envelope; Subject holds the NATS subject used by the drainer
// when shipping the row.
//
// The command service NEVER mutates an existing row: Append creates,
// the separate publisher service updates published_at. Subject/Payload
// are create-only; published_at is not modelled at all here because
// this service has no business reading or writing it.
//
// The outbox.events.headers column is not mapped here: the command
// service writes no headers, and the publisher derives its NATS
// Nats-Msg-Id from the row id (GENERATED ALWAYS AS IDENTITY).
type OutboxEvent struct {
	// ID is assigned by Postgres via GENERATED ALWAYS AS IDENTITY. Do NOT
	// set this field on INSERT — GORM omits zero-value primary keys from
	// the INSERT and Postgres generates the value. Hand-setting it will
	// fail at runtime with "cannot insert into a generated-always column".
	ID        int64     `gorm:"primaryKey"`
	Subject   string    `gorm:"<-:create"`
	Payload   []byte    `gorm:"<-:create"`
	CreatedAt time.Time `gorm:"<-:create"`
}

func (OutboxEvent) TableName() string { return "outbox.events" }
