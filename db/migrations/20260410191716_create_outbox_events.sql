-- migrate:up
CREATE TABLE outbox.events (
    id              BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    event_id        UUID        NOT NULL,
    occurred_at     TIMESTAMPTZ NOT NULL,
    aggregate_type  TEXT        NOT NULL,
    aggregate_id    TEXT        NOT NULL,
    subject         TEXT        NOT NULL,
    headers         JSONB       NOT NULL,
    payload         BYTEA       NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL,
    published_at    TIMESTAMPTZ
);

-- Publisher drives through unpublished rows in created_at order.
CREATE INDEX outbox_events_unpublished_idx
    ON outbox.events (created_at)
    WHERE published_at IS NULL;

-- Aggregate history lookups for audit / debugging.
CREATE INDEX outbox_events_aggregate_idx
    ON outbox.events (aggregate_type, aggregate_id, created_at);

-- migrate:down
DROP TABLE IF EXISTS outbox.events;
