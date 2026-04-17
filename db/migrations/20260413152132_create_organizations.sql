-- migrate:up
CREATE TYPE domain.geo_point AS (
    longitude DOUBLE PRECISION,
    latitude  DOUBLE PRECISION
);

CREATE TYPE domain.address AS (
    text  TEXT,
    point domain.geo_point
);

CREATE TABLE domain.organizations (
    id                UUID PRIMARY KEY,
    name              TEXT        NOT NULL,
    description       TEXT,
    legal_address     domain.address NOT NULL,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- migrate:down
DROP TABLE IF EXISTS domain.organizations;
DROP TYPE IF EXISTS domain.address;
DROP TYPE IF EXISTS domain.geo_point;
