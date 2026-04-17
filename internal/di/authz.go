package di

import (
	"github.com/samber/do/v2"
	"gorm.io/gorm"

	"github.com/medincident/medincident-command-service/internal/service/authz"
)

func provideAuthz(injector do.Injector) (*authz.Authz, error) {
	db, err := do.Invoke[*gorm.DB](injector)
	if err != nil {
		return nil, err
	}
	return authz.New(db), nil
}
