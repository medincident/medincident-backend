package organizationapp

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"

	"github.com/medincident/medincident-command-service/internal/orgstructure/organization"
	"github.com/medincident/medincident-command-service/internal/outbox"
	"github.com/medincident/medincident-command-service/internal/shared/geo"
	"github.com/medincident/medincident-command-service/internal/tx"
)

// Clock is a minimal injection point for time.Now. Injected via the
// constructor so tests can run with a deterministic clock and so
// aggregates see a single "now" for a single request.
type Clock interface {
	Now() time.Time
}

// beginner is the local interface the service depends on for opening
// transactions. The concrete implementation (storage/postgres.Beginner)
// satisfies it. Declaring the interface at the consumer is idiomatic
// Go — the service does not import the postgres package at all.
type beginner interface {
	Begin(ctx context.Context) (tx.Tx, error)
}

// Service is the Organization application service. It orchestrates the
// aggregate, its repository, and the outbox in a single transaction.
//
// Every public method takes an XCommand and returns an XResult —
// MANDATORY per project convention. Domain constructors use positional
// parameters; the application-layer boundary always uses structs.
type Service struct {
	logger      *zerolog.Logger
	beginner    beginner
	repo        Repository
	outboxStore outbox.Store
	outboxReg   outbox.Registry
	clock       Clock
}

// NewService wires up an Organization application service.
func NewService(
	logger *zerolog.Logger,
	bg beginner,
	repo Repository,
	store outbox.Store,
	reg outbox.Registry,
	clock Clock,
) *Service {
	return &Service{
		logger:      logger,
		beginner:    bg,
		repo:        repo,
		outboxStore: store,
		outboxReg:   reg,
		clock:       clock,
	}
}

// Create handles the CreateOrganization command end to end:
//
//  1. Build the *geo.Address VO from the input and run the Organization
//     aggregate constructor, collecting ALL validation errors in one
//     pass via errors.Join. Operation-level validation never stops at
//     the first failure — the client must see every field violation.
//  2. If validation succeeded, open a transaction, save the aggregate,
//     publish its events to the outbox in the same tx, commit.
func (s *Service) Create(ctx context.Context, cmd CreateCommand) (CreateResult, error) {
	var (
		errs         []error
		legalAddress *geo.Address
	)

	if cmd.LegalAddress != nil {
		addr, addrErr := buildAddress(cmd.LegalAddress)
		if addrErr != nil {
			errs = append(errs, addrErr)
		} else {
			legalAddress = addr
		}
	}

	now := s.clock.Now().UTC()
	org, orgErr := organization.New(cmd.Name, cmd.Description, legalAddress, now)
	if orgErr != nil {
		errs = append(errs, orgErr)
	}

	if len(errs) > 0 {
		return CreateResult{}, errors.Join(errs...)
	}

	if err := s.persistAndPublish(ctx, org); err != nil {
		return CreateResult{}, err
	}
	return CreateResult{ID: org.ID}, nil
}

// Rename handles the RenameOrganization command.
func (s *Service) Rename(ctx context.Context, cmd RenameCommand) (RenameResult, error) {
	err := s.mutate(ctx, cmd.ID, func(o *organization.Organization) error {
		return o.Rename(cmd.NewName, s.clock.Now().UTC())
	})
	return RenameResult{}, err
}

// UpdateDescription handles the UpdateDescription command. An empty
// new description clears the stored description.
func (s *Service) UpdateDescription(ctx context.Context, cmd UpdateDescriptionCommand) (UpdateDescriptionResult, error) {
	err := s.mutate(ctx, cmd.ID, func(o *organization.Organization) error {
		return o.UpdateDescription(cmd.NewDescription, s.clock.Now().UTC())
	})
	return UpdateDescriptionResult{}, err
}

// RelocateLegalAddress handles the RelocateLegalAddress command. A nil
// Address input means "remove the address".
func (s *Service) RelocateLegalAddress(ctx context.Context, cmd RelocateLegalAddressCommand) (RelocateLegalAddressResult, error) {
	var newAddr *geo.Address
	if cmd.Address != nil {
		addr, err := buildAddress(cmd.Address)
		if err != nil {
			return RelocateLegalAddressResult{}, err
		}
		newAddr = addr
	}
	err := s.mutate(ctx, cmd.ID, func(o *organization.Organization) error {
		return o.RelocateLegalAddress(newAddr, s.clock.Now().UTC())
	})
	return RelocateLegalAddressResult{}, err
}

// mutate is the shared "load + mutate + save + publish-events" skeleton
// used by every command that targets an existing aggregate. It opens a
// transaction, loads the aggregate, runs the mutator, saves, publishes
// the drained events to the outbox, and commits — all atomically.
func (s *Service) mutate(
	ctx context.Context,
	id uuid.UUID,
	apply func(*organization.Organization) error,
) error {
	t, err := s.beginner.Begin(ctx)
	if err != nil {
		return err
	}
	return tx.Within(ctx, t, func(ctx context.Context, t tx.Tx) error {
		org, err := s.repo.GetByID(ctx, id)
		if err != nil {
			return err
		}
		if err := apply(org); err != nil {
			return err
		}
		if err := s.repo.Save(ctx, org); err != nil {
			return err
		}
		return outbox.Publish(ctx, t, s.outboxStore, s.outboxReg, org)
	})
}

// persistAndPublish is the Create-specific atomic path: save a fresh
// aggregate and publish its Created event in a single transaction.
func (s *Service) persistAndPublish(ctx context.Context, org *organization.Organization) error {
	t, err := s.beginner.Begin(ctx)
	if err != nil {
		return err
	}
	return tx.Within(ctx, t, func(ctx context.Context, t tx.Tx) error {
		if err := s.repo.Save(ctx, org); err != nil {
			return err
		}
		return outbox.Publish(ctx, t, s.outboxStore, s.outboxReg, org)
	})
}

// buildAddress is the glue between the transport-shaped AddressInput
// and the domain *geo.Address VO. It runs BOTH NewPoint and NewAddress
// regardless of individual failures so that every field violation in
// the address tree is surfaced in one go — no fail-fast.
//
// Callers must pre-check for a nil input and skip the call entirely
// when the address is absent — buildAddress always builds a real VO.
func buildAddress(in *AddressInput) (*geo.Address, error) {
	var (
		errs  []error
		point *geo.Point
	)
	if in.Point != nil {
		p, err := geo.NewPoint(in.Point.Longitude, in.Point.Latitude)
		if err != nil {
			errs = append(errs, err)
		} else {
			point = &p
		}
	}
	addr, err := geo.NewAddress(in.Text, point)
	if err != nil {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}
	return &addr, nil
}
