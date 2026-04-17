package authz

import (
	"github.com/samber/oops"
	"gorm.io/gorm"
)

const (
	ErrCodePermissionDenied   = "permission_denied"
	ErrCodeCallerNotFound     = "caller_not_found"
	ErrCodeScopeResolveFailed = "scope_resolve_failed"
)

func errPermissionDenied(callerID string) error {
	return oops.In("service.authz").
		Code(ErrCodePermissionDenied).
		Public("Permission denied.").
		With("caller_id", callerID).
		Errorf("permission denied")
}

type Authz struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Authz {
	return &Authz{db: db}
}
