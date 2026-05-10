package domain

import (
	"context"
	"time"

	"github.com/nats-io/nats.go/jetstream"
	"github.com/rs/zerolog"
	"github.com/samber/oops"
	"google.golang.org/protobuf/proto"
	"gorm.io/gorm"

	clinicv1 "github.com/medincident/medincident-backend/pkg/event/clinic/v1"
	deptv1 "github.com/medincident/medincident-backend/pkg/event/department/v1"
	orgv1 "github.com/medincident/medincident-backend/pkg/event/organization/v1"
	eventv1 "github.com/medincident/medincident-backend/pkg/event/v1"

	"github.com/medincident/medincident-backend/internal/service/query/projector"
)

// Dispatcher decodes an Envelope, resolves the event type from the NATS
// subject, and calls the matching projector function inside a DB transaction.
type Dispatcher struct {
	db   *gorm.DB
	proj *projector.Projectors
	log  *zerolog.Logger
}

// NewDispatcher returns a Dispatcher wired to the given DB and projector set.
func NewDispatcher(db *gorm.DB, proj *projector.Projectors, log *zerolog.Logger) *Dispatcher {
	return &Dispatcher{db: db, proj: proj, log: log}
}

// Dispatch decodes msg, applies the projector inside a transaction, and
// returns nil on success. Returns a malformed error (_malformed code) for
// permanent decode failures; transient DB errors are returned unwrapped so
// the caller naks with backoff.
func (d *Dispatcher) Dispatch(ctx context.Context, msg jetstream.Msg) error {
	env := new(eventv1.Envelope)
	if err := proto.Unmarshal(msg.Data(), env); err != nil {
		return oops.In("consumer.domain").
			Code(ErrCodeEnvelopeUnmarshalFailed).
			With("subject", msg.Subject()).
			Wrap(err)
	}
	aggregateID := env.GetAggregateId()
	occurredAt := envelopeTime(env)

	payload := env.GetPayload()
	if payload == nil {
		return oops.In("consumer.domain").
			Code(ErrCodePayloadUnmarshalFailed).
			With("subject", msg.Subject()).
			Errorf("nil payload")
	}
	evMsg, err := payload.UnmarshalNew()
	if err != nil {
		return oops.In("consumer.domain").
			Code(ErrCodePayloadUnmarshalFailed).
			With("subject", msg.Subject()).
			Wrap(err)
	}

	return d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		switch m := evMsg.(type) {
		// ── Organization ──────────────────────────────────────────────────
		case *orgv1.OrganizationCreated:
			return d.proj.OrganizationCreated(tx, aggregateID, occurredAt, m)
		case *orgv1.OrganizationDetailsChanged:
			return d.proj.OrganizationDetailsChanged(tx, aggregateID, occurredAt, m)
		case *orgv1.OrganizationLegalAddressChanged:
			return d.proj.OrganizationLegalAddressChanged(tx, aggregateID, occurredAt, m)

		// ── Clinic ────────────────────────────────────────────────────────
		case *clinicv1.ClinicCreated:
			return d.proj.ClinicCreated(tx, aggregateID, occurredAt, m)
		case *clinicv1.ClinicDetailsChanged:
			return d.proj.ClinicDetailsChanged(tx, aggregateID, occurredAt, m)
		case *clinicv1.ClinicPhysicalAddressChanged:
			return d.proj.ClinicPhysicalAddressChanged(tx, aggregateID, occurredAt, m)

		// ── Department ────────────────────────────────────────────────────
		case *deptv1.DepartmentCreated:
			return d.proj.DepartmentCreated(tx, aggregateID, occurredAt, m)
		case *deptv1.DepartmentDetailsChanged:
			return d.proj.DepartmentDetailsChanged(tx, aggregateID, occurredAt, m)

		default:
			d.log.Warn().
				Str("subject", msg.Subject()).
				Str("type", payload.GetTypeUrl()).
				Msg("unknown domain event type; skipping")
			return nil
		}
	})
}

// envelopeTime converts the envelope's occurred_at timestamp to time.Time,
// defaulting to time.Now().UTC() when the publisher forgot to set one.
func envelopeTime(env *eventv1.Envelope) time.Time {
	ts := env.GetOccurredAt()
	if ts == nil || !ts.IsValid() {
		return time.Now().UTC()
	}
	return ts.AsTime().UTC()
}
