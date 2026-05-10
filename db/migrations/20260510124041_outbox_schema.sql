-- migrate:up
CREATE SCHEMA IF NOT EXISTS outbox;

CREATE TABLE outbox.events (
    seq          bigserial    PRIMARY KEY,
    id           uuid         NOT NULL DEFAULT gen_random_uuid() UNIQUE,
    subject      text         NOT NULL,
    payload      bytea        NOT NULL,
    created_at   timestamptz  NOT NULL DEFAULT now(),
    published_at timestamptz
);

CREATE INDEX ON outbox.events (seq) WHERE published_at IS NULL;

-- migrate:down
DROP TABLE IF EXISTS outbox.events;
DROP SCHEMA IF EXISTS outbox;
