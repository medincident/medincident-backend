-- migrate:up

CREATE TABLE projections.incident_categories (
    id                  UUID        PRIMARY KEY,
    organization_id     UUID        NOT NULL,
    parent_category_id  UUID,
    name                TEXT        NOT NULL,
    description         TEXT,
    is_active           BOOLEAN     NOT NULL DEFAULT TRUE,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX incident_categories_organization_id_idx
    ON projections.incident_categories (organization_id);

CREATE INDEX incident_categories_parent_category_id_idx
    ON projections.incident_categories (parent_category_id);

CREATE UNIQUE INDEX incident_categories_org_active_name_uniq
    ON projections.incident_categories (organization_id, name)
    WHERE is_active;

CREATE TABLE projections.incident_types (
    id                      UUID        PRIMARY KEY,
    organization_id         UUID        NOT NULL,
    category_id             UUID        NOT NULL,
    name                    TEXT        NOT NULL,
    description             TEXT,
    is_active               BOOLEAN     NOT NULL DEFAULT TRUE,
    is_allowed_for_patients BOOLEAN     NOT NULL DEFAULT FALSE,
    created_at              TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX incident_types_organization_id_idx
    ON projections.incident_types (organization_id);

CREATE INDEX incident_types_category_id_idx
    ON projections.incident_types (category_id);

CREATE UNIQUE INDEX incident_types_org_active_name_uniq
    ON projections.incident_types (organization_id, name)
    WHERE is_active;

-- migrate:down
DROP TABLE IF EXISTS projections.incident_types;
DROP TABLE IF EXISTS projections.incident_categories;
