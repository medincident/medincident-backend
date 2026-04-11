package organizationapp

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/rs/zerolog"

	"github.com/medincident/medincident-command-service/internal/orgstructure/organization"
	"github.com/medincident/medincident-command-service/internal/outbox"
	"github.com/medincident/medincident-command-service/internal/shared/clock"
	"github.com/medincident/medincident-command-service/internal/shared/geo"
	"github.com/medincident/medincident-command-service/internal/shared/tx"
)

// Service is the public port of the Organization application layer.
// All consumers (HTTP handlers, gRPC handlers, tests) depend on this
// interface. The concrete implementation is unexported.
type Service interface {
	Create(ctx context.Context, cmd CreateCommand) (CreateResult, error)
	Rename(ctx context.Context, cmd RenameCommand) (RenameResult, error)
	UpdateDescription(ctx context.Context, cmd UpdateDescriptionCommand) (UpdateDescriptionResult, error)
	RelocateLegalAddress(ctx context.Context, cmd RelocateLegalAddressCommand) (RelocateLegalAddressResult, error)
}

type service struct {
	logger    *zerolog.Logger
	beginner  tx.Beginner
	repo      Repository
	publisher outbox.Publisher
	clk       clock.Clock
}

var _ Service = (*service)(nil)

// NewService wires up an Organization application service. Returns the
// concrete *service; the DI provider exposes it under the Service
// interface to consumers.
func NewService(
	logger *zerolog.Logger,
	bg tx.Beginner,
	repo Repository,
	pub outbox.Publisher,
	clk clock.Clock,
) *service {
	return &service{
		logger:    logger,
		beginner:  bg,
		repo:      repo,
		publisher: pub,
		clk:       clk,
	}
}

// Create builds the address VO and the Organization aggregate,
// collecting every validation failure via errors.Join so the client
// sees all field violations in one response, then persists and
// publishes atomically.
func (s *service) Create(ctx context.Context, cmd CreateCommand) (CreateResult, error) {
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

	now := s.clk.Now().UTC()
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

func (s *service) Rename(ctx context.Context, cmd RenameCommand) (RenameResult, error) {
	err := s.mutate(ctx, cmd.ID, func(o *organization.Organization) error {
		return o.Rename(cmd.Name, s.clk.Now().UTC())
	})
	return RenameResult{}, err
}

// UpdateDescription handles the UpdateDescription command. An empty
// new description clears the stored description.
func (s *service) UpdateDescription(ctx context.Context, cmd UpdateDescriptionCommand) (UpdateDescriptionResult, error) {
	err := s.mutate(ctx, cmd.ID, func(o *organization.Organization) error {
		return o.UpdateDescription(cmd.Description, s.clk.Now().UTC())
	})
	return UpdateDescriptionResult{}, err
}

// RelocateLegalAddress handles the RelocateLegalAddress command. A nil
// Address input means "remove the address".
func (s *service) RelocateLegalAddress(ctx context.Context, cmd RelocateLegalAddressCommand) (RelocateLegalAddressResult, error) {
	var newAddr *geo.Address
	if cmd.Address != nil {
		addr, err := buildAddress(cmd.Address)
		if err != nil {
			return RelocateLegalAddressResult{}, err
		}
		newAddr = addr
	}
	err := s.mutate(ctx, cmd.ID, func(o *organization.Organization) error {
		return o.RelocateLegalAddress(newAddr, s.clk.Now().UTC())
	})
	return RelocateLegalAddressResult{}, err
}

// mutate is the shared load-mutate-save-publish skeleton used by every
// command that targets an existing aggregate. Everything runs inside a
// single transaction.
func (s *service) mutate(
	ctx context.Context,
	id uuid.UUID,
	apply func(*organization.Organization) error,
) error {
	t, err := s.beginner.Begin(ctx)
	if err != nil {
		return err
	}
	return tx.Within(ctx, t, func(ctx context.Context, t tx.Tx) error {
		org, err := s.repo.Find(ctx, id)
		if err != nil {
			return err
		}
		if err := apply(org); err != nil {
			return err
		}
		if err := s.repo.Save(ctx, org); err != nil {
			return err
		}
		return s.publisher.Publish(ctx, t, org)
	})
}

// persistAndPublish is the Create-specific atomic path: save a fresh
// aggregate and publish its Created event in a single transaction.
func (s *service) persistAndPublish(ctx context.Context, org *organization.Organization) error {
	t, err := s.beginner.Begin(ctx)
	if err != nil {
		return err
	}
	return tx.Within(ctx, t, func(ctx context.Context, t tx.Tx) error {
		if err := s.repo.Save(ctx, org); err != nil {
			return err
		}
		return s.publisher.Publish(ctx, t, org)
	})
}

// buildAddress constructs the domain VO from a non-nil AddressInput,
// running both NewPoint and NewAddress unconditionally so every field
// violation surfaces in one joined error.
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
