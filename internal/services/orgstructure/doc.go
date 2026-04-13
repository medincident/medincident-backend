// Package orgstructure contains the command-side business logic for
// the Organization, Clinic, and Department aggregates. Each method
// of each service lives in its own file (<aggregate>_<method>.go);
// shared validators and the outbox helper live in address.go and
// outbox.go respectively. There are no repository interfaces —
// services call *gorm.DB directly and open transactions inline.
package orgstructure
