package identity

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/nats-io/nats.go/jetstream"
	"github.com/rs/zerolog"
	"github.com/samber/oops"
	"google.golang.org/protobuf/proto"

	zitadelv1 "github.com/medincident/medincident-zitadel-actions/gen/zitadel/events/v1"
	sessionsv1 "github.com/medincident/medincident-zitadel-actions/gen/zitadel/sessions/v1"
	usersv1 "github.com/medincident/medincident-zitadel-actions/gen/zitadel/users/v1"

	"github.com/medincident/medincident-backend/internal/config"
)

// Zitadel event_type discriminators emitted by the zitadel-actions
// gateway.
const (
	ZitadelEventUserHumanAdded          = "user.human.added"
	ZitadelEventUserHumanProfileChanged = "user.human.profile.changed"
	ZitadelEventUserHumanEmailChanged   = "user.human.email.changed"
	ZitadelEventUserHumanEmailVerified  = "user.human.email.verified"
	ZitadelEventSessionAdded            = "session.added"
	ZitadelEventSessionUserChecked      = "session.user.checked"
)

// Error codes emitted by consumer runtime paths. Malformed payloads
// end with _malformed so the interceptor maps them to permanent
// termination; transient DB failures surface with _failed.
const (
	ErrCodeStreamLookupFailed      = "stream_lookup_failed"
	ErrCodeConsumerCreateFailed    = "consumer_create_failed"
	ErrCodeConsumeStartFailed      = "consume_start_failed"
	ErrCodeEnvelopeUnmarshalFailed = "envelope_unmarshal_malformed"
	ErrCodePayloadUnmarshalFailed  = "payload_unmarshal_malformed"
)

// nakRetryDelay is the delay applied when a dispatch fails with a
// transient error (DB blip etc.). Kept short so backlogs clear quickly
// once the downstream recovers.
const nakRetryDelay = 5 * time.Second

// Consumer is the NATS JetStream consumer for the zitadel.> stream.
// Lifecycle is managed via Start / Shutdown; Start spins up a Consume
// loop in the background, Shutdown halts the loop and drains in-flight
// messages.
type Consumer struct {
	js        jetstream.JetStream
	cfg       *config.NATSConfig
	projector *Projector
	logger    *zerolog.Logger

	mu       sync.Mutex
	started  bool
	stopped  bool
	consume  jetstream.ConsumeContext
	ctx      context.Context //nolint:containedctx // stored for per-message cancellation.
	cancel   context.CancelFunc
	inflight sync.WaitGroup
}

// NewConsumer returns a new Consumer wired to the given JetStream
// client and projector.
func NewConsumer(
	js jetstream.JetStream,
	cfg *config.NATSConfig,
	projector *Projector,
	logger *zerolog.Logger,
) *Consumer {
	return &Consumer{js: js, cfg: cfg, projector: projector, logger: logger}
}

// Start subscribes a durable JetStream consumer to the configured
// zitadel stream and begins processing messages in the background.
// Idempotent: calling Start twice returns nil on the second call.
func (c *Consumer) Start(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.started {
		return nil
	}

	stream, err := c.js.Stream(ctx, c.cfg.Stream)
	if err != nil {
		return oops.In("consumer.identity").
			Code(ErrCodeStreamLookupFailed).
			With("stream", c.cfg.Stream).
			Wrap(err)
	}

	filter := ""
	if len(c.cfg.Subjects) == 1 {
		filter = c.cfg.Subjects[0]
	}
	cons, err := stream.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
		Durable:           c.cfg.DurableName,
		Name:              c.cfg.DurableName,
		AckPolicy:         jetstream.AckExplicitPolicy,
		AckWait:           consumerAckWait(nakRetryDelay),
		DeliverPolicy:     jetstream.DeliverAllPolicy,
		MaxDeliver:        -1,
		FilterSubject:     filter,
		FilterSubjects:    subjectsFilter(c.cfg.Subjects),
		MaxAckPending:     1024,
		InactiveThreshold: 24 * time.Hour,
	})
	if err != nil {
		return oops.In("consumer.identity").
			Code(ErrCodeConsumerCreateFailed).
			With("stream", c.cfg.Stream).
			With("durable", c.cfg.DurableName).
			Wrap(err)
	}

	c.ctx, c.cancel = context.WithCancel(context.Background())
	cc, err := cons.Consume(c.handleMsg)
	if err != nil {
		c.cancel()
		return oops.In("consumer.identity").
			Code(ErrCodeConsumeStartFailed).
			With("stream", c.cfg.Stream).
			With("durable", c.cfg.DurableName).
			Wrap(err)
	}
	c.consume = cc
	c.started = true
	c.logger.Info().
		Str("stream", c.cfg.Stream).
		Str("durable", c.cfg.DurableName).
		Msg("identity consumer started")
	return nil
}

// Shutdown stops the consumer loop and waits for in-flight messages
// to finish. ctx's deadline bounds the drain wait.
func (c *Consumer) Shutdown(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if !c.started || c.stopped {
		return nil
	}
	c.stopped = true
	if c.consume != nil {
		c.consume.Stop()
	}
	if c.cancel != nil {
		c.cancel()
	}
	done := make(chan struct{})
	go func() {
		c.inflight.Wait()
		close(done)
	}()
	select {
	case <-done:
		c.logger.Info().Msg("identity consumer stopped")
		return nil
	case <-ctx.Done():
		c.logger.Warn().Msg("identity consumer shutdown timeout")
		return ctx.Err()
	}
}

// handleMsg runs once per incoming message.
func (c *Consumer) handleMsg(msg jetstream.Msg) {
	c.inflight.Add(1)
	defer c.inflight.Done()

	ctx, cancel := context.WithCancel(c.ctx)
	defer cancel()

	err := c.dispatch(ctx, msg)
	if err == nil {
		if ackErr := msg.Ack(); ackErr != nil {
			c.logger.Warn().Err(ackErr).Msg("identity consumer ack failed")
		}
		return
	}

	// Shutdown: don't ack. Redelivery on next boot.
	if errors.Is(err, context.Canceled) {
		return
	}

	// Permanent failure (malformed payload) → term so server stops
	// redelivering. Transient failure → nak with backoff.
	if isMalformed(err) {
		c.logger.Error().
			Err(err).
			Str("subject", msg.Subject()).
			Msg("identity dispatch malformed; term")
		if termErr := msg.TermWithReason(err.Error()); termErr != nil {
			c.logger.Warn().Err(termErr).Msg("identity consumer term failed")
		}
		return
	}

	c.logger.Error().
		Err(err).
		Str("subject", msg.Subject()).
		Msg("identity dispatch failed; nak with backoff")
	if nakErr := msg.NakWithDelay(nakRetryDelay); nakErr != nil {
		c.logger.Warn().Err(nakErr).Msg("identity consumer nak failed")
	}
}

// dispatch decodes the envelope and routes to the right projector.
func (c *Consumer) dispatch(ctx context.Context, msg jetstream.Msg) error {
	env := new(zitadelv1.Envelope)
	if err := proto.Unmarshal(msg.Data(), env); err != nil {
		return oops.In("consumer.identity").
			Code(ErrCodeEnvelopeUnmarshalFailed).
			Wrap(err)
	}
	return c.dispatchZitadel(ctx, env)
}

// dispatchZitadel handles one decoded envelope.
func (c *Consumer) dispatchZitadel(ctx context.Context, env *zitadelv1.Envelope) error {
	occurredAt := envelopeOccurredAt(env.GetCreatedAt())
	aggregateID := env.GetAggregateId()
	eventType := env.GetEventType()

	// EmailVerified is a marker event with no payload — short-circuit.
	if eventType == ZitadelEventUserHumanEmailVerified {
		return c.projector.ApplyUserHumanEmailVerified(ctx, aggregateID, occurredAt)
	}

	payload := env.GetPayload()
	if payload == nil {
		return oops.In("consumer.identity").
			Code(ErrCodePayloadUnmarshalFailed).
			With("event_type", eventType).
			With("aggregate_id", aggregateID).
			Errorf("nil payload")
	}
	msg, err := payload.UnmarshalNew()
	if err != nil {
		return oops.In("consumer.identity").
			Code(ErrCodePayloadUnmarshalFailed).
			With("type_url", payload.GetTypeUrl()).
			With("event_type", eventType).
			Wrap(err)
	}

	switch eventType {
	case ZitadelEventUserHumanAdded:
		ev, ok := msg.(*usersv1.UserHumanAdded)
		if !ok {
			return typeMismatch(eventType, msg)
		}
		return c.projector.ApplyUserHumanAdded(ctx, aggregateID, occurredAt, ev)
	case ZitadelEventUserHumanProfileChanged:
		ev, ok := msg.(*usersv1.UserHumanProfileChanged)
		if !ok {
			return typeMismatch(eventType, msg)
		}
		return c.projector.ApplyUserHumanProfileChanged(ctx, aggregateID, occurredAt, ev)
	case ZitadelEventUserHumanEmailChanged:
		ev, ok := msg.(*usersv1.UserHumanEmailChanged)
		if !ok {
			return typeMismatch(eventType, msg)
		}
		return c.projector.ApplyUserHumanEmailChanged(ctx, aggregateID, occurredAt, ev)
	case ZitadelEventSessionAdded:
		ev, ok := msg.(*sessionsv1.SessionAdded)
		if !ok {
			return typeMismatch(eventType, msg)
		}
		return c.projector.ApplySessionAdded(ctx, aggregateID, occurredAt, ev)
	case ZitadelEventSessionUserChecked:
		ev, ok := msg.(*sessionsv1.SessionUserChecked)
		if !ok {
			return typeMismatch(eventType, msg)
		}
		return c.projector.ApplySessionUserChecked(ctx, aggregateID, occurredAt, ev)
	default:
		c.logger.Warn().
			Str("event_type", eventType).
			Str("aggregate_id", aggregateID).
			Msg("unknown zitadel event type; skipping")
		return nil
	}
}

// typeMismatch returns a malformed-payload error when the unmarshaled
// payload's concrete type does not match event_type. Always a
// publisher bug or schema drift.
func typeMismatch(eventType string, got any) error {
	return oops.In("consumer.identity").
		Code(ErrCodePayloadUnmarshalFailed).
		With("event_type", eventType).
		Errorf("payload type %T does not match event_type", got)
}

// isMalformed reports whether an error is a permanent decode failure.
func isMalformed(err error) bool {
	oe, ok := oops.AsOops(err)
	if !ok {
		return false
	}
	code := oe.Code()
	return code == ErrCodeEnvelopeUnmarshalFailed || code == ErrCodePayloadUnmarshalFailed
}

// envelopeOccurredAt converts a protobuf timestamp to time.Time,
// defaulting to time.Now().UTC() when the publisher forgot to set one.
func envelopeOccurredAt(ts timestampLike) time.Time {
	if ts == nil || !ts.IsValid() {
		return time.Now().UTC()
	}
	return ts.AsTime().UTC()
}

// timestampLike is the minimal interface the protobuf Timestamp type
// satisfies — avoids importing timestamppb here directly.
type timestampLike interface {
	AsTime() time.Time
	IsValid() bool
}

// consumerAckWait gives the projector a little more time than the nak
// delay so a message is not redelivered while the previous Nak is
// still inflight on the server side.
func consumerAckWait(nakDelay time.Duration) time.Duration {
	minAckWait := 30 * time.Second
	if nakDelay*2 > minAckWait {
		return nakDelay * 2
	}
	return minAckWait
}

// subjectsFilter returns the FilterSubjects slice to pass to the
// JetStream API. When a single subject is configured we prefer the
// singular FilterSubject field for broader server-compat; otherwise
// return the slice as-is.
func subjectsFilter(subjects []string) []string {
	if len(subjects) <= 1 {
		return nil
	}
	return append([]string(nil), subjects...)
}
