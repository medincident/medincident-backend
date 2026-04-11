-- migrate:up
CREATE TABLE outbox.events (
    id             UUID PRIMARY KEY,
    aggregate_type TEXT NOT NULL,
    aggregate_id   TEXT NOT NULL,
    event_type     TEXT NOT NULL,
    payload        JSONB NOT NULL,
    created_at     TIMESTAMPTZ NOT NULL,
    published_at   TIMESTAMPTZ
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
