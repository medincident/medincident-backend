package model

import "time"

// SystemAdmin marks a Zitadel user as a global system administrator.
// No foreign keys — the zitadel_user_id is an opaque reference to
// Zitadel, not ours to manage. No UpdatedAt — the row has no
// mutable state.
type SystemAdmin struct {
	ZitadelUserID string    `gorm:"primaryKey;<-:create;column:zitadel_user_id"`
	CreatedAt     time.Time `gorm:"<-:create"`
}

func (SystemAdmin) TableName() string { return "domain.system_admins" }
