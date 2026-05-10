-- migrate:up
CREATE SCHEMA IF NOT EXISTS outbox;

CREATE TABLE outbox.events (
    seq          BIGSERIAL    PRIMARY KEY,
    id           UUID         NOT NULL DEFAULT gen_random_uuid() UNIQUE,
    subject      TEXT         NOT NULL,
    payload      BYTEA        NOT NULL,
    headers      JSONB,
    created_at   TIMESTAMPTZ  NOT NULL DEFAULT now(),
    published_at TIMESTAMPTZ
);

CREATE INDEX outbox_events_unpublished_seq_idx
    ON outbox.events (seq)
    WHERE published_at IS NULL;

-- migrate:down
DROP TABLE IF EXISTS outbox.events;
DROP SCHEMA IF EXISTS outbox;
