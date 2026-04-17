-- migrate:up
CREATE TYPE domain.geo_point AS (
    longitude DOUBLE PRECISION,
    latitude  DOUBLE PRECISION
);

CREATE TABLE domain.organizations (
    id                       UUID PRIMARY KEY,
    name                     TEXT        NOT NULL,
    description              TEXT,
    legal_address_text       TEXT        NOT NULL,
    legal_address_point      domain.geo_point,
    created_at               TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at               TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- migrate:down
DROP TABLE IF EXISTS domain.organizations;
DROP TYPE IF EXISTS domain.geo_point;
