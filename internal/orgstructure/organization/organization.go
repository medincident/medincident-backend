// Package organization is the Organization aggregate within the
// orgstructure bounded context. It exposes the aggregate root, its
// domain events, and its validators. Application-layer types (commands,
// results, service, repository interface) live in the app sub-package;
// infra adapters (proto mappers) live in the infra sub-package.
package organization

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/samber/oops"

	"github.com/medincident/medincident-command-service/internal/shared/aggregate"
	"github.com/medincident/medincident-command-service/internal/shared/geo"
)

// ErrCodeOrganizationIDGenerationFailed is emitted when uuid.NewV7
// fails inside the lifecycle constructor.
const ErrCodeOrganizationIDGenerationFailed = "id_generation_failed"

// Organization is the aggregate root. UpdatedAt is only ever moved via
// Root.Raise — direct writes are a contract violation.
type Organization struct {
	aggregate.Root

	ID           uuid.UUID
	Name         string
	Description  string       // "" = not set
	LegalAddress *geo.Address // nil = not set
}

// New is the primary lifecycle constructor. It validates name and
// description (multi-error via errors.Join) and trusts legalAddress as
// already-built (the caller is expected to have constructed it via
// geo.NewAddress or decided it's nil).
func New(
	name string,
	description string,
	legalAddress *geo.Address,
	now time.Time,
) (*Organization, error) {
	name = trim(name)
	description = trim(description)

	var errs []error
	if err := validateName(name); err != nil {
		errs = append(errs, err)
	}
	if err := validateDescription(description); err != nil {
		errs = append(errs, err)
	}
	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	id, err := uuid.NewV7()
	if err != nil {
		return nil, oops.In("orgstructure.organization").
			Code(ErrCodeOrganizationIDGenerationFailed).
			Public("Failed to create organization. Please try again.").
			Hint("uuid.NewV7 returned an error — check system clock / entropy source").
			Wrap(err)
	}

	o := &Organization{
		Root:         aggregate.NewRoot(now),
		ID:           id,
		Name:         name,
		Description:  description,
		LegalAddress: legalAddress,
	}
	o.Raise(&Created{
		ID:           id,
		Name:         name,
		Description:  description,
		LegalAddress: legalAddress,
		At:           now,
	}, now)
	return o, nil
}

// Hydrate restores an Organization from persisted state. It does NOT
// validate (the repository is trusted) and does NOT raise events.
// Repositories reach it via the app.Repository interface.
func Hydrate(
	id uuid.UUID,
	name string,
	description string,
	legalAddress *geo.Address,
	createdAt, updatedAt time.Time,
) *Organization {
	return &Organization{
		Root:         aggregate.HydrateRoot(createdAt, updatedAt),
		ID:           id,
		Name:         name,
		Description:  description,
		LegalAddress: legalAddress,
	}
}

// AggregateType is the stable short name used in the outbox envelope.
// Renaming is a breaking change for consumers that filter by it.
func (o *Organization) AggregateType() string { return "organization" }

// AggregateID renders the organization id as a string for the
// outbox.EventSource contract; PullEvents is inherited from aggregate.Root.
func (o *Organization) AggregateID() string { return o.ID.String() }

// Rename changes the organization's name. Validates the new name,
// no-ops if identical, raises Renamed on successful change.
func (o *Organization) Rename(name string, now time.Time) error {
	name = trim(name)
	if err := validateName(name); err != nil {
		return err
	}
	if o.Name == name {
		return nil
	}
	o.Name = name
	o.Raise(&Renamed{ID: o.ID, Name: name, At: now}, now)
	return nil
}

// UpdateDescription replaces the description. Empty string clears it.
// No-ops if the new description equals the current one. Raises
// DescriptionUpdated on successful change.
func (o *Organization) UpdateDescription(description string, now time.Time) error {
	description = trim(description)
	if err := validateDescription(description); err != nil {
		return err
	}
	if o.Description == description {
		return nil
	}
	o.Description = description
	o.Raise(&DescriptionUpdated{
		ID:          o.ID,
		Description: description,
		At:          now,
	}, now)
	return nil
}

// RelocateLegalAddress replaces the legal address wholesale. nil is a
// valid input meaning "remove the address". No-ops if the new value is
// equal to the current one. Raises LegalAddressRelocated on successful
// change. The VO is trusted: the caller built it via geo.NewAddress or
// decided it's nil.
func (o *Organization) RelocateLegalAddress(address *geo.Address, now time.Time) error {
	if o.LegalAddress == nil && address == nil {
		return nil
	}
	if o.LegalAddress != nil && address != nil && o.LegalAddress.Equal(*address) {
		return nil
	}
	o.LegalAddress = address
	o.Raise(&LegalAddressRelocated{ID: o.ID, LegalAddress: address, At: now}, now)
	return nil
}
