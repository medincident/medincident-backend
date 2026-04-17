-- migrate:up

-- Read-model copy of domain.incident_types. No cross-schema FKs.
CREATE TABLE projections.incident_types (
    id               UUID        PRIMARY KEY,
    organization_id  UUID        NOT NULL,
    category_id      UUID        NOT NULL,
    name             TEXT        NOT NULL,
    description      TEXT,
    is_active        BOOLEAN     NOT NULL DEFAULT TRUE,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX incident_types_organization_id_idx
    ON projections.incident_types (organization_id);

CREATE INDEX incident_types_category_id_idx
    ON projections.incident_types (category_id);

-- migrate:down
DROP TABLE IF EXISTS projections.incident_types;
