-- migrate:up
CREATE TABLE domain.organizations (
    id                       UUID PRIMARY KEY,
    name                     TEXT        NOT NULL,
    description              TEXT,
    legal_address_text       TEXT        NOT NULL,
    legal_address_longitude  DOUBLE PRECISION,
    legal_address_latitude   DOUBLE PRECISION,
    created_at               TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at               TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- migrate:down
DROP TABLE IF EXISTS domain.organizations;
