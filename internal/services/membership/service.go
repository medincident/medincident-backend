// Package membership implements command-side business logic for
// membership-style entities (currently: Employee + Vacation). Future
// entities (department heads, clinic heads) will live alongside
// EmployeeService in this package.
package membership

import (
	"github.com/rs/zerolog"
	"gorm.io/gorm"

	"github.com/medincident/medincident-command-service/internal/services/zitadel"
)

// EmployeeService holds the command-side methods for managing
// employees and their vacations. Every method opens its own gorm
// transaction and publishes one or more events through the outbox in
// that same transaction.
type EmployeeService struct {
	db       *gorm.DB
	verifier *zitadel.Service
	logger   *zerolog.Logger
}

// NewEmployeeService wires an EmployeeService. The verifier is called
// exactly once inside HireEmployee, before the DB transaction.
func NewEmployeeService(db *gorm.DB, verifier *zitadel.Service, logger *zerolog.Logger) *EmployeeService {
	return &EmployeeService{db: db, verifier: verifier, logger: logger}
}
