-- migrate:up
CREATE TABLE outbox.events (
    id           BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    subject      TEXT NOT NULL,
    payload      BYTEA NOT NULL,
    headers      JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    published_at TIMESTAMPTZ
);

-- Publisher drains unpublished rows in created_at order.
CREATE INDEX outbox_events_unpublished_idx
    ON outbox.events (created_at)
    WHERE published_at IS NULL;

-- migrate:down
DROP TABLE IF EXISTS outbox.events;
