// Package domain provides the NATS JetStream consumer that applies domain
// events to the query-side projections. It is the query-side counterpart
// of the internal/service/query/identity consumer; the structure is
// intentionally identical.
package domain

import (
	"context"
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/nats-io/nats.go/jetstream"
	"github.com/rs/zerolog"
	"github.com/samber/oops"
	"gorm.io/gorm"

	"github.com/medincident/medincident-backend/internal/config"
	"github.com/medincident/medincident-backend/internal/service/query/projector"
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

// Consumer is the NATS JetStream consumer for the medincident_events stream.
// Lifecycle is managed via Start / Shutdown; Start spins up a Consume
// loop in the background, Shutdown halts the loop and drains in-flight
// messages.
type Consumer struct {
	js         jetstream.JetStream
	cfg        *config.NATSConfig
	db         *gorm.DB
	dispatcher *Dispatcher
	logger     *zerolog.Logger

	mu       sync.Mutex
	started  bool
	stopped  bool
	consume  jetstream.ConsumeContext
	ctx      context.Context //nolint:containedctx // stored for per-message cancellation.
	cancel   context.CancelFunc
	inflight sync.WaitGroup
}

// NewConsumer returns a Consumer wired to the given JetStream client and
// query DB.
func NewConsumer(
	js jetstream.JetStream,
	cfg *config.NATSConfig,
	db *gorm.DB,
	proj *projector.Projectors,
	logger *zerolog.Logger,
) *Consumer {
	return &Consumer{
		js:         js,
		cfg:        cfg,
		db:         db,
		dispatcher: NewDispatcher(db, proj, logger),
		logger:     logger,
	}
}

// Start subscribes a durable JetStream consumer to the configured stream.
// Idempotent: calling Start twice returns nil on the second call.
func (c *Consumer) Start(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.started {
		return nil
	}

	stream, err := c.js.Stream(ctx, c.cfg.Stream)
	if err != nil {
		return oops.In("consumer.domain").
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
		MaxAckPending:     256,
		InactiveThreshold: 24 * time.Hour,
	})
	if err != nil {
		return oops.In("consumer.domain").
			Code(ErrCodeConsumerCreateFailed).
			With("stream", c.cfg.Stream).
			With("durable", c.cfg.DurableName).
			Wrap(err)
	}

	c.ctx, c.cancel = context.WithCancel(ctx)
	cc, err := cons.Consume(c.handleMsg)
	if err != nil {
		c.cancel()
		return oops.In("consumer.domain").
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
		Msg("domain consumer started")
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
		c.logger.Info().Msg("domain consumer stopped")
		return nil
	case <-ctx.Done():
		c.logger.Warn().Msg("domain consumer shutdown timeout")
		return ctx.Err()
	}
}

// handleMsg runs once per incoming message.
func (c *Consumer) handleMsg(msg jetstream.Msg) {
	c.inflight.Add(1)
	defer c.inflight.Done()

	ctx, cancel := context.WithCancel(c.ctx)
	defer cancel()

	err := c.dispatcher.Dispatch(ctx, msg)
	if err == nil {
		if ackErr := msg.Ack(); ackErr != nil {
			c.logger.Warn().Err(ackErr).Msg("domain consumer ack failed")
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
			Msg("domain dispatch malformed; term")
		if termErr := msg.TermWithReason(err.Error()); termErr != nil {
			c.logger.Warn().Err(termErr).Msg("domain consumer term failed")
		}
		return
	}

	c.logger.Error().
		Err(err).
		Str("subject", msg.Subject()).
		Msg("domain dispatch failed; nak with backoff")
	if nakErr := msg.NakWithDelay(nakRetryDelay); nakErr != nil {
		c.logger.Warn().Err(nakErr).Msg("domain consumer nak failed")
	}
}

// isMalformed reports whether an error is a permanent decode failure.
func isMalformed(err error) bool {
	oe, ok := oops.AsOops(err)
	if !ok {
		return false
	}
	code, ok := oe.Code().(string)
	if !ok {
		return false
	}
	return strings.HasSuffix(code, "_malformed")
}

// consumerAckWait gives the projector a little more time than the nak
// delay so a message is not redelivered while the previous Nak is
// still inflight on the server side.
func consumerAckWait(nakDelay time.Duration) time.Duration {
	const minAckWait = 30 * time.Second
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
