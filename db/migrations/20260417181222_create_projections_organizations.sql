-- migrate:up
CREATE TABLE projections.organizations (
    id                       UUID PRIMARY KEY,
    name                     TEXT        NOT NULL,
    description              TEXT,
    legal_address_text       TEXT        NOT NULL,
    legal_address_longitude  DOUBLE PRECISION,
    legal_address_latitude   DOUBLE PRECISION,
    created_at               TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at               TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX organizations_name_idx
    ON projections.organizations USING btree (name);

-- migrate:down
DROP TABLE IF EXISTS projections.organizations;
